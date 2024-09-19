package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleResponseError(w http.ResponseWriter, message string, responseCode int) {
	log.Println(message)
	w.WriteHeader(responseCode)

	if message == "" {
		return
	}

	resp := map[string]string{"error": message}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}
