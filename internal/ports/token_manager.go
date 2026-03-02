package ports

import "github.com/fedotovmax/mkk-luna-test/internal/adapters/auth/jwt"

type TokenManager interface {
	Create(p *jwt.CreateParams) (token string, err error)
	Verify(token, issuer string) (jti string, uid string, err error)
}
