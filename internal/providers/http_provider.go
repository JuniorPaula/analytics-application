package providers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ErrorApi struct {
	Error string `json:"error"`
}

// WriteJSON is a helper function to write JSON responses
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Fatal(err)
		}
	}
}

// VerifyStatusCodeErrors is a helper function to write JSON responses
func VerifyStatusCodeErrors(w http.ResponseWriter, r *http.Response) {
	var err ErrorApi
	json.NewDecoder(r.Body).Decode(&err)
	WriteJSON(w, r.StatusCode, err)
}

func MakeRquest(r *http.Request, method, url string, data io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
