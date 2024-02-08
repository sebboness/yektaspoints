package request_points

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sebboness/yektaspoints/handlers"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/storage"
	"github.com/sebboness/yektaspoints/util"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/segmentio/ksuid"
)

type pointsHandlerRequest struct {
	events.APIGatewayProxyRequest

	Points int    `json:"points"`
	Reason string `json:"reason"`
	UserID string `json:"-"`
}

type pointsHandlerResponse struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

type RequestPointsController struct {
	pointsDB storage.IPointsStorage
}

func (c *RequestPointsController) RequestPointsHandler(ctx context.Context, event *pointsHandlerRequest) (events.APIGatewayProxyResponse, error) {

	logger.WithContext(ctx).Infof("authorizer: %+v", event.RequestContext.Authorizer)
	event.UserID = handlers.GetUserIDFromLambdaRequest(&event.APIGatewayProxyRequest)

	resp, err := c.handleRequestPoints(ctx, event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			// Body: nil,
		}, err
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			// Body: nil,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func (c *RequestPointsController) handleRequestPoints(ctx context.Context, req *pointsHandlerRequest) (pointsHandlerResponse, error) {
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
		return apierr.New(apierr.Unauthorized, "missing user ID")
	}

	apierr := apierr.New(apierr.InvalidInput, "failed to validate request")

	if req.Points <= 0 {
		apierr.AppendError("points must be a positive integer")
	}

	// Arbitrary check for some valid reason text
	// TODO: make it better
	if req.Reason == "" || len(req.Reason) <= 5 {
		apierr.AppendError("reason for requesting points must not be empty")
	}

	if len(apierr.Errors) > 0 {
		return apierr
	}

	return nil
}
