package keccak

import (
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	CouncilName    = "Kinetix Anchor Council"
	CouncilVersion = "1"

	approvalType = "EpochApproval(address ledger,uint64 number,bytes32 digest)"
)

var _ notary.ApprovalScheme = Scheme{}

func (Scheme) ApprovalDigest(
	council notary.CouncilIdentity,
	ledger notary.Address,
	number uint64,
	epoch notary.Digest,
) notary.Digest {
	structure := words(
		digest([]byte(approvalType)),
		addressWord(ledger),
		uintWord(number),
		epoch,
	)

	return typedData(councilDomain, council.ChainID, council.Address, structure)
}
