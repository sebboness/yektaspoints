package userauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/util/auth"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type getUserAuthResponse struct {
	auth.AuthResult
}

// GetUserAuthHandler authenticates a user depending on the request grant_type
func (c *UserAuthController) GetUserAuthHandler(cgin *gin.Context) {

	authInfo := handlers.GetAuthorizerInfo(cgin)
	if !authInfo.HasInfo() {
		cgin.JSON(http.StatusUnauthorized, handlers.ErrorResult(apierr.Unauthorized))
		return
	}

	resp, err := c.handleGetUserAuth(cgin.Request.Context(), authInfo.GetUsername())
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

func (c *UserAuthController) handleGetUserAuth(ctx context.Context, username string) error {
	resp := getUserAuthResponse{}

	result, err := c.auth.RefreshToken(ctx, username)
	if err != nil {
		return resp, fmt.Errorf("failed to refresh token: %w", err)
	}

	return nil
}
