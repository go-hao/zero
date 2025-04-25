package xjwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is not valid")
	ErrTokenExpired = jwt.ErrTokenExpired
)
