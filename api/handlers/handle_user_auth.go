package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sebboness/yektaspoints/util/auth"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userAuthResponse struct {
	auth.AuthResult
}

// UserAuthHandler authenticates a user using a username/password auth flow
func (c *LambdaController) UserAuthHandler(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var req userAuthRequest

	// try to unmarshal from request body
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal json body: %w", err)
		return ApiResponseInternalServerError(err), err
	}

	resp, err := c.handleUserAuth(ctx, &req)
	if err != nil {
		if apierr := apierr.IsApiError(err); apierr != nil {
			return ApiErrorResponse(apierr), apierr
		}

		return ApiResponseInternalServerError(err), err
	}

	return ApiResponseOK(resp), nil
}

func (c *LambdaController) handleUserAuth(ctx context.Context, req *userAuthRequest) (userAuthResponse, error) {
	resp := userAuthResponse{}

	if err := validateUserAuth(req); err != nil {
		return resp, err
	}

	result, err := c.auth.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return resp, fmt.Errorf("failed to authenticate: %w", err)
	}

	resp.AuthResult = result
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

	// pwResult := auth.ValidatePassword(req.Password)
	// if !pwResult.WithinLength {
	// 	apierr.AppendError("password must be within 8 and 256 characters in length")
	// }
	// if !pwResult.Lower {
	// 	apierr.AppendError("password must have at least one upper case letter")
	// }
	// if !pwResult.Upper {
	// 	apierr.AppendError("password must have at least one lower case letter")
	// }
	// if !pwResult.Number {
	// 	apierr.AppendError("password must have at least one digit")
	// }
	// if !pwResult.Special {
	// 	apierr.AppendError("password must have at least one special character")
	// }

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
