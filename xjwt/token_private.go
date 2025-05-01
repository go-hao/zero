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
	if len(tokenSplit) >= 2 {
		return tokenSplit[1]
	}

	return tokenString
}

func parseToken(algorithm Algorithm, secretKey any, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{},
		func(token *jwt.Token) (any, error) {
			return secretKey, nil
		},
		jwt.WithValidMethods([]string{string(algorithm)}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func createToken(issuer string, subject string, lifetimeInSec int64, algorithm Algorithm, secretKey any) (*Token, error) {
	// init token
	token := jwt.New(jwt.GetSigningMethod(string(algorithm)))

	// init claims
	claims := newTokenClaims(issuer, subject, lifetimeInSec)

	// add cliams to token
	token.Claims = claims

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:     tokenString,
		ExpiresIn: lifetimeInSec,
	}, nil
}
