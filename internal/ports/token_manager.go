package ports

import "time"

type TokenManager interface {
	Create(issuer, uid, sid string) (token string, exp time.Time, err error)
	Verify(token, issuer, secret string) (jti string, uid string, err error)
}
