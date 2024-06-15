package points

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type getUserPointsHandlerRequest struct {
	// UserID is for the user that owns the points (child)
	UserID string `json:"-"`

	// User ID that makes the request (parent)
	RequestingUserID string `json:"-"`
}

type getUserPointsHandlerResponse struct {
	Points []models.Point `json:"points"`
}

func (c *PointsController) GetUserPointsHandler(cgin *gin.Context) {

	userID := cgin.Param("user_id")
	if userID == "" {
		apiErr := apierr.New(apierr.InvalidInput).WithError("user_id is a required query parameter")
		cgin.JSON(apiErr.StatusCode(), handlers.ErrorResult(apiErr))
		return
	}

	authInfo := c.AuthContext.GetAuthorizerInfo(cgin)

	req := &getUserPointsHandlerRequest{
		UserID:           userID,
		RequestingUserID: authInfo.GetUserID(),
	}

	resp, err := c.handleGetUserPoints(cgin.Request.Context(), req)
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

func (c *PointsController) handleGetUserPoints(ctx context.Context, req *getUserPointsHandlerRequest) (getUserPointsHandlerResponse, error) {
	resp := getUserPointsHandlerResponse{}

	if req.UserID == "" {
		return resp, apierr.New(apierr.AccessDenied).WithError("missing user id")
	}

	// check that the requested user (a parent) has access to data owned by user (a child)
	hasAccess, err := c.userDB.ParentHasAccessToChild(ctx, req.RequestingUserID, req.UserID)
	if err != nil {
		return resp, fmt.Errorf("failed to check access permissions: %w", err)
	}
	if !hasAccess {
		logger.WithFields(map[string]any{
			"requesting_user_id": req.RequestingUserID,
			"user_id":            req.UserID,
		})
		return resp, fmt.Errorf("user does not have permissions to get points from user: %w", err)
	}

	points, err := c.pointsDB.GetPointsByUserID(ctx, req.UserID, models.QueryPointsFilter{})
	if err != nil {
		return resp, fmt.Errorf("failed to get points: %w", err)
	}

	resp.Points = points
	return resp, nil
}
