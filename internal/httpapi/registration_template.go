package httpapi

type registrationTemplate struct {
	Types       map[string][]typedDataField `json:"types"`
	PrimaryType string                      `json:"primaryType"`
	Domain      templateDomain              `json:"domain"`
	Message     map[string]string           `json:"message"`
}
