package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/result"
)

func WithAuthorizedUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		authInfo := handlers.GetAuthorizerInfo(c)

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

		if !authInfo.IsEmailVerified() {
			// reject request
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("unverified user")))
			return
		}

		c.Next()
	}
}
