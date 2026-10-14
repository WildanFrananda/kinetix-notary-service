package keccak

import (
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type domain struct {
	name    string
	version string
}

var (
	courtDomain       = domain{name: CourtName, version: CourtVersion}
	withholdingDomain = domain{name: WithholdingCourtName, version: WithholdingCourtVersion}
	councilDomain     = domain{name: CouncilName, version: CouncilVersion}
)

func typedData(of domain, chainID uint64, contract notary.Address, structure notary.Digest) notary.Digest {
	separator := words(
		digest([]byte(domainType)),
		digest([]byte(of.name)),
		digest([]byte(of.version)),
		uintWord(chainID),
		addressWord(contract),
	)

	packed := make([]byte, 0, 2+32+32)
	packed = append(packed, 0x19, 0x01)
	packed = append(packed, separator[:]...)
	packed = append(packed, structure[:]...)

	return digest(packed)
}
