package sponsorcli

type registrationsFile struct {
	Court         courtField          `json:"court"`
	Epoch         string              `json:"epoch"`
	Registrations []registrationField `json:"registrations"`
}
