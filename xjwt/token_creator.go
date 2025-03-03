package xjwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenCreatorConf struct {
	Algorithm                 Algorithm `json:",options=HS256|HS384|HS512|RS256|RS384|RS512|ES256|ES384|ES512"`
	SecretKey                 string    `json:",optional"`
	SecretKeyPath             string    `json:",default=./certs/key.pem"`
	AccessTokenLifetimeInSec  int64     `json:",range=(0:86400],default=3600"`
	RefreshTokenLifetimeInSec int64     `json:",range=(0:31536000],default=3600"`
}

type TokenCreator struct {
	algorithm            Algorithm
	secretKey            interface{}
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

func MustNewTokenCreator(c TokenCreatorConf) *TokenCreator {
	tokenCreator := &TokenCreator{}
	tokenCreator.algorithm = c.Algorithm
	secretKey, err := getPrivateKey(c.Algorithm, c.SecretKey, c.SecretKeyPath)
	if err != nil {
		log.Fatalf("error: MustNewTokenCreator: %s", err.Error())
	}
	tokenCreator.secretKey = secretKey
	tokenCreator.AccessTokenLifetime = time.Duration(c.AccessTokenLifetimeInSec) * time.Second
	tokenCreator.RefreshTokenLifetime = time.Duration(c.RefreshTokenLifetimeInSec) * time.Second

	return tokenCreator
}

func (t *TokenCreator) CreateAccessToken(issuer string, subject string) (string, error) {
	// init token
	token := jwt.New(jwt.GetSigningMethod(string(t.algorithm)))

	// init claims
	claims := newTokenClaims(issuer, subject, t.AccessTokenLifetime)

	// add cliams to token
	token.Claims = claims

	return token.SignedString(t.secretKey)
}

func (t *TokenCreator) CreateRefreshToken(issuer string) (string, error) {
	// init token
	token := jwt.New(jwt.GetSigningMethod(string(t.algorithm)))

	// init claims
	claims := newTokenClaims(issuer, "", t.RefreshTokenLifetime)

	// add cliams to token
	token.Claims = claims

	return token.SignedString(t.secretKey)
}
