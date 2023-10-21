package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SIGN_METHOD = jwt.SigningMethodHS256

type JWTTokenProvider struct {
	JWTSecret []byte
}

func (provider *JWTTokenProvider) Create(claims TokenClaims) (string, error) {
	jwtClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
	}

	token := jwt.NewWithClaims(JWT_SIGN_METHOD, jwtClaims)
	tokenString, err := token.SignedString(provider.JWTSecret)

	return tokenString, err
}

func (provider *JWTTokenProvider) Verify(token string) (TokenClaims, bool, error) {
	jwtClaims := jwt.RegisteredClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &jwtClaims, func(token *jwt.Token) (any, error) {
		if token.Method != JWT_SIGN_METHOD {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(provider.JWTSecret), nil
	})

	if err != nil {
		return TokenClaims{}, false, err
	}

	if !jwtToken.Valid {
		return TokenClaims{}, false, nil
	}

	if jwtClaims.ExpiresAt.Time.Before(time.Now()) {
		return TokenClaims{}, false, nil
	}

	result := TokenClaims{
		ExpiresAt: jwtClaims.ExpiresAt.Time,
	}

	return result, true, nil
}

func CreateJWTTokenProvider(secret string) (TokenProvider, error) {
	if secret == "" {
		return nil, errors.New("secret is empty")
	}

	provider := JWTTokenProvider{
		JWTSecret: []byte(secret),
	}

	return &provider, nil
}
