package userauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type getUserAuthResponse struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Groups   []string `json:"groups"`
}

// GetUserAuthHandler returns auth info from the currently logged in user.
// We get this info from the passed in auth token via Cognito
func (c *UserAuthController) GetUserAuthHandler(cgin *gin.Context) {

	authInfo := handlers.GetAuthorizerInfo(cgin)
	if !authInfo.HasInfo() {
		cgin.JSON(http.StatusUnauthorized, handlers.ErrorResult(apierr.Unauthorized))
		return
	}

	resp := getUserAuthResponse{
		UserID:   authInfo.GetUserID(),
		Username: authInfo.GetUsername(),
		Email:    authInfo.GetEmail(),
		Name:     authInfo.GetName(),
		Groups:   authInfo.GetGroups(),
	}

	cgin.JSON(http.StatusOK, handlers.SuccessResult(resp))
}
