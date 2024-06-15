package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/result"
)

// WithRolesAny checks if the current user's roles matches any of the given roles.
// If no match is found, the request returns a 401; Otherwise the request continues.
func WithRolesAny(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authContext, err := handlers.GetAuthContext()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, result.ErrorResult(fmt.Errorf("failed to get auth context: %w", err)))
			return
		}

		authInfo := authContext.GetAuthorizerInfo(c)

		matches := false
		logger.Infof("one of roles required: %v; auth roles: %v", roles, authInfo.GetGroups())
		for _, role := range roles {
			logger.Infof("auth roles (%v) contains %v? %v", authInfo.GetGroups(), role, slices.Contains(authInfo.GetGroups(), role))
			if slices.Contains(authInfo.GetGroups(), role) {
				matches = true
				break
			}
		}

		if !matches {
			logger.Infof("hit here")
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("user not a parent")))
			return
		}

		c.Next()
	}
}
