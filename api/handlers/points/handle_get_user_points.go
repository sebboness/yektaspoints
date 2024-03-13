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

	authInfo := handlers.GetAuthorizerInfo(cgin)

	req := &getUserPointsHandlerRequest{
		UserID: authInfo.GetUserID(),
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
