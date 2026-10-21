package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const RegistrationsRoute = "/api/v1/proofs/registrations"

const maxRegistrationBody = 4 << 10

var registrationTypes = map[string][]typedDataField{
	"EIP712Domain": {
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"RegistrationRequest": {
		{Name: "holder", Type: "address"},
		{Name: "principal", Type: "string"},
		{Name: "issuedAt", Type: "uint64"},
	},
}

type RegistrationHandler struct {
	auth      Authenticator
	registrar Registrant
	log       *slog.Logger
}

func NewRegistrationHandler(auth Authenticator, registrar Registrant, log *slog.Logger) *RegistrationHandler {
	return &RegistrationHandler{auth: auth, registrar: registrar, log: log}
}

func (h *RegistrationHandler) ServeTemplate(w http.ResponseWriter, r *http.Request) {
	principal, err := h.auth.Principal(r.Header.Get("Authorization"))
	if err != nil {
		h.fail(w, r, err)

		return
	}

	domain := h.registrar.Domain()
	writeJSON(w, http.StatusOK, registrationTemplate{
		Types:       registrationTypes,
		PrimaryType: "RegistrationRequest",
		Domain: templateDomain{
			Name:              domain.Name,
			Version:           domain.Version,
			ChainID:           domain.ChainID,
			VerifyingContract: domain.VerifyingContract.Hex(),
		},
		Message: map[string]string{"principal": principal.String()},
	})
}

func (h *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	principal, err := h.auth.Principal(r.Header.Get("Authorization"))
	if err != nil {
		h.fail(w, r, err)

		return
	}

	request, signature, err := readRegistration(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", notary.MessageOf(err))

		return
	}

	request.Principal = principal

	issued, err := h.registrar.Register(r.Context(), request, signature)
	if err != nil {
		h.fail(w, r, err)

		return
	}

	response := registrationResponse{
		Court: courtResponse{
			ChainID: strconv.FormatUint(issued.Court.ChainID, 10),
			Address: issued.Court.Address.Hex(),
		},
		Epoch:         strconv.FormatUint(issued.Epoch, 10),
		Registrations: make([]issuedRegistration, 0, len(issued.Registrations)),
	}

	for _, one := range issued.Registrations {
		response.Registrations = append(response.Registrations, issuedRegistration{
			Holder:    one.Registration.Holder.Hex(),
			Account:   one.Registration.Account.Hex(),
			FromEpoch: strconv.FormatUint(one.Registration.FromEpoch, 10),
			Kind:      one.Kind.String(),
			Digest:    one.Digest.Hex(),
			Signature: one.Signature.Hex(),
		})
	}

	writeJSON(w, http.StatusOK, response)
}

func readRegistration(body io.Reader) (notary.RegistrationRequest, notary.Signature, error) {
	const op = "httpapi.readRegistration"

	decoder := json.NewDecoder(io.LimitReader(body, maxRegistrationBody+1))
	decoder.DisallowUnknownFields()

	var raw registrationBody
	if err := decoder.Decode(&raw); err != nil {
		return notary.RegistrationRequest{}, notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"the body is not a registration request",
		)
	}

	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return notary.RegistrationRequest{}, notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"the body carries more than one registration request",
		)
	}

	holder, err := notary.ParseAddress(raw.Holder)
	if err != nil {
		return notary.RegistrationRequest{}, notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"holder is not an address",
		)
	}

	if holder.IsZero() {
		return notary.RegistrationRequest{}, notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			nil,
			"holder is the zero address",
		)
	}

	seconds, err := strconv.ParseUint(raw.IssuedAt, 10, 64)
	if err != nil {
		return notary.RegistrationRequest{}, notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"issuedAt is not a decimal number of seconds",
		)
	}

	issuedAt, err := notary.UnixTimeFromSeconds(seconds)
	if err != nil {
		return notary.RegistrationRequest{}, notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"issuedAt is not a time",
		)
	}

	signature, err := notary.ParseSignature(raw.Signature)
	if err != nil {
		return notary.RegistrationRequest{}, notary.Signature{}, notary.Errf(
			notary.KindInvalidArgument,
			op,
			err,
			"signature is not 65 bytes of hex",
		)
	}

	return notary.RegistrationRequest{Holder: holder, IssuedAt: issuedAt}, signature, nil
}

func (h *RegistrationHandler) fail(w http.ResponseWriter, r *http.Request, err error) {
	status, code := registrationStatusOf(notary.KindOf(err))
	message := notary.MessageOf(err)

	if message == "" {
		message = http.StatusText(status)
	}

	if status >= http.StatusInternalServerError {
		h.log.ErrorContext(
			r.Context(),
			"a registration failed",
			"kind",
			notary.KindOf(err).String(),
			"err",
			err.Error(),
		)

		message = "the registration could not be issued"

		switch code {
		case "unavailable":
			message = "the notary cannot read its records right now"
		case "signer_unavailable":
			message = "the notary cannot sign a registration right now"
		case "chain_unavailable":
			message = "the notary cannot reach the chain right now"
		case "identity_unavailable":
			message = identityUnavailableMessage
		}
	}

	writeError(w, status, code, message)
}

func registrationStatusOf(kind notary.Kind) (int, string) {
	switch kind {
	case notary.KindUnauthenticated:
		return http.StatusUnauthorized, "unauthenticated"
	case notary.KindIdentityUnavailable:
		return http.StatusServiceUnavailable, "identity_unavailable"
	case notary.KindRegistrationRefused:
		return http.StatusForbidden, "registration_refused"
	case notary.KindNoAccount:
		return http.StatusNotFound, "no_account"
	case notary.KindRegistrationNotOpen:
		return http.StatusConflict, "registration_not_open"
	case notary.KindStoreUnavailable:
		return http.StatusServiceUnavailable, "unavailable"
	case notary.KindSignerUnavailable:
		return http.StatusServiceUnavailable, "signer_unavailable"
	case notary.KindLedgerUnavailable:
		return http.StatusServiceUnavailable, "chain_unavailable"
	case notary.KindDisagreement:
		return http.StatusInternalServerError, "disagreement"
	default:
		return http.StatusInternalServerError, "internal"
	}
}
