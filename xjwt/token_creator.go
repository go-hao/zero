package xjwt

import (
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokenCreatorConfig struct {
	Algorithm                 Algorithm `json:",options=HS256|HS384|HS512|RS256|RS384|RS512|ES256|ES384|ES512"`
	SecretKey                 string    `json:",optional"`
	SecretKeyPath             string    `json:",default=./certs/key.pem"`
	AccessTokenLifetimeInSec  int64     `json:",range=(0:86400],default=3600"`
	RefreshTokenLifetimeInSec int64     `json:",range=(0:31536000],default=86400"`
}

func (c TokenCreatorConfig) Validate() error {
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

func MustNewTokenCreator(c TokenCreatorConfig) *TokenCreator {
	// validate config
	err := c.Validate()
	if err != nil {
		logx.Errorf("MustNewTokenCreator: %v", err)
		logx.Must(err)
	}

	tokenCreator := &TokenCreator{}
	tokenCreator.algorithm = c.Algorithm
	secretKey, err := getPrivateKey(c.Algorithm, c.SecretKey, c.SecretKeyPath)
	if err != nil {
		logx.Errorf("MustNewTokenCreator: %v", err)
		logx.Must(err)
	}

	tokenCreator.secretKey = secretKey
	tokenCreator.AccessTokenLifetimeInSec = c.AccessTokenLifetimeInSec
	tokenCreator.RefreshTokenLifetimeInSec = c.RefreshTokenLifetimeInSec

	return tokenCreator
}

func (t *TokenCreator) CreateAccessToken(issuer string, subject string) (*Token, error) {
	return createToken(issuer, subject, t.AccessTokenLifetimeInSec, t.algorithm, t.secretKey)
}

func (t *TokenCreator) CreateRefreshToken(issuer string, subject string) (*Token, error) {
	return createToken(issuer, subject, t.RefreshTokenLifetimeInSec, t.algorithm, t.secretKey)
}
