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

	authedRoutes := r.Group("/v1", mw.WithAuthorizedUser())

	// health
	authedRoutes.GET("/health", lambdaCtrl.HealthCheckHandler)

	// family
	authedRoutes.GET("/family/:family_id", familyCtrl.GetFamilyHandler, mw.WithRolesAny(groupsChildParent))

	// Points
	// authedRoutes.GET("/v1/points/:point_id", mw.WithRolesAny(groupsChildParent), pointsCtrl.GetPointsHandler)
	authedRoutes.POST("/points", mw.WithRolesAny(groupsChild), pointsCtrl.RequestPointsHandler)
	authedRoutes.PUT("/points/:point_id/approve", mw.WithRolesAny(groupsParent), pointsCtrl.ApprovePointsHandler)

	// Points (User)
	authedRoutes.GET("/user/:user_id/points", mw.WithRolesAny(groupsChildParent), pointsCtrl.GetUserPointsHandler)
	authedRoutes.GET("/user/:user_id/points-summary", mw.WithRolesAny(groupsChildParent), pointsCtrl.GetPointsSummaryHandler)

	// User
	authedRoutes.GET("/user", userCtrl.GetUserHandler)

	return r
}
