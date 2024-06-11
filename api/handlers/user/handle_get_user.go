package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/log"
)

type getUserResponse struct {
	Email     string   `json:"email"`
	FamilyIDs []string `json:"family_ids"`
	Name      string   `json:"name"`
	UserID    string   `json:"user_id"`
	Username  string   `json:"username"`
	Roles     []string `json:"roles"`
}

// GetUserHandler returns user data from the currently logged-in user.
func (c *UserController) GetUserHandler(cgin *gin.Context) {

	authInfo := c.AuthContext.GetAuthorizerInfo(cgin)
	if !authInfo.HasInfo() {
		cgin.JSON(http.StatusUnauthorized, handlers.ErrorResult(apierr.Unauthorized))
		return
	}

	resp, err := c.handleGetUser(cgin.Request.Context(), authInfo.GetUserID())
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			cgin.JSON(apierr.StatusCode(), handlers.ErrorResult(apierr))
			return
		}

		cgin.JSON(http.StatusInternalServerError, handlers.ErrorResult(err))
		return
	}

	cgin.JSON(http.StatusOK, handlers.SuccessResult(resp))
}

func (c *UserController) handleGetUser(ctx context.Context, userId string) (getUserResponse, error) {
	resp := getUserResponse{}
	logger := log.Get().WithContext(ctx).AddFields(map[string]any{
		"user_id": userId,
	})

	user, err := c.userDB.GetUserByID(ctx, userId)
	if err != nil {
		logger.WithFields(map[string]any{"error": err.Error()}).Errorf("failed to get user")
		return resp, fmt.Errorf("failed to get user: %w", err)
	}

	resp.Email = user.Email
	resp.FamilyIDs = user.FamilyIDs
	resp.Name = user.Name
	resp.Roles = user.Roles
	resp.UserID = user.UserID
	resp.Username = user.Username

	return resp, nil
}
