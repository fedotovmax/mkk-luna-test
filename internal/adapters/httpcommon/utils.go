package httpcommon

import (
	"bytes"
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

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, `{"message": "failed to encode json"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

type MessageResponse struct {
	Message string `json:"message"`
}

func Message(m string) MessageResponse {
	return MessageResponse{Message: m}
}
