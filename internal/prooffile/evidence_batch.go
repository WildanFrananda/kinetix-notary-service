package prooffile

type evidenceBatch struct {
	Stream             string `json:"stream"`
	Number             string `json:"number"`
	PreviousDigest     string `json:"previousDigest"`
	Root               string `json:"root"`
	LeafCount          string `json:"leafCount"`
	ThroughTransaction string `json:"throughTransaction"`
	ThroughSequence    string `json:"throughSequence"`
	ClosedAt           string `json:"closedAt"`
	Digest             string `json:"digest"`
}
