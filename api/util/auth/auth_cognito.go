package auth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	apierr "github.com/sebboness/yektaspoints/util/error"

	"github.com/sebboness/yektaspoints/util/env"
)

type CognitoController struct {
	authClient          AuthClient
	cognitoClientID     string
	cognitoClientSecret string
	userPoolID          string
}

func New(ctx context.Context) (AuthController, error) {
	cognitoClientID := env.GetEnv("COGNITO_CLIENT_ID")
	cognitoClientSecret := env.GetEnv("COGNITO_CLIENT_SECRET")
	userPoolID := env.GetEnv("COGNITO_USER_POOL_ID")
	return NewWithClient(ctx, cognitoClientID, cognitoClientSecret, userPoolID)
}

func NewWithClient(ctx context.Context, cognitoClientID, cognitoClientSecret, userPoolID string) (AuthController, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	cognitoClient := cognito.NewFromConfig(cfg)

	return &CognitoController{
		authClient:          cognitoClient,
		cognitoClientID:     cognitoClientID,
		cognitoClientSecret: cognitoClientSecret,
		userPoolID:          userPoolID,
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
		apiErr := apierr.GetAwsError(err)
		return result, apiErr
	}

	if resp.ChallengeName == types.ChallengeNameTypeNewPasswordRequired {
		result.NewPasswordRequired = true
		result.Session = *resp.Session
		return result, nil
	}

	result.AccessToken = *resp.AuthenticationResult.AccessToken
	result.IdToken = *resp.AuthenticationResult.IdToken
	result.ExpiresIn = resp.AuthenticationResult.ExpiresIn
	result.RefreshToken = *resp.AuthenticationResult.RefreshToken

	// We need to grab the user record after authentication in order to store the "username" (aka the "sub") value
	// which we need for token refreshes later
	userResp, err := c.authClient.GetUser(ctx, &cognito.GetUserInput{AccessToken: &result.AccessToken})
	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return result, apiErr
	}

	result.Username = *userResp.Username

	return result, nil
}

func (c *CognitoController) AssignUserToRole(ctx context.Context, username, role string) error {

	resp, err := c.authClient.AdminAddUserToGroup(ctx, &cognito.AdminAddUserToGroupInput{
		GroupName:  aws.String(role),
		Username:   aws.String(username),
		UserPoolId: aws.String(c.userPoolID),
	})

	logger.WithField("resp", resp).Infof("user confirm signup response")

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return apiErr
	}

	return nil
}

func (c *CognitoController) ConfirmRegistration(ctx context.Context, username, code string) error {

	resp, err := c.authClient.ConfirmSignUp(ctx, &cognito.ConfirmSignUpInput{
		ConfirmationCode: aws.String(code),
		Username:         aws.String(username),
		ClientId:         aws.String(c.cognitoClientID),
		SecretHash:       aws.String(c.computeSecretHash(username)),
	})

	logger.WithField("resp", resp).Infof("user confirm signup response")

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return apiErr
	}

	return nil
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
		apiErr := apierr.GetAwsError(err)
		return AuthResult{}, apiErr
	}

	return AuthResult{
		AccessToken: *resp.AuthenticationResult.IdToken,
		ExpiresIn:   resp.AuthenticationResult.ExpiresIn,
	}, nil
}

func (c *CognitoController) Register(ctx context.Context, req UserRegisterRequest) (UserRegisterResult, error) {
	result := UserRegisterResult{}

	resp, err := c.authClient.SignUp(ctx, &cognito.SignUpInput{
		Username:   aws.String(req.Username),
		Password:   aws.String(req.Password),
		ClientId:   aws.String(c.cognitoClientID),
		SecretHash: aws.String(c.computeSecretHash(req.Username)),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(req.Name),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(req.Email),
			},
		},
	})

	logger.WithField("resp", resp).Infof("user signup response")

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return result, apiErr
	}

	result.IsConfirmed = resp.UserConfirmed
	result.Username = *resp.UserSub
	result.ConfirmationType = string(resp.CodeDeliveryDetails.DeliveryMedium)
	result.ConfirmationSentTo = *resp.CodeDeliveryDetails.Destination

	return result, nil
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
		apiErr := apierr.GetAwsError(err)
		return apiErr
	}

	return nil
}

func (c *CognitoController) computeSecretHash(username string) string {
	return computeSecretHash(username, c.cognitoClientID, c.cognitoClientSecret)
}
