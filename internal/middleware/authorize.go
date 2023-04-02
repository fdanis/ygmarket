package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fdanis/yg-loyalsys/internal/common"
	"github.com/golang-jwt/jwt/v4"
)

type AuthorizeMiddleware struct {
	secretKey string
}

func NewAuthorizeMiddleware(key string) (a *AuthorizeMiddleware) {
	return &AuthorizeMiddleware{secretKey: key}
}

func (a *AuthorizeMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w)
			return
		} else {
			jwtToken := authHeader
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(a.secretKey), nil
			})
			if err != nil {
				unauthorized(w)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				err := claims.Valid()
				if err != nil {
					unauthorized(w)
					return
				}
				c := common.UserClaims{
					Login: claims["Login"].(string),
					ID:    int(claims["ID"].(float64)),
				}
				ctx := context.WithValue(r.Context(), common.Auth, c)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				unauthorized(w)
				return
			}
		}
	})
}

func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized"))
}
