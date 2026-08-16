package jwt_pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {

	testCase := []struct {
		name        string
		claimsInput jwt.MapClaims
		generateKey func() (*rsa.PrivateKey, *rsa.PublicKey)
		verifyFunc  func(t *testing.T, claims jwt.MapClaims, err error)
	}{
		{
			name:        "normal case",
			claimsInput: jwt.MapClaims{"test": "test"},
			generateKey: func() (*rsa.PrivateKey, *rsa.PublicKey) {
				privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
				publicKey := privateKey.Public().(*rsa.PublicKey)
				return privateKey, publicKey
			},
			verifyFunc: func(t *testing.T, claimsOutput jwt.MapClaims, err error) {
				assert.NotNil(t, claimsOutput)
				assert.Nil(t, err)
			},
		},
		{
			name:        "err case - invalid token",
			claimsInput: jwt.MapClaims{"test": "test"},
			generateKey: func() (*rsa.PrivateKey, *rsa.PublicKey) {
				privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
				publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte("AllYourBase"))
				if err != nil {
					t.Fatalf("Failed to parse public key: %v", err)
				}
				return privateKey, publicKey
			},
			verifyFunc: func(t *testing.T, claimsOutput jwt.MapClaims, err error) {
				assert.Nil(t, claimsOutput)
				assert.Equal(t, err, InvalidToken)
			},
		},
	}

	for _, tc := range testCase {
		t.Parallel()
		privateKey, publicKey := tc.generateKey()
		jwtGenerator := &jwtGenerator{privateKey: privateKey}
		tokenString, _ := jwtGenerator.GenerateToken(tc.claimsInput)
		jwtValidator := &jwtValidator{publicKey: publicKey}
		claims, err := jwtValidator.ValidateToken(tokenString)
		tc.verifyFunc(t, claims, err)
	}

}
