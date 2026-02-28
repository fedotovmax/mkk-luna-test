package httpcommon

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON(r io.Reader, v any) error {
	defer io.Copy(io.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, `{"message": "failed to encode json"}`, http.StatusInternalServerError)
	}
}

type MessageResponse struct {
	Message string `json:"message"`
}

func Message(m string) MessageResponse {
	return MessageResponse{Message: m}
}
