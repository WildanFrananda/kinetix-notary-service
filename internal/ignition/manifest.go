package ignition

type Manifest struct {
	Name                string          `json:"name"`
	NumG1Points         int             `json:"numG1Points"`
	NumG2Points         int             `json:"numG2Points"`
	PointsPerTranscript int             `json:"pointsPerTranscript"`
	Participants        []Participant   `json:"participants"`
	Crs                 ReferenceString `json:"crs"`
}
