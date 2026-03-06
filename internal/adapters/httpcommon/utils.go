package httpcommon

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
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

var ErrUnauthorized = errors.New("unauthorized")

func GetLocalSession(r *http.Request) (*domain.Local, error) {
	user, ok := r.Context().Value(SessionCtxKey).(*domain.Local)
	if !ok {
		return nil, ErrUnauthorized
	}
	return user, nil
}
