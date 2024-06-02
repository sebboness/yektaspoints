package main

import (
	"github.com/gin-gonic/gin"
	mw "github.com/sebboness/yektaspoints/middleware"
)

const groupChild = "child"
const groupParent = "parent"

var groupsChild = []string{groupChild}
var groupsChildParent = []string{groupChild, groupParent}
var groupsParent = []string{groupParent}

func RegisterRoutes(r *gin.Engine) *gin.Engine {
	r.Use(gin.Recovery()).Use(mw.CORSMiddleware())

	// Health
	r.GET("/", lambdaCtrl.HealthCheckHandler)
	r.GET("/health", lambdaCtrl.HealthCheckHandler)

	// Auth
	r.POST("/auth/token", authCtrl.UserAuthHandler)

	// User registration
	r.POST("/v1/user/register", userCtrl.UserRegisterHandler)
	r.POST("/v1/user/register/confirm", userCtrl.UserRegisterConfirmHandler)

	authedUserRoutes := r.Group("/v1")
	authedUserRoutes.Use(mw.WithAuthorizedUser())
	{
		r.GET("/v1/health", lambdaCtrl.HealthCheckHandler)

		// family
		r.GET("/v1/family", familyCtrl.GetFamilyHandler, mw.WithRolesAny(groupsChildParent))

		// Points
		r.GET("/v1/points/:point_id", pointsCtrl.GetUserPointsHandler, mw.WithRolesAny(groupsChildParent))
		r.GET("/v1/points/summary/:user_id", pointsCtrl.GetPointsSummaryHandler, mw.WithRolesAny(groupsChildParent))
		r.POST("/v1/points", pointsCtrl.RequestPointsHandler, mw.WithRolesAny(groupsChild))

		// Points (User)
		r.GET("/v1/points/user/:user_id", pointsCtrl.GetUserPointsHandler, mw.WithRolesAny(groupsChildParent))
		r.POST("/v1/points/user/:user_id/approve", pointsCtrl.ApprovePointsHandler, mw.WithRolesAny(groupsParent))

		// User
		r.GET("/v1/user", userCtrl.GetUserHandler)
	}

	return r
}
