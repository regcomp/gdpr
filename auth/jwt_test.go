package auth

import (
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateThenValidateJWT(t *testing.T) {
	t.Parallel()

	privateKey, err := generateRSAKeyPair()
	if err != nil {
		// TODO:
	}

	tests := []struct {
		name  string
		input CustomClaims
	}{
		{
			"empty claims",
			CustomClaims{},
		},
		{
			"single claim",
			CustomClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer: "test",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := generateJWTWithClaims(&test.input, privateKey)
			if err != nil {
				t.Errorf("could not generate JWT with claims=%v, %s", test.input, err.Error())
			}

			claims, err := verifyJWTWithClaims(token, &CustomClaims{}, &privateKey.PublicKey)
			if err != nil {
				t.Errorf("could not verify jwt, %s", err.Error())
			}

			if !reflect.DeepEqual(claims, &test.input) {
				t.Errorf(
					"expected=%+v, got=%+v",
					test.input,
					claims,
				)
			}
		})
	}
}
