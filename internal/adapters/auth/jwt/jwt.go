package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type manager struct {
	secret  string
	expires time.Duration
}

func New(secret string, expires time.Duration) *manager {
	return &manager{secret: secret, expires: expires}
}

func (m *manager) Create(issuer, uid, sid string) (token string, exp time.Time, err error) {

	const op = "adapters.auth.jwt.Create"

	now := time.Now().UTC()

	exp = now.Add(m.expires)

	accessClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    issuer,
		ID:        sid,
		Subject:   uid,
	}

	accessTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	token, err = accessTokenObject.SignedString([]byte(m.secret))

	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		exp = time.Time{}
		return "", exp, err
	}

	return token, exp, nil
}

func (m *manager) Verify(token, issuer, secret string) (jti string, uid string, err error) {

	const op = "adapters.auth.jwt.Verify"

	opts := []jwt.ParserOption{
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	}

	result, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, opts...)

	if err != nil {
		err = fmt.Errorf("%s: %w: %v", op, ErrParseClaims, err)
		return jti, uid, err
	}

	if !result.Valid {
		err = fmt.Errorf("%s: %w", op, ErrInvalidToken)
		return jti, uid, err
	}

	claims, ok := result.Claims.(*jwt.RegisteredClaims)

	if !ok {
		err = fmt.Errorf("%s: %w", op, ErrParseClaims)
		return jti, uid, err
	}

	if claims.Subject == "" || claims.ID == "" {
		err = fmt.Errorf("%s: %w", op, ErrInvalidToken)
		return "", "", err
	}

	jti, uid = claims.ID, claims.Subject

	return jti, uid, nil
}
