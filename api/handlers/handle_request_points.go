package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/segmentio/ksuid"
)

type pointsHandlerRequest struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
	UserID string `json:"-"`
}

type pointsHandlerResponse struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

func (c *LambdaController) RequestPointsHandler(cgin *gin.Context) {

	var req pointsHandlerRequest

	// try to unmarshal from request body
	err := cgin.BindJSON(&req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal json body: %w", err)
		cgin.JSON(http.StatusBadRequest, ErrorResult(err))
		return
	}

	authInfo := GetAuthorizerInfo(cgin)
	req.UserID = authInfo.GetUserID()

	resp, err := c.handleRequestPoints(cgin.Request.Context(), &req)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			cgin.JSON(apierr.StatusCode(), ErrorResult(apierr))
			return
		}

		cgin.JSON(http.StatusInternalServerError, ErrorResult(err))
		return
	}

	cgin.JSON(http.StatusOK, SuccessResult(resp))
}

func (c *LambdaController) handleRequestPoints(ctx context.Context, req *pointsHandlerRequest) (pointsHandlerResponse, error) {
	resp := pointsHandlerResponse{}

	if err := validateRequestPoints(req); err != nil {
		return resp, err
	}

	point := models.Point{
		ID:             ksuid.New().String(),
		UserID:         req.UserID,
		Points:         req.Points,
		Reason:         req.Reason,
		StatusID:       models.PointStatusRequested,
		Type:           models.PointTypeAdd,
		RequestedOnStr: util.ToFormattedUTC(time.Now()),
	}

	err := c.pointsDB.SavePoint(ctx, point)
	if err != nil {
		return resp, fmt.Errorf("failed to save points: %w", err)
	}

	return resp, nil
}

func validateRequestPoints(req *pointsHandlerRequest) error {
	if req.UserID == "" {
		return apierr.New(fmt.Errorf("%w: missing user ID", apierr.Unauthorized))
	}

	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if req.Points <= 0 {
		apierr.AppendError("points must be a positive integer")
	}

	// Arbitrary check for some valid reason text
	// TODO: make it better
	if req.Reason == "" || len(req.Reason) <= 5 {
		apierr.AppendError("reason for requesting points must not be empty")
	}

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
