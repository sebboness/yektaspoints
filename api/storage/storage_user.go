package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type IUserStorage interface {
	GetUserByID(ctx context.Context, userId string) (models.User, error)
	SaveUser(ctx context.Context, user models.User) error
}

func (s *DynamoDbStorage) GetUserByID(ctx context.Context, userId string) (models.User, error) {
	user := models.User{}

	keyEx := expression.Key("user_id").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		return user, fmt.Errorf("failed to build query expression: %w", err)
	}

	resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(s.tableUser),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ScanIndexForward:          aws.Bool(false),
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return user, apiErr
	}

	if len(resp.Items) == 0 {
		logger.WithContext(ctx).WithField("userId", userId).Warnf("item (user_id:%s) not found", userId)
		return user, apierr.New(apierr.NotFound).WithError(fmt.Sprintf("user (user_id=%s)", userId))
	}

	err = attributevalue.UnmarshalMap(resp.Items[0], &user)
	if err != nil {
		return user, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	user.ParseTimes()

	return user, nil
}

func (s *DynamoDbStorage) SaveUser(ctx context.Context, user models.User) error {

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal map from point: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableUser),
		Item:      item,
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return apiErr
	}

	return nil
}
