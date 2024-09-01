package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
    "github.com/Nicktsim/kodetest/storage/psql"

)

type Claims struct {
	UserID   int    `json:"user_id"`
	Login string `json:"username"`
	jwt.RegisteredClaims
}

var secretKey = []byte(os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
func CreateToken(user *psql.User) (string,error) {
	const op = "utils.CreateToken"

	claims:= Claims{
		UserID: user.ID,
		Login: user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token :=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", fmt.Errorf("%s: %w", op, err)
    }

    return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
    const op = "utils.ValidateToken"

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    if !token.Valid {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return claims, nil
}