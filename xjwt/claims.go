package xjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenClaims:
//
//  1. Issuer string `json:"iss,omitempty"`: issuer of the jwt
//
//  2. Subject string `json:"sub,omitempty"`: user id
//
//  3. (NA) Audience ClaimStrings `json:"aud,omitempty"`:
//
//     - access token: app ids the jwt can access
//
//     - id token: app id which the jwt is assigned to (id token is not supported in this wrapper)
//
//  4. ExpiresAt *jwt.NumericDate `json:"exp,omitempty"`: jwt expire time
//
//  5. IssuedAt *jwt.NumericDate `json:"iat,omitempty"`: jwt issue time
//
//  6. (NA) NotBefore *jwt.NumericDate `json:"nbf,omitempty"`: jwt invalid before the time
//
//  7. ID string `json:"jti,omitempty"`: jwt id
type TokenClaims struct {
	jwt.RegisteredClaims
}

func newTokenClaims(issuer string, subject string, lifetime time.Duration) *TokenClaims {
	now := time.Now()
	claims := &TokenClaims{}
	claims.Issuer = issuer
	claims.Subject = subject
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(lifetime))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ID = uuid.New().String()

	return claims
}
