package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/middleware"
)

func RegisterRoutes(r *gin.Engine) *gin.Engine {
	r.Use(gin.Recovery()).Use(middleware.CORSMiddleware())

	// Health
	r.GET("/", lambdaCtrl.HealthCheckHandler)
	r.GET("/health", lambdaCtrl.HealthCheckHandler)

	// Auth
	r.POST("/auth/token", authCtrl.UserAuthHandler)

	// User registration
	r.POST("/v1/user/register", userCtrl.UserRegisterHandler)
	r.POST("/v1/user/register/confirm", userCtrl.UserRegisterConfirmHandler)

	authedUserRoutes := r.Group("/v1")
	authedUserRoutes.Use(middleware.WithAuthorizedUser())
	{
		r.GET("/v1/health", lambdaCtrl.HealthCheckHandler)

		// family
		r.GET("/v1/family", familyCtrl.GetFamilyHandler)

		// Points
		r.GET("/v1/points", pointsCtrl.GetUserPointsHandler)
		r.GET("/v1/points/:point_id", pointsCtrl.GetUserPointsHandler)
		r.POST("/v1/points", pointsCtrl.RequestPointsHandler)

		// User
		r.GET("/v1/user", userCtrl.GetUserHandler)
	}

	return r
}
