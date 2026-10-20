package proverproc

type wireRequest struct {
	Number          string      `json:"number"`
	Depth           int         `json:"depth"`
	TreeRoot        string      `json:"treeRoot"`
	TotalCommitment string      `json:"totalCommitment"`
	Blinding        string      `json:"blinding"`
	Entries         []wireEntry `json:"entries"`
}
