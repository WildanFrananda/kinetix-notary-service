package httpapi

type registrationResponse struct {
	Court         courtResponse        `json:"court"`
	Epoch         string               `json:"epoch"`
	Registrations []issuedRegistration `json:"registrations"`
}
