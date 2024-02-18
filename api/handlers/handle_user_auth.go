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
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
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
	result := auth.AuthResult{}

	if err := validateUserAuth(req); err != nil {
		return resp, err
	}

	// We validated that grant type can only be one of the following cases
	if req.GrantType == auth.GrantTypePassword {
		authResult, err := c.auth.Authenticate(ctx, req.Username, req.Password)
		if err != nil {
			return resp, fmt.Errorf("failed to authenticate: %w", err)
		}
		result = authResult
	} else if req.GrantType == auth.GrantTypeRefreshToken {
		authResult, err := c.auth.RefreshToken(ctx, req.Username, req.RefreshToken)
		if err != nil {
			return resp, fmt.Errorf("failed to refresh token: %w", err)
		}
		result = authResult
	}

	resp.AuthResult = result
	return resp, nil
}

func validateUserAuth(req *userAuthRequest) error {
	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if _, ok := auth.SupportedGrantTypes[req.GrantType]; !ok {
		apierr.AppendErrorf("unsupported grant_type \"%s\"", req.GrantType)
	}

	if req.GrantType == auth.GrantTypePassword {
		if req.Username == "" {
			apierr.AppendError("missing username")
		}
		if req.Password == "" {
			apierr.AppendError("missing password")
		}
	} else if req.GrantType == auth.GrantTypeRefreshToken {
		if req.Username == "" {
			apierr.AppendError("missing username")
		}
		if req.RefreshToken == "" {
			apierr.AppendError("missing refresh_token")
		}
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
