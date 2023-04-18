package common

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	Login string
	ID    int
	jwt.RegisteredClaims
}
