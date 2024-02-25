package auth

import (
	"context"

	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type AuthClient interface {
	AdminAddUserToGroup(ctx context.Context, params *cognito.AdminAddUserToGroupInput, optFns ...func(*cognito.Options)) (*cognito.AdminAddUserToGroupOutput, error)
	ConfirmSignUp(ctx context.Context, params *cognito.ConfirmSignUpInput, optFns ...func(*cognito.Options)) (*cognito.ConfirmSignUpOutput, error)
	GetUser(ctx context.Context, params *cognito.GetUserInput, optFns ...func(*cognito.Options)) (*cognito.GetUserOutput, error)
	InitiateAuth(ctx context.Context, params *cognito.InitiateAuthInput, optFns ...func(*cognito.Options)) (*cognito.InitiateAuthOutput, error)
	RespondToAuthChallenge(ctx context.Context, params *cognito.RespondToAuthChallengeInput, optFns ...func(*cognito.Options)) (*cognito.RespondToAuthChallengeOutput, error)
	SignUp(ctx context.Context, params *cognito.SignUpInput, optFns ...func(*cognito.Options)) (*cognito.SignUpOutput, error)
	UpdateUserAttributes(ctx context.Context, params *cognito.UpdateUserAttributesInput, optFns ...func(*cognito.Options)) (*cognito.UpdateUserAttributesOutput, error)
}
