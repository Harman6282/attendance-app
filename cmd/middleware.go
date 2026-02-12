package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/Harman6282/attendance-app/internal/store"
	"github.com/Harman6282/attendance-app/internal/token"
)

type ctxKey string

const claimsContextKey ctxKey = "auth_claims"

func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader == "" {
			app.writeJSONError(w, http.StatusUnauthorized, "authorization header required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			app.writeJSONError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := app.tokenMaker.VerifyToken(strings.TrimSpace(parts[1]))
		if err != nil {
			app.writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireRole(allowed ...store.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := getClaimsFromContext(r.Context())
			if !ok {
				app.writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			for _, role := range allowed {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			app.writeJSONError(w, http.StatusForbidden, "forbidden for this role")
		})
	}
}

func (app *application) teacherOnly(next http.Handler) http.Handler {
	return app.requireRole(store.Teacher)(next)
}

func (app *application) studentOnly(next http.Handler) http.Handler {
	return app.requireRole(store.Student)(next)
}

func getClaimsFromContext(ctx context.Context) (*token.JWTClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*token.JWTClaims)
	return claims, ok
}
