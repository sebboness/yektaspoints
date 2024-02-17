package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userAuthResponse struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

func (c *PointsController) UserAuthHandler(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var req userAuthRequest

	// try to unmarshal from request body
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal json body: %w", err)
		return ApiResponseInternalServerError(err), err
	}

	req.Password = GetUserIDFromLambdaRequest(event)

	resp, err := c.handleUserAuth(ctx, &req)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			return ApiErrorResponse(apierr), apierr
		}

		return ApiResponseInternalServerError(err), err
	}

	return ApiResponseOK(resp), nil
}

func (c *PointsController) handleUserAuth(ctx context.Context, req *userAuthRequest) (userAuthResponse, error) {
	resp := userAuthResponse{}

	if err := validateUserAuth(req); err != nil {
		return resp, err
	}

	return resp, nil
}

func validateUserAuth(req *userAuthRequest) error {
	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if req.Username == "" {
		apierr.AppendError("missing username")
	}
	if req.Password == "" {
		apierr.AppendError("missing password")
	}

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
