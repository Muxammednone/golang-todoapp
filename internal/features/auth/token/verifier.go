package auth_token

import (
	"strconv"

	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	secret []byte
}

func NewVerifier(secret []byte) *Verifier {
	return &Verifier{
		secret: secret,
	}
}

func (v *Verifier) Verify(tokenString string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return v.secret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return 0, core_errors.ErrUnauthenticated
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, core_errors.ErrUnauthenticated
	}

	return userID, nil
}
