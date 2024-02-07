package main

import (
	"github.com/gin-gonic/gin"

	"github.com/sebboness/yektaspoints/handlers/user"
)

func main() {
	router := gin.Default()

	userController := &user.UserController{}

	router.POST("/users", userController.SaveUserHandler)

	router.Run("localhost:8080")
}
