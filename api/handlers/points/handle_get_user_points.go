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
	UserID string `json:"-"`
}

type getUserPointsHandlerResponse struct {
	Points []models.Point `json:"points"`
}

func (c *PointsController) GetUserPointsHandler(cgin *gin.Context) {

	// authInfo := handlers.GetAuthorizerInfo(cgin)
	// TODO: Implement a way to check that the requested user (a parent) has access to retrieve the points for the user
	//       identified in this request via the user_id query parameter (their child).

	userID := cgin.Param("user_id")
	if userID == "" {
		apiErr := apierr.New(apierr.InvalidInput).WithError("user_id is a required query parameter")
		cgin.JSON(apiErr.StatusCode(), handlers.ErrorResult(apiErr))
		return
	}

	req := &getUserPointsHandlerRequest{
		UserID: userID,
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

	points, err := c.pointsDB.GetPointsByUserID(ctx, req.UserID, models.QueryPointsFilter{})
	if err != nil {
		return resp, fmt.Errorf("failed to get points: %w", err)
	}

	resp.Points = points
	return resp, nil
}
