package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sebboness/yektaspoints/storage"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/result"
)

type Result struct {
	Status  string   `json:"status"`
	Errors  []string `json:"errors"`
	Message string   `json:"message"`
	Data    any      `json:"data,omitempty"`
}

type PointsController struct {
	pointsDB storage.IPointsStorage
}

func NewPointsController(env string) (*PointsController, error) {
	storageCfg := storage.Config{Env: env}

	pointsDB, err := storage.NewDynamoDbStorage(storageCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize points db: %w", err)
	}

	return &PointsController{
		pointsDB: pointsDB,
	}, nil
}

// GetUserIDFromLambdaRequest returns the user ID from the lambda request.
// Specifically it returns the cognito:username claim value from the authorizer claims map
func GetUserIDFromLambdaRequest(req *events.APIGatewayProxyRequest) string {
	userID := ""
	if claims, ok := req.RequestContext.Authorizer["claims"]; ok {
		claimsMap := claims.(map[string]any)
		userID = claimsMap["cognito:username"].(string)
	}
	return userID
}

func ApiResponseOK(data any) events.APIGatewayProxyResponse {
	r := result.SuccessResult(data)
	rjson, err := json.Marshal(r)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("failed to marshal json response: %v", err),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(rjson),
	}
}

func ApiResponseWithError(statusCode int, err error) events.APIGatewayProxyResponse {
	r := result.ErrorResult(err)
	rjson, jerr := json.Marshal(r)
	if jerr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("failed to marshal json response: %v. original error: %v", jerr, err),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(rjson),
	}
}

func ApiResponseBadRequest(err error) events.APIGatewayProxyResponse {
	return ApiResponseWithError(http.StatusBadRequest, err)
}

func ApiResponseNotFound(err error) events.APIGatewayProxyResponse {
	return ApiResponseWithError(http.StatusNotFound, err)
}

func ApiResponseUnauthorized(err error) events.APIGatewayProxyResponse {
	return ApiResponseWithError(http.StatusUnauthorized, err)
}

func ApiResponseInternalServerError(err error) events.APIGatewayProxyResponse {
	return ApiResponseWithError(http.StatusInternalServerError, err)
}

func ApiErrorResponse(err *apierr.ApiError) events.APIGatewayProxyResponse {
	return ApiResponseWithError(err.StatusCode(), err.Err)
}
