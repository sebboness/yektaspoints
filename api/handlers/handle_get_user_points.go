package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/segmentio/ksuid"
)

type getUserPointsHandlerRequest struct {
	events.APIGatewayProxyRequest

	Points int    `json:"points"`
	Reason string `json:"reason"`
	UserID string `json:"-"`
}

type getUserPointsHandlerResponse struct {
	Points []models.Point `json:"points"`
}

func (c *PointsController) GetUserPointsHandler(ctx context.Context, event *getUserPointsHandlerRequest) (events.APIGatewayProxyResponse, error) {

	logger.WithContext(ctx).Infof("authorizer: %+v", event.RequestContext.Authorizer)
	event.UserID = GetUserIDFromLambdaRequest(&event.APIGatewayProxyRequest)

	resp, err := c.handleGetUserPoints(ctx, event)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			return ApiErrorResponse(apierr), apierr
		}

		return ApiResponseInternalServerError(err), err
	}

	return ApiResponseOK(resp), nil
}

func (c *PointsController) handleGetUserPoints(ctx context.Context, req *getUserPointsHandlerRequest) (getUserPointsHandlerResponse, error) {
	resp := getUserPointsHandlerResponse{}

	if err := validateGetUserPoints(req); err != nil {
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

func validateGetUserPoints(req *getUserPointsHandlerRequest) error {
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
