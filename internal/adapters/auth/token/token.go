package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type CreateTokenResult struct {
	Nohashed string
	Hashed   string
}

func CreateToken() (*CreateTokenResult, error) {
	token, err := GenerateToken()

	if err != nil {
		return nil, err
	}

	hash := HashToken(token)

	resulst := &CreateTokenResult{
		Nohashed: token,
		Hashed:   hash,
	}

	return resulst, nil
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
