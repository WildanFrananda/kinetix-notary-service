package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholding"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholdinginterface"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type WithholdingCourt struct {
	client      *ethclient.Client
	binding     *withholding.WithholdingCourt
	contract    *bind.BoundContract
	expected    notary.CourtIdentity
	court       notary.Address
	ledger      notary.Address
	interfaceID [4]byte
	timeout     time.Duration
}

var _ notary.WithholdingCourt = (*WithholdingCourt)(nil)

func NewWithholdingCourt(client *ethclient.Client, settings Settings) *WithholdingCourt {
	binding := withholding.NewWithholdingCourt()

	return &WithholdingCourt{
		client:      client,
		binding:     binding,
		contract:    binding.Instance(client, common.Address(settings.Withholding)),
		expected:    notary.CourtIdentity{ChainID: settings.Ledger.ChainID, Address: settings.Withholding},
		court:       settings.Court,
		ledger:      settings.Ledger.Address,
		interfaceID: withholdingInterfaceID(),
		timeout:     settings.Timeout,
	}
}

func (w *WithholdingCourt) Identity(ctx context.Context) (notary.CourtIdentity, error) {
	const op = "ethereum.WithholdingCourt.Identity"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	chainID, err := w.client.ChainID(ctx)
	if err != nil {
		return notary.CourtIdentity{}, unavailable(op, err)
	}

	if !chainID.IsUint64() || chainID.Uint64() != w.expected.ChainID {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"the RPC endpoint serves chain %s, the withholding court is configured on chain %d",
			chainID,
			w.expected.ChainID,
		)
	}

	code, err := w.client.CodeAt(ctx, common.Address(w.expected.Address), nil)
	if err != nil {
		return notary.CourtIdentity{}, unavailable(op, err)
	}

	if len(code) == 0 {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"chain %d has no contract at %s",
			w.expected.ChainID,
			w.expected.Address.Hex(),
		)
	}

	supported, err := probe(
		ctx,
		w.client,
		w.expected.Address,
		w.binding.PackSupportsInterface(w.interfaceID),
		w.binding.UnpackSupportsInterface,
	)
	if err != nil {
		return notary.CourtIdentity{}, w.wrap(op, err)
	}

	if !supported {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"%s does not announce the withholding court interface",
			w.expected.Address.Hex(),
		)
	}

	beside, err := probe(ctx, w.client, w.expected.Address, w.binding.PackCourt(), w.binding.UnpackCourt)
	if err != nil {
		return notary.CourtIdentity{}, w.wrap(op, err)
	}

	if notary.Address(beside) != w.court {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"%s stands beside the court %s, not %s",
			w.expected.Address.Hex(),
			beside.Hex(),
			w.court.Hex(),
		)
	}

	judged, err := probe(ctx, w.client, w.expected.Address, w.binding.PackLedger(), w.binding.UnpackLedger)
	if err != nil {
		return notary.CourtIdentity{}, w.wrap(op, err)
	}

	if notary.Address(judged) != w.ledger {
		return notary.CourtIdentity{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"%s reads the ledger %s, not %s",
			w.expected.Address.Hex(),
			judged.Hex(),
			w.ledger.Hex(),
		)
	}

	return w.expected, nil
}

func (w *WithholdingCourt) RegistrationDigest(
	ctx context.Context,
	registration notary.Registration,
) (notary.Digest, error) {
	const op = "ethereum.WithholdingCourt.RegistrationDigest"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	digest, err := bind.Call(
		w.contract,
		&bind.CallOpts{Context: ctx},
		w.binding.PackRegistrationDigest(boundRegistration(registration)),
		w.binding.UnpackRegistrationDigest,
	)
	if err != nil {
		return notary.Digest{}, w.callFailed(op, err)
	}

	return notary.Digest(digest), nil
}

func (w *WithholdingCourt) LatestBlock(ctx context.Context) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	number, err := w.client.BlockNumber(ctx)
	if err != nil {
		return 0, unavailable("ethereum.WithholdingCourt.LatestBlock", err)
	}

	return number, nil
}

func (w *WithholdingCourt) Requests(ctx context.Context, from, to uint64) ([]notary.WithholdingRequest, error) {
	const op = "ethereum.WithholdingCourt.Requests"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	logs, err := w.client.FilterLogs(ctx, goethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: []common.Address{common.Address(w.expected.Address)},
		Topics:    [][]common.Hash{{w.binding.GetABI().Events["StatementRequested"].ID}},
	})
	if err != nil {
		return nil, unavailable(op, err)
	}

	requests := make([]notary.WithholdingRequest, 0, len(logs))

	for i := range logs {
		event, err := w.binding.UnpackStatementRequestedEvent(&logs[i])
		if err != nil {
			return nil, notary.Errf(notary.KindWrongCourt, op, err, "a StatementRequested log does not decode")
		}

		dueAt, err := notary.UnixTimeFromSeconds(event.DueAt)
		if err != nil {
			return nil, notary.Errf(notary.KindWrongCourt, op, err, "a request carries no valid deadline")
		}

		requests = append(requests, notary.WithholdingRequest{
			Registration: event.Registration,
			Account:      event.Account,
			Epoch:        event.Epoch,
			Holder:       notary.Address(event.Holder),
			DueAt:        dueAt,
		})
	}

	return requests, nil
}

func (w *WithholdingCourt) RequestOf(
	ctx context.Context,
	registration notary.Digest,
	epoch uint64,
) (notary.RequestState, notary.UnixTime, error) {
	const op = "ethereum.WithholdingCourt.RequestOf"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	request, err := bind.Call(
		w.contract,
		&bind.CallOpts{Context: ctx},
		w.binding.PackRequestOf(registration, epoch),
		w.binding.UnpackRequestOf,
	)
	if err != nil {
		return 0, notary.UnixTime{}, w.callFailed(op, err)
	}

	state := notary.RequestState(request.State)
	if state > notary.RequestUnanswered {
		return 0, notary.UnixTime{}, notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"the court reports state %d",
			request.State,
		)
	}

	if state == notary.RequestNone {
		return state, notary.UnixTime{}, nil
	}

	dueAt, err := notary.UnixTimeFromSeconds(request.DueAt)
	if err != nil {
		return 0, notary.UnixTime{}, notary.Errf(notary.KindWrongCourt, op, err, "a request carries no valid deadline")
	}

	return state, dueAt, nil
}

func (w *WithholdingCourt) TotalSilences(ctx context.Context) (uint64, error) {
	const op = "ethereum.WithholdingCourt.TotalSilences"

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	total, err := bind.Call(
		w.contract,
		&bind.CallOpts{Context: ctx},
		w.binding.PackTotalSilences(),
		w.binding.UnpackTotalSilences,
	)
	if err != nil {
		return 0, w.callFailed(op, err)
	}

	if !total.IsUint64() {
		return 0, notary.Errf(notary.KindWrongCourt, op, nil, "the court reports %s silences", total)
	}

	return total.Uint64(), nil
}

func (w *WithholdingCourt) callFailed(op string, err error) error {
	if errors.Is(err, bind.ErrNoCode) {
		return notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"chain %d has no contract at %s",
			w.expected.ChainID,
			w.expected.Address.Hex(),
		)
	}

	return unavailable(op, err)
}

func (w *WithholdingCourt) wrap(op string, err error) error {
	if errors.Is(err, errNotALedger) {
		return notary.Errf(
			notary.KindWrongCourt,
			op,
			nil,
			"%s does not answer as a withholding court",
			w.expected.Address.Hex(),
		)
	}

	return unavailable(op, fmt.Errorf("withholding court: %w", err))
}

func boundRegistration(registration notary.Registration) withholding.Registration {
	return withholding.Registration{
		Holder:    common.Address(registration.Holder),
		Account:   registration.Account,
		FromEpoch: registration.FromEpoch,
	}
}

func withholdingInterfaceID() [4]byte {
	var id [4]byte
	for _, method := range withholdinginterface.NewIWithholdingCourt().GetABI().Methods {
		for i := range id {
			id[i] ^= method.ID[i]
		}
	}

	return id
}
