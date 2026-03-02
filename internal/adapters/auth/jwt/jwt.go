package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type manager struct {
	secret string
}

func New(secret string) *manager {
	return &manager{secret: secret}
}

type CreateParams struct {
	Issuer         string
	Uid            string
	Sid            string
	TokenExpiresAt time.Time
	Now            time.Time
}

func (m *manager) Create(p *CreateParams) (token string, err error) {

	const op = "adapters.auth.jwt.Create"

	accessClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(p.TokenExpiresAt),
		IssuedAt:  jwt.NewNumericDate(p.Now),
		Issuer:    p.Issuer,
		ID:        p.Sid,
		Subject:   p.Uid,
	}

	accessTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	token, err = accessTokenObject.SignedString([]byte(m.secret))

	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)

		return "", err
	}

	return token, nil
}

func (m *manager) Verify(token, issuer string) (jti string, uid string, err error) {

	const op = "adapters.auth.jwt.Verify"

	opts := []jwt.ParserOption{
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	}

	result, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(m.secret), nil
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
