package keccak

import (
	"encoding/binary"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

const (
	CourtName    = "Kinetix Statement Court"
	CourtVersion = "1"

	domainType    = "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
	statementType = "EpochStatement(address ledger,uint64 epoch,bytes32 commitment,uint64 amount,uint64 index,bytes32 siblingsHash)"
)

var _ notary.StatementScheme = Scheme{}

func (Scheme) SiblingsHash(siblings []notary.Node) notary.Digest {
	packed := make([]byte, 0, len(siblings)*(32+8))
	for _, sibling := range siblings {
		packed = append(packed, sibling.Hash[:]...)
		packed = binary.BigEndian.AppendUint64(packed, sibling.Sum)
	}

	return digest(packed)
}

func (Scheme) StatementDigest(court notary.CourtIdentity, statement notary.Statement) notary.Digest {
	structure := words(
		digest([]byte(statementType)),
		addressWord(statement.Ledger),
		uintWord(statement.Epoch),
		statement.Commitment,
		uintWord(statement.Amount),
		uintWord(statement.Index),
		statement.SiblingsHash,
	)

	return typedData(CourtName, CourtVersion, court, structure)
}

func typedData(name, version string, court notary.CourtIdentity, structure notary.Digest) notary.Digest {
	domain := words(
		digest([]byte(domainType)),
		digest([]byte(name)),
		digest([]byte(version)),
		uintWord(court.ChainID),
		addressWord(court.Address),
	)

	packed := make([]byte, 0, 2+32+32)
	packed = append(packed, 0x19, 0x01)
	packed = append(packed, domain[:]...)
	packed = append(packed, structure[:]...)

	return digest(packed)
}

func words(parts ...notary.Digest) notary.Digest {
	packed := make([]byte, 0, len(parts)*32)
	for _, part := range parts {
		packed = append(packed, part[:]...)
	}

	return digest(packed)
}

func uintWord(value uint64) notary.Digest {
	var word notary.Digest
	binary.BigEndian.PutUint64(word[24:], value)

	return word
}

func addressWord(address notary.Address) notary.Digest {
	var word notary.Digest
	copy(word[12:], address[:])

	return word
}
