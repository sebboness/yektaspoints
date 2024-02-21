package auth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/sebboness/yektaspoints/util/env"
)

type CognitoController struct {
	authClient          AuthClient
	cognitoClientID     string
	cognitoClientSecret string
}

func New(ctx context.Context) (AuthController, error) {
	cognitoClientID := env.GetEnv("COGNITO_CLIENT_ID")
	cognitoClientSecret := env.GetEnv("COGNITO_CLIENT_SECRET")
	return NewWithClient(ctx, cognitoClientID, cognitoClientSecret)
}

func NewWithClient(ctx context.Context, cognitoClientID, cognitoClientSecret string) (AuthController, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	cognitoClient := cognito.NewFromConfig(cfg)

	return &CognitoController{
		authClient:          cognitoClient,
		cognitoClientID:     cognitoClientID,
		cognitoClientSecret: cognitoClientSecret,
	}, nil
}

func (c *CognitoController) Authenticate(ctx context.Context, username, password string) (AuthResult, error) {
	result := AuthResult{}

	resp, err := c.authClient.InitiateAuth(ctx, &cognito.InitiateAuthInput{
		ClientId: &c.cognitoClientID,
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME":    username,
			"PASSWORD":    password,
			"SECRET_HASH": c.computeSecretHash(username),
		},
	})

	logger.WithContext(ctx).WithField("resp", resp).Infof("initiate auth response")

	if err != nil {
		return result, fmt.Errorf("failed to authenticate: %w", err)
	}

	if resp.ChallengeName == types.ChallengeNameTypeNewPasswordRequired {
		result.NewPasswordRequired = true
		result.Session = *resp.Session
		return result, nil
	}

	result.AccessToken = *resp.AuthenticationResult.AccessToken
	result.IdToken = *resp.AuthenticationResult.IdToken
	result.ExpiresIn = resp.AuthenticationResult.ExpiresIn

	// We need to grab the user record after authentication in order to store the "username" (aka the "sub") value
	// which we need for token refreshes later
	userResp, err := c.authClient.GetUser(ctx, &cognito.GetUserInput{AccessToken: &result.AccessToken})
	if err != nil {
		return result, fmt.Errorf("failed to get user after authentication: %w", err)
	}

	result.Username = *userResp.Username

	return result, nil
}

func (c *CognitoController) RefreshToken(ctx context.Context, username, refreshToken string) (AuthResult, error) {

	resp, err := c.authClient.InitiateAuth(ctx, &cognito.InitiateAuthInput{
		ClientId: &c.cognitoClientID,
		AuthFlow: types.AuthFlowTypeRefreshToken,
		AuthParameters: map[string]string{
			// "DEVICE_KEY": "", // TODO does this need to be set?
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   c.computeSecretHash(username),
		},
	})

	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	return AuthResult{
		AccessToken: *resp.AuthenticationResult.IdToken,
		ExpiresIn:   resp.AuthenticationResult.ExpiresIn,
	}, nil
}

func (c *CognitoController) UpdatePassword(ctx context.Context, session, username, password string) error {
	// accessToken := ""
	// attribPwName := ""
	// attribPwVal := ""

	resp, err := c.authClient.RespondToAuthChallenge(ctx, &cognito.RespondToAuthChallengeInput{
		Session:       &session,
		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
		ClientId:      &c.cognitoClientID,
		ChallengeResponses: map[string]string{
			// "USER_ID_FOR_SRP":     "sebboness",
			"USERNAME":            username,
			"NEW_PASSWORD":        password,
			"userAttributes.name": "Sebastian",
			// "userAttributes.password": password,
			"SECRET_HASH": c.computeSecretHash(username),
		},
	})

	logger.WithContext(ctx).WithField("resp", resp).Infof("update password response")

	if err != nil {
		logger.WithContext(ctx).WithField("error", err).Errorf("update password error")
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}

func (c *CognitoController) computeSecretHash(username string) string {
	return computeSecretHash(username, c.cognitoClientID, c.cognitoClientSecret)
}
