package handlers

import (
	"context"

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
)

type AuthorizerInfo struct {
	Claims map[string]string `json:"claims"`
}

func PrepareAuthorizedContext(ctx context.Context, req events.APIGatewayProxyRequest) context.Context {
	authorizer := AuthorizerInfo{}

	if len(req.RequestContext.Authorizer) > 0 {
		if claimsObj, ok := req.RequestContext.Authorizer["claims"]; ok {
			if claims, ok := claimsObj.(map[string]string); ok {
				authorizer.Claims = claims
			}
		}
	}

	return context.WithValue(ctx, ctxKeyAuthInfo, authorizer)
}

func GetAuthorizerInfo(c *gin.Context) AuthorizerInfo {
	if c.Request != nil {
		ctx := c.Request.Context()
		info := ctx.Value(ctxKeyAuthInfo)
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
		return i.Claims[key]
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

func (i AuthorizerInfo) GetName() string {
	return i.ValueOrEmpty(claimKeyName)
}

func (i AuthorizerInfo) GetUserID() string {
	return i.ValueOrEmpty(claimKeyUserID)
}

func (i AuthorizerInfo) GetUsername() string {
	return i.ValueOrEmpty(claimKeyUsername)
}
