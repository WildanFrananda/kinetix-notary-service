package httpapi

import (
	"encoding/json"
	"net/http"
)

const identityUnavailableMessage = "the notary cannot check tokens right now"

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	if status == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", `Bearer realm="kinetix"`)
	}

	writeJSON(w, status, map[string]string{"error": code, "message": message})
}
