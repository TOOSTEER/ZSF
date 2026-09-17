package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"quest-backend/internal/auth/reqctx"
)

// Middleware validates the Authorization: Bearer <token> header and
// stores the parsed claims in the request context (via reqctx).
func Middleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing or malformed Authorization header")
				return
			}

			tokenString := strings.TrimPrefix(header, "Bearer ")
			claims, err := ParseToken(cfg, tokenString)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := reqctx.WithClaims(r.Context(), reqctx.Claims{
				UserID: claims.UserID,
				Role:   claims.Role,
				TeamID: claims.TeamID,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireHost ensures the caller's role is "host".
func RequireHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := reqctx.FromContext(r.Context())
		if !ok || claims.Role != "host" {
			writeError(w, http.StatusForbidden, "host role required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireMember ensures the caller's role is "member".
func RequireMember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := reqctx.FromContext(r.Context())
		if !ok || claims.Role != "member" {
			writeError(w, http.StatusForbidden, "member role required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
