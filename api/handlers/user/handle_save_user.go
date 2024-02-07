package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/models"
)

type saveUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (n *UserController) SaveUserHandler(c *gin.Context) {

	var req saveUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{})
	}

	user, err := handleSaveUser(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{})
	}

	c.IndentedJSON(http.StatusOK, user)
}

func handleSaveUser(req saveUserRequest) (models.User, error) {
	user := models.User{}

	return user, nil
}

func validateSaveUser() error {
	return nil
}
