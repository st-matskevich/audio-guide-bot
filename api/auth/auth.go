package auth

import "time"

type TokenClaims struct {
	ExpiresAt time.Time
}

type TokenProvider interface {
	Create(claims TokenClaims) (string, error)
	Verify(token string) (TokenClaims, bool, error)
}
