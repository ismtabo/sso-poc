package middleware

import (
	"net/http"

	"github.com/ismtabo/sso-poc/auth/pkg/context"
	"github.com/ismtabo/sso-poc/auth/pkg/service"
)

func WithSession(sessions service.SessionService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.WithSessionContext(r.Context())
		isAuthenticated, _ := sessions.IsAuthenticated(rw, r)
		session := context.GetContextSession(ctx)
		session.Authenticated = isAuthenticated
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
