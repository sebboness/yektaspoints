package handlers

import "github.com/aws/aws-lambda-go/events"

// Use this mock API Gateway event in testing routes that require a user ID
var mockApiGWEvent = events.APIGatewayProxyRequest{
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
