package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/result"
)

func WithAuthorizer() gin.HandlerFunc {
	return func(c *gin.Context) {

		authInfo := handlers.GetAuthorizerInfo(c)
		if !authInfo.HasInfo() {
			// reject request
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("unauthorized")))
			return
		}

		c.Next()
	}
}
