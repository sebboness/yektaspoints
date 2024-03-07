package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type IUserStorage interface {
	GetUserByID(ctx context.Context, userId string) (models.User, error)
	SaveUser(ctx context.Context, user models.User) error
	UpdateUserStatus(ctx context.Context, userId string, status models.UserStatus) error
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

type UpdateUserFamilyRequest struct {
	UserID   string
	FamilyID string
	Add      bool
}

// TODO need to finish this when ready to implement (-_o)
func (s *DynamoDbStorage) UpdateUserFamily(ctx context.Context, req UpdateUserFamilyRequest) error {

	familyIds := []string{}

	update := expression.Set(expression.Name("family_ids"), expression.Value(familyIds))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	keyEx, err := attributevalue.Marshal(req.UserID)
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(s.tableUser),
		Key:                       map[string]types.AttributeValue{"user_id": keyEx},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return apiErr
	}

	return nil
}

func (s *DynamoDbStorage) UpdateUserStatus(ctx context.Context, userId string, status models.UserStatus) error {

	update := expression.Set(expression.Name("status"), expression.Value(status))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	keyEx, err := attributevalue.Marshal(userId)
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(s.tableUser),
		Key:                       map[string]types.AttributeValue{"user_id": keyEx},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return apiErr
	}

	return nil
}
