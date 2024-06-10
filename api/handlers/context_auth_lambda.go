package handlers

import (
	"encoding/json"

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

	authInfoJson, _ := json.Marshal(info)
	logger.Infof("GetAuthorizerInfo authInfo?: " + string(authInfoJson))

	if _info != nil {
		return _info.(AuthorizerInfo)
	}

	return info
}
