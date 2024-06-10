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

		authInfo := handlers.GetAuthorizerInfo(c)

		matches := false
		for _, role := range roles {
			if slices.Contains(authInfo.GetGroups(), role) {
				matches = true
				break
			}
		}

		if !matches {
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.ErrorResult(fmt.Errorf("user not a parent")))
			return
		}

		c.Next()
	}
}