package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/util/log"
)

var allowList = map[string]bool{
	"http://localhost:3000":                 true,
	"https://mypoints.hexonite.net":         true,
	"https://mypoints-dev.hexonite.net":     true,
	"https://mypoints-staging.hexonite.net": true,
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.Get()
		logger.Infof("Request Origin is " + c.Request.Header.Get("Origin"))

		if origin := c.Request.Header.Get("Origin"); allowList[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,Origin,Cache-Control,X-Requested-With,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-CSRF-Token")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,GET,OPTIONS,PATCH,POST,PUT")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
