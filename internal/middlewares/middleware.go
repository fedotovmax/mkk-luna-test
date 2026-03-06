package middlewares

import "net/http"

type Middleware interface {
	Middleware(next http.Handler) http.Handler
}
