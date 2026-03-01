package httpcommon

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeJSON_Success(t *testing.T) {
	input := map[string]string{"name": "Maxim", "surname": "Fedotov"}
	body, _ := json.Marshal(input)

	var out map[string]string
	err := DecodeJSON(bytes.NewReader(body), &out)
	require.NoError(t, err)
	require.Equal(t, input, out)
}

func TestDecodeJSON_InvalidJSON(t *testing.T) {
	body := bytes.NewReader([]byte(`{invalid json data...}`))

	var out map[string]string
	err := DecodeJSON(body, &out)
	require.Error(t, err)
}

func TestWriteJSON_Success(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"status": "ok"}

	WriteJSON(rec, http.StatusOK, data)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, ContentTypeJSON, res.Header.Get(HeaderContentType))

	var decoded map[string]string
	err := json.NewDecoder(res.Body).Decode(&decoded)
	require.NoError(t, err)
	require.Equal(t, data, decoded)
}

func TestWriteJSON_EncodeError(t *testing.T) {
	rec := httptest.NewRecorder()

	data := func() {}

	WriteJSON(rec, http.StatusOK, data)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var msg map[string]string
	body, _ := io.ReadAll(res.Body)
	_ = json.Unmarshal(body, &msg)
	require.Equal(t, "failed to encode json", msg["message"])
}

func TestMessage(t *testing.T) {
	m := "hello"
	res := Message(m)
	require.Equal(t, m, res.Message)
}
