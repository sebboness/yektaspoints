package storage

import (
	"context"
	"fmt"
	"slices"

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
	ParentHasAccessToChild(ctx context.Context, parentId string, childId string) (bool, error)
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

func (s *DynamoDbStorage) ParentHasAccessToChild(ctx context.Context, parentId string, childId string) (bool, error) {
	if parentId == "" || childId == "" {
		return false, apierr.New(apierr.BadRequest).WithError("one or all user ids are empty")
	}

	parent, err := s.GetUserByID(ctx, parentId)
	if err != nil {
		return false, fmt.Errorf("failed to get parent user record for %v: %w", parentId, err)
	}

	child, err := s.GetUserByID(ctx, childId)
	if err != nil {
		return false, fmt.Errorf("failed to get child user record for %v: %w", childId, err)
	}

	// Check if parent has any associated family IDs
	if parent.FamilyIDs == nil {
		return false, fmt.Errorf("parent user %v has no associated families", parentId)
	}

	// Check if parent has parent role
	if parent.Roles == nil || !parent.IsParent() {
		return false, fmt.Errorf("user %v is not in role %v", parentId, models.RoleParent)
	}

	if child.FamilyIDs == nil {
		return false, fmt.Errorf("child user %v has no associated families", childId)
	}

	// Check if child has child role
	if child.Roles == nil || !child.IsChild() {
		return false, fmt.Errorf("user %v is not in role %v", childId, models.RoleChild)
	}

	parentHasAccess := false

	// Check if child belongs to any of parent's associated families
	for _, familyId := range parent.FamilyIDs {
		if slices.Contains(child.FamilyIDs, familyId) {
			parentHasAccess = true
			break
		}
	}

	return parentHasAccess, nil
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
