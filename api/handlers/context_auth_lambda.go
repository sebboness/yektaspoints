package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

type LambdaAuthContext struct {
	AuthContext
}

func NewLambdaAuthContext() (AuthContext, error) {
	return &LambdaAuthContext{}, nil
}

func (a *LambdaAuthContext) GetAuthorizerInfo(c *gin.Context) AuthorizerInfo {
	info := AuthorizerInfo{}

	if c.Request == nil {
		return info
	}

	ctx := c.Request.Context()
	_info := ctx.Value(CtxKeyAuthInfo)

	if _info != nil {
		return _info.(AuthorizerInfo)
	}

	return info
}

// PrepareAuthorizedContext prepares the given lambda api gateway context with authorizer info
func PrepareAuthorizedContext(ctx context.Context, req events.APIGatewayProxyRequest) context.Context {
	authorizer := AuthorizerInfo{}

	if len(req.RequestContext.Authorizer) > 0 {
		if claimsObj, ok := req.RequestContext.Authorizer["claims"]; ok {
			if claims, ok := claimsObj.(map[string]interface{}); ok {
				authorizer.Claims = claims
			}
		}
	}

	return context.WithValue(ctx, CtxKeyAuthInfo, authorizer)
}
