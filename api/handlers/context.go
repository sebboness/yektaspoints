package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
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

func PrepareAuthorizedContext(ctx context.Context, req events.APIGatewayProxyRequest) context.Context {
	authorizer := AuthorizerInfo{}

	if len(req.RequestContext.Authorizer) > 0 {
		if claimsObj, ok := req.RequestContext.Authorizer["claims"]; ok {
			if claims, ok := claimsObj.(map[string]interface{}); ok {
				authorizer.Claims = claims
			}
		}
	} else {
		// logger.Infof("PrepareAuthorizedContext B")
	}

	return context.WithValue(ctx, CtxKeyAuthInfo, authorizer)
}

func GetAuthorizerInfo(c *gin.Context) AuthorizerInfo {
	info := AuthorizerInfo{}

	if c.Request == nil {
		return info
	}

	ctx := c.Request.Context()
	_info := ctx.Value(CtxKeyAuthInfo)

	authInfoJson, _ := json.Marshal(info)
	logger.Infof("GetAuthorizerInfo authInfo?: " + string(authInfoJson))

	if _info != nil {
		return _info.(AuthorizerInfo)
	}

	// if we got here, check if token is in authorization header (if running web api)
	// if env.GetEnv("RUN_AS_WEB_API") == "true" {
	// 	reqToken := getTokenFromHeader(c.Request)

	// 	if jwtParser == nil {
	// 		_jwtParser, err := jwt.NewJwtParser()
	// 		if err != nil {
	// 			logger.Errorf("failed to initialize jwt parser: %v", err.Error())
	// 			return AuthorizerInfo{}
	// 		}

	// 		jwtParser = _jwtParser
	// 	}

	// 	claims, err := jwtParser.GetJwtClaims(reqToken)
	// 	if err != nil {
	// 		logger.Errorf("failed to initialize jwt parser: %v", err.Error())
	// 	}

	// 	logger.Infof("claims = %v", claims)
	// 	info = AuthorizerInfo{
	// 		Claims: claims,
	// 	}

	// 	// Store claims in context
	// 	c.Request = c.Request.Clone(context.WithValue(ctx, CtxKeyAuthInfo, info))

	// 	logger.Infof("authorized info = %v", info)
	// }

	return info
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
