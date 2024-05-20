package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/result"
)

func WithParentUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		authInfo := handlers.GetAuthorizerInfo(c)

		if !slices.Contains(authInfo.GetGroups(), "parent") {
			// reject request
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("user not a parent")))
			return
		}

		c.Next()
	}
}
