package sponsoring

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const estimationSignature = "fffffffffffffffffffffffffffffff0000000000000000000000000000000007" +
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1c"

type Sponsor struct {
	chain     notary.SponsorChain
	paymaster notary.Paymaster
	bundler   notary.Bundler
	signer    notary.HolderSigner
	scheme    notary.OperationScheme
	calls     notary.SponsoredCalls
	poll      time.Duration
	patience  time.Duration
}

func NewSponsor(
	chain notary.SponsorChain,
	paymaster notary.Paymaster,
	bundler notary.Bundler,
	signer notary.HolderSigner,
	scheme notary.OperationScheme,
	calls notary.SponsoredCalls,
	poll time.Duration,
	patience time.Duration,
) *Sponsor {
	return &Sponsor{
		chain:     chain,
		paymaster: paymaster,
		bundler:   bundler,
		signer:    signer,
		scheme:    scheme,
		calls:     calls,
		poll:      poll,
		patience:  patience,
	}
}

func (s *Sponsor) Holder() notary.Address {
	return s.signer.Holder()
}

func (s *Sponsor) Calls() notary.SponsoredCalls {
	return s.calls
}

func (s *Sponsor) Identity(ctx context.Context) (notary.PaymasterIdentity, error) {
	return s.paymaster.Identity(ctx)
}

func (s *Sponsor) Submit(
	ctx context.Context,
	identity notary.PaymasterIdentity,
	call notary.SponsoredCall,
) (notary.OperationReceipt, error) {
	const op = "sponsoring.Sponsor.Submit"

	if err := s.admitted(ctx, call); err != nil {
		return notary.OperationReceipt{}, err
	}

	operation, delegation, err := s.prepare(ctx, identity, call)
	if err != nil {
		return notary.OperationReceipt{}, err
	}

	gas, err := s.bundler.Estimate(ctx, operation, identity.EntryPoint, delegation)
	if err != nil {
		return notary.OperationReceipt{}, err
	}

	operation.PreVerificationGas = gas.PreVerificationGas
	operation.VerificationGasLimit = gas.VerificationGasLimit
	operation.CallGasLimit = gas.CallGasLimit
	operation.PaymasterVerificationGasLimit = gas.PaymasterVerificationGasLimit

	if operation.MaxCost().Cmp(identity.Ceiling) > 0 {
		return notary.OperationReceipt{}, notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"the operation may cost up to %s wei, more than the paymaster's ceiling of %s",
			operation.MaxCost(),
			identity.Ceiling,
		)
	}

	hash := s.scheme.OperationHash(
		notary.EntryPointIdentity{ChainID: identity.ChainID, Address: identity.EntryPoint},
		operation,
	)

	signature, err := s.signer.SignOperation(hash)
	if err != nil {
		return notary.OperationReceipt{}, err
	}

	operation.Signature = signature[:]

	sent, err := s.bundler.Send(ctx, operation, identity.EntryPoint, delegation)
	if err != nil {
		return notary.OperationReceipt{}, err
	}

	if sent != hash {
		return notary.OperationReceipt{}, notary.Errf(
			notary.KindDisagreement,
			op,
			nil,
			"the bundler took the operation as %s, this sponsor hashed it as %s",
			sent.Hex(),
			hash.Hex(),
		)
	}

	return s.await(ctx, hash)
}

func (s *Sponsor) admitted(ctx context.Context, call notary.SponsoredCall) error {
	window, err := s.paymaster.Sponsorship(ctx, s.signer.Holder(), s.calls.Execute(call), call.PaymasterData)
	if err != nil {
		return err
	}

	now, err := s.chain.BlockTime(ctx)
	if err != nil {
		return err
	}

	return window.Admit(now)
}

func (s *Sponsor) prepare(
	ctx context.Context,
	identity notary.PaymasterIdentity,
	call notary.SponsoredCall,
) (notary.UserOperation, *notary.Delegation, error) {
	const op = "sponsoring.Sponsor.prepare"

	holder := s.signer.Holder()

	delegate, err := s.chain.DelegateOf(ctx, holder)
	if err != nil {
		return notary.UserOperation{}, nil, err
	}

	operation := notary.UserOperation{
		Sender:        holder,
		CallData:      s.calls.Execute(call),
		Paymaster:     identity.Address,
		PaymasterData: call.PaymasterData,
	}

	var delegation *notary.Delegation

	switch {
	case delegate == identity.Account:
	case delegate.IsZero():
		nonce, err := s.chain.AccountNonce(ctx, holder)
		if err != nil {
			return notary.UserOperation{}, nil, err
		}

		signed, err := s.signer.Authorize(identity.ChainID, identity.Account, nonce)
		if err != nil {
			return notary.UserOperation{}, nil, err
		}

		operation.Delegate = identity.Account
		delegation = &signed
	default:
		return notary.UserOperation{}, nil, notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"%s delegates to %s; the paymaster sponsors only accounts delegated to %s",
			holder.Hex(),
			delegate.Hex(),
			identity.Account.Hex(),
		)
	}

	operation.Nonce, err = s.chain.OperationNonce(ctx, identity.EntryPoint, holder)
	if err != nil {
		return notary.UserOperation{}, nil, err
	}

	operation.MaxFeePerGas, operation.MaxPriorityFeePerGas, err = s.chain.Fees(ctx)
	if err != nil {
		return notary.UserOperation{}, nil, err
	}

	operation.Signature, err = hex.DecodeString(estimationSignature)
	if err != nil {
		return notary.UserOperation{}, nil, notary.Errf(notary.KindUnknown, op, err, "")
	}

	return operation, delegation, nil
}

func (s *Sponsor) await(ctx context.Context, hash notary.Digest) (notary.OperationReceipt, error) {
	const op = "sponsoring.Sponsor.await"

	ctx, cancel := context.WithTimeout(ctx, s.patience)
	defer cancel()

	ticker := time.NewTicker(s.poll)
	defer ticker.Stop()

	for {
		receipt, found, err := s.bundler.Receipt(ctx, hash)
		if err == nil && found {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return notary.OperationReceipt{}, notary.Errf(
				notary.KindSourceUnavailable,
				op,
				err,
				"the bundler took operation %s but showed no receipt within %s; it may still land",
				hash.Hex(),
				s.patience,
			)
		case <-ticker.C:
		}
	}
}
