package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/util/jwt"
)

type ApiAuthContext struct {
	AuthContext
	JwtParser jwt.JwtParser
}

func NewApiAuthContext() (AuthContext, error) {
	jwtParser, err := jwt.NewJwtParser()
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

	// if we got here, check if token is in authorization header (if running web api)
	reqToken := a.getTokenFromHeader(c.Request)

	if jwtParser == nil {
		_jwtParser, err := jwt.NewJwtParser()
		if err != nil {
			logger.Errorf("failed to initialize jwt parser: %v", err.Error())
			return AuthorizerInfo{}
		}

		jwtParser = _jwtParser
	}

	claims, err := jwtParser.GetJwtClaims(reqToken)
	if err != nil {
		logger.Errorf("failed to initialize jwt parser: %v", err.Error())
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
	splitToken := strings.Split(reqToken, "Bearer ")
	return splitToken[1]
}
