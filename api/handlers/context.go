package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	ctxKeyAuthInfo ctxKey = "api:auth"

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

	return context.WithValue(ctx, ctxKeyAuthInfo, authorizer)
}

func GetAuthorizerInfo(c *gin.Context) AuthorizerInfo {
	if c.Request != nil {
		ctx := c.Request.Context()
		info := ctx.Value(ctxKeyAuthInfo)

		authInfoJson, _ := json.Marshal(info)
		logger.Infof("GetAuthorizerInfo authInfo?: " + string(authInfoJson))

		if info != nil {
			return info.(AuthorizerInfo)
		}
	}

	return AuthorizerInfo{}
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
