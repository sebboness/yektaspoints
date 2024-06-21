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

		authContext, err := handlers.GetAuthContext()
		if err != nil {
			logger.Errorf("failed to get auth context: %v", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, result.ErrorResult(fmt.Errorf("failed to get auth context: %w", err)))
			return
		}

		authInfo := authContext.GetAuthorizerInfo(c)

		if !authInfo.HasInfo() {
			// reject request
			logger.Errorf("unauthorized request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("unauthorized")))
			return
		}

		if authInfo.GetUserID() == "" {
			// reject request
			logger.Errorf("unauthorized request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("unknown user ID")))
			return
		}

		if !authInfo.IsEmailVerified() {
			// reject request
			logger.Errorf("unverified user")
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("unverified user")))
			return
		}

		logger.Infof("getting here?")

		c.Next()
	}
}
