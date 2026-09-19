// Package reqctx carries the authenticated caller's identity through the
// request context. It is intentionally dependency-free so that both the
// auth package (which sets it after validating a JWT) and every other
// domain package (which only needs to read it) can import it without
// creating an import cycle.
package reqctx

import "context"

// Claims is the minimal identity extracted from a valid JWT.
type Claims struct {
	UserID int
	Role   string // "host" | "member"
	TeamID int
}

type ctxKey string

const claimsKey ctxKey = "reqctx_claims"

// WithClaims returns a new context carrying the given claims.
func WithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

// FromContext extracts the claims stored by WithClaims.
func FromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey).(Claims)
	return c, ok
}
