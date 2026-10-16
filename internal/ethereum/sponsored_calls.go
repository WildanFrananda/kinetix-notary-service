package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/account7702"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/court"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholding"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type SponsoredCalls struct {
	account     *account7702.Simple7702Account
	court       *court.StatementCourt
	withholding *withholding.WithholdingCourt
}

var _ notary.SponsoredCalls = SponsoredCalls{}

func NewSponsoredCalls() SponsoredCalls {
	return SponsoredCalls{
		account:     account7702.NewSimple7702Account(),
		court:       court.NewStatementCourt(),
		withholding: withholding.NewWithholdingCourt(),
	}
}

func (c SponsoredCalls) Execute(call notary.SponsoredCall) []byte {
	return c.account.PackExecute(common.Address(call.Target), new(big.Int), call.Data)
}

func (c SponsoredCalls) RequestStatement(
	target notary.Address,
	issued notary.IssuedRegistration,
	epoch uint64,
) notary.SponsoredCall {
	return notary.SponsoredCall{
		Target: target,
		Data:   c.withholding.PackRequestStatement(registrationOf(issued), issued.Signature[:], epoch),
	}
}

func (c SponsoredCalls) RecordUnanswered(
	target notary.Address,
	issued notary.IssuedRegistration,
	epoch uint64,
) notary.SponsoredCall {
	return notary.SponsoredCall{
		Target:        target,
		Data:          c.withholding.PackRecordUnanswered(registrationOf(issued), epoch),
		PaymasterData: issued.Signature[:],
	}
}

func (c SponsoredCalls) ProveBadProof(
	target notary.Address,
	statement notary.Statement,
	signature notary.Signature,
	proof notary.InclusionProof,
) notary.SponsoredCall {
	siblings := make([]court.SumNode, len(proof.Siblings))
	for i, sibling := range proof.Siblings {
		siblings[i] = court.SumNode{Hash: sibling.Hash, Sum: sibling.Sum}
	}

	return notary.SponsoredCall{
		Target: target,
		Data:   c.court.PackProveBadProof(bound(statement), signature[:], siblings),
	}
}

func (c SponsoredCalls) ProveEquivocation(
	target notary.Address,
	first notary.Statement,
	firstSignature notary.Signature,
	second notary.Statement,
	secondSignature notary.Signature,
) notary.SponsoredCall {
	return notary.SponsoredCall{
		Target: target,
		Data:   c.court.PackProveEquivocation(bound(first), firstSignature[:], bound(second), secondSignature[:]),
	}
}

func registrationOf(issued notary.IssuedRegistration) withholding.Registration {
	return withholding.Registration{
		Holder:    common.Address(issued.Registration.Holder),
		Account:   issued.Registration.Account,
		FromEpoch: issued.Registration.FromEpoch,
	}
}
