// util/auth.go
package util

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

// Claims defines the JWT payload.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT for a given user ID.
func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// VerifyJWT verifies the given JWT and returns the user ID.
func VerifyJWT(tokenString string, secretKey []byte) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New(err.Error())
	}
}
