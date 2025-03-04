package xjwt

import (
	"log"
)

type TokenParserConf struct {
	Algorithm     Algorithm `json:",options=HS256|HS384|HS512|RS256|RS384|RS512|ES256|ES384|ES512"`
	SecretKey     string    `json:",optional"`
	SecretKeyPath string    `json:",default=./certs/key.pem.pub"`
}

type TokenParser struct {
	algorithm Algorithm
	secretKey interface{}
}

func MustNewTokenParser(c TokenParserConf) *TokenParser {
	tokenParser := &TokenParser{}
	tokenParser.algorithm = c.Algorithm
	secretKey, err := getPublicKey(c.Algorithm, c.SecretKey, c.SecretKeyPath)
	if err != nil {
		log.Fatalf("error: MustNewTokenParser: %s", err.Error())
	}

	tokenParser.secretKey = secretKey

	return tokenParser
}

func (t *TokenParser) ParseAccessTokenForAuth(tokenString string) (*TokenClaims, error) {
	tokenToParse := getAccessToken(tokenString)

	// parse access token
	return parseToken(t.algorithm, t.secretKey, tokenToParse, false)
}

func (t *TokenParser) ParseTokensForRefresh(accessTokenString string, refreshTokenString string) (*TokenClaims, error) {
	accessTokenToParse := getAccessToken(accessTokenString)

	// parse access token
	accessTokenClaims, err := parseToken(t.algorithm, t.secretKey, accessTokenToParse, true)
	// access token can be expired
	if err != nil {
		return nil, err
	}

	// parse refresh token
	_, err = parseToken(t.algorithm, t.secretKey, refreshTokenString, false)
	if err != nil {
		return nil, err
	}

	return accessTokenClaims, nil
}
