package handlers

import "github.com/aws/aws-lambda-go/events"

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
