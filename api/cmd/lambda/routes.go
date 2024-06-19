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
		// r.GET("/v1/points/:point_id", mw.WithRolesAny(groupsChildParent), pointsCtrl.GetPointsHandler)
		r.POST("/v1/points", mw.WithRolesAny(groupsChild), pointsCtrl.RequestPointsHandler)
		r.PUT("/v1/points/:point_id/approve", mw.WithRolesAny(groupsParent), pointsCtrl.ApprovePointsHandler)

		// Points (User)
		r.GET("/v1/user/:user_id/points", mw.WithRolesAny(groupsChildParent), pointsCtrl.GetUserPointsHandler)
		r.GET("/v1/user/:user_id/points-summary", mw.WithRolesAny(groupsChildParent), pointsCtrl.GetPointsSummaryHandler)

		// User
		r.GET("/v1/user", userCtrl.GetUserHandler)
	}

	return r
}
