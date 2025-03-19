package xjwt

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type TokenCreatorConf struct {
	Algorithm                 Algorithm `json:",default=HS256"`
	SecretKey                 string    `json:",optional"`
	SecretKeyPath             string    `json:",optional"`
	AccessTokenLifetimeInSec  int64     `json:",default=3600"`
	RefreshTokenLifetimeInSec int64     `json:",default=86400"`
}

func (c TokenCreatorConf) Validate() error {
	switch c.Algorithm {
	case RS256, RS384, RS512, ES256, ES384, ES512:
		if len(c.SecretKeyPath) == 0 {
			return errors.New(fmt.Sprintf("SecretKeyPath cannot be empty, when Algorithm is %s", c.Algorithm))
		}
	case HS256, HS384, HS512:
		if len(c.SecretKey) == 0 && len(c.SecretKeyPath) == 0 {
			return errors.New(fmt.Sprintf("either SecretKeyPath or SecretKey cannot be empty, when Algorithm is %s", c.Algorithm))
		}
	default:
		return errors.New("Algorithm must be one of HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512")
	}
	if c.AccessTokenLifetimeInSec <= 0 || c.AccessTokenLifetimeInSec > 86400 {
		return errors.New("AccessTokenLifetimeInSec must be in range (0, 86400]")
	}
	if c.RefreshTokenLifetimeInSec <= 0 || c.RefreshTokenLifetimeInSec > 31536000 {
		return errors.New("RefreshTokenLifetimeInSec must be in range (0, 31536000]")
	}

	return nil
}

type TokenCreator struct {
	algorithm                 Algorithm
	secretKey                 any
	AccessTokenLifetimeInSec  int64
	RefreshTokenLifetimeInSec int64
}

type Token struct {
	Value     string
	ExpiresIn int64
}

func MustNewTokenCreator(c TokenCreatorConf) *TokenCreator {
	tokenCreator := &TokenCreator{}
	tokenCreator.algorithm = c.Algorithm
	secretKey, err := getPrivateKey(c.Algorithm, c.SecretKey, c.SecretKeyPath)
	if err != nil {
		log.Fatalf("error: MustNewTokenCreator: %s", err.Error())
	}

	tokenCreator.secretKey = secretKey
	tokenCreator.AccessTokenLifetimeInSec = c.AccessTokenLifetimeInSec
	tokenCreator.RefreshTokenLifetimeInSec = c.RefreshTokenLifetimeInSec

	return tokenCreator
}

func (t *TokenCreator) CreateAccessToken(issuer string, subject string) (*Token, error) {
	// init token
	token := jwt.New(jwt.GetSigningMethod(string(t.algorithm)))

	// init claims
	claims := newTokenClaims(issuer, subject, t.AccessTokenLifetimeInSec)

	// add cliams to token
	token.Claims = claims

	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:     tokenString,
		ExpiresIn: t.AccessTokenLifetimeInSec,
	}, nil
}

func (t *TokenCreator) CreateRefreshToken(issuer string) (*Token, error) {
	// init token
	token := jwt.New(jwt.GetSigningMethod(string(t.algorithm)))

	// init claims
	claims := newTokenClaims(issuer, "", t.RefreshTokenLifetimeInSec)

	// add cliams to token
	token.Claims = claims

	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:     tokenString,
		ExpiresIn: t.RefreshTokenLifetimeInSec,
	}, nil
}
