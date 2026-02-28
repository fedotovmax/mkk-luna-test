package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/httpcommon"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
)

var emptyHeaderErr = errors.New("empty authorization header")
var badHeaderFormatErr = errors.New("bad authorization header format")

var Unauthorized = "Unauthorized"

func validateAuthHeader(header string) (string, error) {
	if header == "" {
		return "", emptyHeaderErr
	}

	authHeaderParts := strings.Split(header, " ")

	if len(authHeaderParts) != 2 {
		return "", badHeaderFormatErr
	}

	if authHeaderParts[0] != httpcommon.HeaderAuthorizationBearer {
		return "", badHeaderFormatErr
	}

	return authHeaderParts[1], nil
}

func NewAuthMiddleware(
	log *slog.Logger,
	tokenManager ports.TokenManager,
	tokenSecret string,
	issuer string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get(httpcommon.HeaderAuthorization)

			accessToken, err := validateAuthHeader(authHeader)

			if err != nil {
				log.Error("auth failed middleware failed", logger.Err(err))
				httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(Unauthorized))
				return
			}

			sid, uid, err := tokenManager.Verify(accessToken, issuer, tokenSecret)

			if err != nil {
				log.Error("auth failed middleware failed", logger.Err(err))
				httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(Unauthorized))
				return
			}

			ctx := context.WithValue(r.Context(), httpcommon.SessionCtxKey, &domain.Local{UserID: uid, SessionID: sid})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
