package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/middleware"
)

func RegisterRoutes(r *gin.Engine) *gin.Engine {
	r.Use(gin.Recovery())

	// Health
	r.GET("/", lambdaCtrl.HealthCheckHandler)
	r.GET("/health", lambdaCtrl.HealthCheckHandler)

	// Auth
	r.POST("/auth/token", lambdaCtrl.UserAuthHandler)

	// User registration
	r.POST("/v1/user/register", userCtrl.UserRegisterHandler)
	r.POST("/v1/user/register/confirm", userCtrl.UserRegisterConfirmHandler)

	authedUserRoutes := r.Group("/v1")
	authedUserRoutes.Use(middleware.WithAuthorizedUser())
	{
		r.GET("/v1/health", lambdaCtrl.HealthCheckHandler)

		// Points
		r.GET("/v1/points", lambdaCtrl.GetUserPointsHandler)
		r.GET("/v1/points/:point_id", lambdaCtrl.GetUserPointsHandler)
		r.POST("/v1/points", lambdaCtrl.RequestPointsHandler)
	}

	return r
}
