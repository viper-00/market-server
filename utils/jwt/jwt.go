package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

var (
	simpleSecret = []byte("predict-market")
)

func CreateJWT(claims map[string]interface{}) (string, error) {
	allClaims := jwt.MapClaims{}
	if len(claims) > 0 {
		for key, value := range claims {
			allClaims[key] = value
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, allClaims)

	tokenString, err := token.SignedString(simpleSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return simpleSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("Claims not found")

}
