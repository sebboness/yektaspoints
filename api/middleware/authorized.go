package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/log"
	"github.com/sebboness/yektaspoints/util/result"
)

func WithAuthorizedUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		logger := log.Get()

		authInfo := handlers.GetAuthorizerInfo(c)

		authInfoJson, _ := json.Marshal(authInfo)
		logger.Infof("WithAuthorizedUser: " + string(authInfoJson))

		if !authInfo.HasInfo() {
			// reject request
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("unauthorized")))
			return
		}

		if authInfo.GetUserID() == "" {
			// reject request
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("unknown user ID")))
			return
		}

		c.Next()
	}
}
