package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/jwt"
)

type ApiAuthContext struct {
	AuthContext
	JwtParser jwt.JwtParser
}

func NewApiAuthContext() (AuthContext, error) {
	region := env.GetEnv("AWS_REGION")
	userPoolId := env.GetEnv("COGNITO_USER_POOL_ID")

	jwtParser, err := jwt.NewAwsJwtParser(region, userPoolId)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize jwt parser: %v", err.Error())
	}

	return &ApiAuthContext{
		JwtParser: jwtParser,
	}, nil
}

func (a *ApiAuthContext) GetAuthorizerInfo(c *gin.Context) AuthorizerInfo {
	info := AuthorizerInfo{}

	if c.Request == nil {
		return info
	}

	ctx := c.Request.Context()
	_info := ctx.Value(CtxKeyAuthInfo)

	// return info if already set
	if _info != nil {
		return _info.(AuthorizerInfo)
	}

	// if we got here, check if token is in authorization header (if running web api)
	reqToken := a.getTokenFromHeader(c.Request)
	if reqToken == "" {
		return info
	}

	// parse claims
	claims, err := a.JwtParser.GetJwtClaims(reqToken)
	if err != nil {
		logger.Errorf("failed to get jwt claims: %v", err.Error())
		return info
	}

	logger.Infof("claims = %v", claims)
	info = AuthorizerInfo{
		Claims: claims,
	}

	// Store claims in context
	c.Request = c.Request.Clone(context.WithValue(ctx, CtxKeyAuthInfo, info))

	logger.Infof("authorized info = %v", info)

	return info
}

// getTokenFromHeader returns the bearer token from the Authorization header
func (a *ApiAuthContext) getTokenFromHeader(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return ""
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	return splitToken[1]
}
