package xjwt

import (
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokenParserConfig struct {
	Algorithm     Algorithm `json:",options=HS256|HS384|HS512|RS256|RS384|RS512|ES256|ES384|ES512"`
	SecretKey     string    `json:",optional"`
	SecretKeyPath string    `json:",default=./certs/key.pem.pub"`
}

func (c TokenParserConfig) Validate() error {
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

	return nil
}

type TokenParser struct {
	algorithm Algorithm
	secretKey any
}

func MustNewTokenParser(c TokenParserConfig) *TokenParser {
	// validate config
	err := c.Validate()
	if err != nil {
		logx.Errorf("MustNewTokenParser: %v", err)
		logx.Must(err)
	}

	tokenParser := &TokenParser{}
	tokenParser.algorithm = c.Algorithm
	secretKey, err := getPublicKey(c.Algorithm, c.SecretKey, c.SecretKeyPath)
	if err != nil {
		logx.Errorf("MustNewTokenParser: %v", err)
		logx.Must(err)
	}

	tokenParser.secretKey = secretKey

	return tokenParser
}

func (t *TokenParser) ValidateAndParseAccessToken(tokenString string) (*TokenClaims, error) {
	tokenToParse := getAccessToken(tokenString)

	// validate and parse access token
	return validateAndParseToken(t.algorithm, t.secretKey, tokenToParse)
}

func (t *TokenParser) ValidateAndParseRefreshToken(tokenString string) (*TokenClaims, error) {
	// validate and parse refresh token
	return validateAndParseToken(t.algorithm, t.secretKey, tokenString)
}

func (t *TokenParser) ParseAccessToken(tokenString string) (*TokenClaims, error) {
	tokenToParse := getAccessToken(tokenString)

	// parse access token
	return parseToken(t.algorithm, t.secretKey, tokenToParse)
}

func (t *TokenParser) ParseRefreshToken(tokenString string) (*TokenClaims, error) {
	// parse refresh token
	return parseToken(t.algorithm, t.secretKey, tokenString)
}
