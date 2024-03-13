package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/log"
	"github.com/sebboness/yektaspoints/util/result"
)

var logger = log.Get()

type Result struct {
	Status  string   `json:"status"`
	Errors  []string `json:"errors"`
	Message string   `json:"message"`
	Data    any      `json:"data,omitempty"`
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

func ErrorResult(err error) *result.Result {
	return result.ErrorResult(err)
}

func SuccessResult(data any) *result.Result {
	return result.SuccessResult(data)
}

// Use this mock API Gateway event in testing routes that require a user ID
var MockApiGWEvent = events.APIGatewayProxyRequest{
	RequestContext: events.APIGatewayProxyRequestContext{
		Authorizer: map[string]interface{}{
			"claims": map[string]interface{}{
				"cognito:username": "john",
				"email":            "john@info.co",
				"email_verified":   "true",
				"name":             "John",
				"sub":              "123",
			},
		},
	},
}
