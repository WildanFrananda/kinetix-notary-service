package ignition

type Ceremony struct {
	Name                string
	Points              int
	PointsPerTranscript int
	Transcripts         int
}

func MainIgnition() Ceremony {
	return Ceremony{
		Name:                "MAIN IGNITION",
		Points:              100_800_000,
		PointsPerTranscript: 5_040_000,
		Transcripts:         20,
	}
}
