package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// NOTE: This struct may end up being a giant amalgamation of fields or other
// embedded structs to hold differing claims fields from different provider
// implementations
type CustomClaims struct {
	jwt.RegisteredClaims
}

func generateRSAKeyPair() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		// TODO:
	}
	return privateKey, nil
}

func generateJWTWithClaims(claims *CustomClaims, privateKey *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		// TODO:
	}

	return tokenString, nil
}

func verifyJWTWithClaims(
	tokenString string,
	customClaims *CustomClaims,
	publicKey *rsa.PublicKey,
) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, customClaims, func(token *jwt.Token) (interface{}, error) {
		// NOTE: the signing method may need to be parameterized somehow
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method=%v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func VerifyAccessToken(token string) bool {
	// TODO:
	return true
}
