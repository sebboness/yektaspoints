package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/jwt"
)

type ctxKey string

const (
	CtxKeyAuthInfo ctxKey = "api:auth"

	claimKeyUserID        = "sub"
	claimKeyUsername      = "cognito:username"
	claimKeyEmail         = "email"
	claimKeyEmailVerified = "email_verified"
	claimKeyName          = "name"
	claimKeyGroups        = "cognito:groups"
)

type AuthorizerInfo struct {
	Claims map[string]interface{} `json:"claims"`
}

type AuthContext interface {
	GetAuthorizerInfo(c *gin.Context) AuthorizerInfo
}

var jwtParser jwt.JwtParser

// GetAuthContext returns a new AuthContext
func GetAuthContext() (AuthContext, error) {
	if env.GetEnv("RUN_AS_WEB_API") == "true" {
		return NewApiAuthContext()
	} else {
		return NewLambdaAuthContext()
	}
}

func (i AuthorizerInfo) HasInfo() bool {
	return len(i.Claims) > 0
}

func (i AuthorizerInfo) ValueOrEmpty(key string) string {
	if i.HasInfo() {
		if _val, ok := i.Claims[key]; ok {
			return fmt.Sprintf("%v", _val)
		}
	}
	return ""
}

func (i AuthorizerInfo) IsEmailVerified() bool {
	emailVerifiedStr := i.ValueOrEmpty(claimKeyEmailVerified)
	return emailVerifiedStr == "true"
}

func (i AuthorizerInfo) GetEmail() string {
	return i.ValueOrEmpty(claimKeyEmail)
}

func (i AuthorizerInfo) GetGroups() []string {
	groupStr := i.ValueOrEmpty(claimKeyGroups)
	return strings.Split(groupStr, ",")
}

func (i AuthorizerInfo) GetName() string {
	return i.ValueOrEmpty(claimKeyName)
}

func (i AuthorizerInfo) GetUserID() string {
	return i.ValueOrEmpty(claimKeyUserID)
}

func (i AuthorizerInfo) GetUsername() string {
	return i.ValueOrEmpty(claimKeyUsername)
}
