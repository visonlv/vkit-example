package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	secret = []byte("halo~")
)

type Account struct {
	Id       string   `json:"id"`
	Roles    []string `json:"roles"`
	Metadata string   `json:"metadata"`
}

type AuthClaims struct {
	Id       string   `json:"id"`
	Roles    []string `json:"roles"`
	Metadata string   `json:"metadata"`

	jwt.StandardClaims
}

func GenToken(a *Account, td time.Duration) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		Id:    a.Id,
		Roles: a.Roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   a.Id,
			Issuer:    "infore",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(td).Unix(),
		},
	})

	tokenString, err := t.SignedString(secret)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func ParseToken(tokenString string) (*AuthClaims, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	tt, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tt.Claims.(*AuthClaims); ok && tt.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("token expired!")
	}
}
