package xjwt

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type TokenCreatorConf struct {
	Algorithm                 Algorithm `json:",options=HS256|HS384|HS512|RS256|RS384|RS512|ES256|ES384|ES512"`
	SecretKey                 string    `json:",optional"`
	SecretKeyPath             string    `json:",default=./certs/key.pem"`
	AccessTokenLifetimeInSec  int64     `json:",range=(0:86400],default=3600"`
	RefreshTokenLifetimeInSec int64     `json:",range=(0:31536000],default=86400"`
}

type TokenCreator struct {
	algorithm                 Algorithm
	secretKey                 interface{}
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
