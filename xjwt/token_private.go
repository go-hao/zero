package xjwt

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func getPrivateKey(algorithm Algorithm, secretKey string, secretKeyPath string) (any, error) {
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

func getPublicKey(algorithm Algorithm, secretKey string, secretKeyPath string) (any, error) {
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

func parseToken(algorithm Algorithm, secretKey any, tokenString string, canExpire bool) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{},
		func(token *jwt.Token) (any, error) {
			return secretKey, nil
		},
		jwt.WithValidMethods([]string{string(algorithm)}),
		jwt.WithExpirationRequired(),
	)

	switch err {
	case jwt.ErrTokenExpired:
		if !canExpire {
			return nil, ErrTokenExpired
		}
	case nil:
		if !canExpire && !token.Valid {
			return nil, ErrInvalidToken
		}
	default:
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
