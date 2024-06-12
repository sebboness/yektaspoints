package jwt

import (
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type JwtReader interface {
	// Parse returns the parsed token which includes its claims
	Parse(tokenString string, keyFunc jwtv5.Keyfunc) (*jwtv5.Token, error)
}

type Jwto struct {
	JwtReader
}

func NewJwt() JwtReader {
	return &Jwto{}
}

func (j *Jwto) Parse(tokenString string, keyFunc jwtv5.Keyfunc) (*jwtv5.Token, error) {
	return jwtv5.Parse(tokenString, keyFunc)
}
