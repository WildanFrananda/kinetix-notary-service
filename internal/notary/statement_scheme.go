package notary

type StatementScheme interface {
	SiblingsHash(siblings []Node) Digest
	StatementDigest(court CourtIdentity, statement Statement) Digest
}
