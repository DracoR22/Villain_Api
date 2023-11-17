package utils

import (
	"fmt"
	"os"

	"github.com/DracoR22/villain_api/app/types"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(account *types.Account) (string, error) {
	claims := &jwt.MapClaims{
		"expires_at":     15000,
		"account_number": account.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}
