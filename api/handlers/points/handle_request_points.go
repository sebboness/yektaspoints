package points

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type pointsHandlerRequest struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

type pointsHandlerResponse struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

func RequestPointsHandler(ctx context.Context, event *pointsHandlerRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := handleRequestPoints(event)
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

func handleRequestPoints(req *pointsHandlerRequest) (pointsHandlerResponse, error) {
	resp := pointsHandlerResponse{}

	if err := validateRequestPoints(req); err != nil {
		return resp, err
	}

	return resp, nil
}

func validateRequestPoints(req *pointsHandlerRequest) error {
	return nil
}

func main() {
	lambda.Start(RequestPointsHandler)
}
