package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) *gin.Engine {
	// Auth
	r.POST("/auth/token", lambdaCtrl.UserAuthHandler)

	// Health
	r.GET("/", lambdaCtrl.HealthCheckHandler)
	r.GET("/health", lambdaCtrl.HealthCheckHandler)
	r.GET("/v1/health", lambdaCtrl.HealthCheckHandler)

	// Points
	r.GET("/v1/points", lambdaCtrl.GetUserPointsHandler)
	r.GET("/v1/points/:point_id", lambdaCtrl.GetUserPointsHandler)
	r.POST("/v1/points", lambdaCtrl.RequestPointsHandler)

	// User
	r.POST("/v1/user/register", userCtrl.UserRegisterHandler)

	return r
}
