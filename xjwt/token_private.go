package xjwt

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func getPrivateKey(algorithm Algorithm, secretKey string, secretKeyPath string) (interface{}, error) {
	switch algorithm {
	case RS256, RS384, RS512:
		key, err := os.ReadFile(secretKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPrivateKeyFromPEM(key)
	case ES256, ES384, ES512:
		key, err := os.ReadFile(secretKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseECPrivateKeyFromPEM(key)
	default:
		if len(secretKey) > 0 {
			return []byte(secretKey), nil
		}
		return os.ReadFile(secretKeyPath)
	}
}

func getPublicKey(algorithm Algorithm, secretKey string, secretKeyPath string) (interface{}, error) {
	switch algorithm {
	case RS256, RS384, RS512:
		key, err := os.ReadFile(secretKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(key)
	case ES256, ES384, ES512:
		key, err := os.ReadFile(secretKeyPath)
		if err != nil {
			return nil, err
		}
		return jwt.ParseECPublicKeyFromPEM(key)
	default:
		if len(secretKey) > 0 {
			return []byte(secretKey), nil
		}
		return os.ReadFile(secretKeyPath)
	}
}

func getAccessToken(tokenString string) string {
	tokenSplit := strings.Split(tokenString, " ")
	if len(tokenSplit) == 2 {
		return tokenSplit[1]
	}

	return tokenString
}

func parseToken(algorithm Algorithm, secretKey interface{}, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		jwt.WithValidMethods([]string{string(algorithm)}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
