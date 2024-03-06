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
	GetFamilyByFamilyID(ctx context.Context, family_id string) (models.Family, error)
	GetUserByID(ctx context.Context, userId string) (models.User, error)
}

func (s *DynamoDbStorage) GetFamilyByFamilyID(ctx context.Context, family_id string) (models.Family, error) {
	family := models.Family{
		FamilyID: family_id,
		Parents:  map[string]models.FamilyUser{},
		Children: map[string]models.FamilyUser{},
	}

	params, err := attributevalue.MarshalList([]interface{}{family_id})
	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return family, fmt.Errorf("failed to marshal params: %w", apiErr)
	}

	resp, err := s.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE family_id=?", s.tableUser)),
		Parameters: params,
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return family, fmt.Errorf("failed to execute statement: %w", apiErr)
	}

	// fetch user IDs from family table first

	// userIds := []string{
	// 	"58214350-c071-704a-82b0-1b83f248d1bd",
	// 	"98a1c330-4051-701a-57c1-e8debd152f2b",
	// }

	// now fetch users matching all user IDs in family
	// filterEx := expression.Name("family_id").Equal(expression.Value(family_id))
	// expr, err := expression.NewBuilder().WithFilter(filterEx).Build()

	// if err != nil {
	// 	return family, fmt.Errorf("failed to build query expression: %w", err)
	// }

	// resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
	// 	TableName:                 aws.String(s.tableUser),
	// 	ExpressionAttributeNames:  expr.Names(),
	// 	ExpressionAttributeValues: expr.Values(),
	// 	// KeyConditionExpression:    expr.KeyCondition(),
	// 	FilterExpression: expr.Filter(),
	// 	ScanIndexForward: aws.Bool(false),
	// })

	// if err != nil {
	// 	apiErr := apierr.GetAwsError(err)
	// 	return family, apiErr
	// }

	if len(resp.Items) == 0 {
		logger.WithContext(ctx).WithField("family_id", family_id).Warnf("no family members found (family_id:%s)", family_id)
		return family, apierr.New(apierr.NotFound).WithError(fmt.Sprintf("family (family_id=%s)", family_id))
	}

	apiErr := apierr.New(fmt.Errorf("failed parsing users"))
	for idx, item := range resp.Items {
		user := models.User{}
		if err = attributevalue.UnmarshalMap(item, &user); err != nil {
			apiErr.AppendErrorf("user[%d] unmarshal error: %s", idx, err.Error())
			continue
		}

		if user.IsParent() {
			family.Parents[user.UserID] = models.NewFamilyUser(user)
		} else if user.IsChild() {
			family.Children[user.UserID] = models.NewFamilyUser(user)
		}
	}

	if len(apiErr.Errors()) > 0 {
		return family, apiErr
	}

	return family, nil
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
		// FilterExpression:          expr.Filter(),
		ScanIndexForward: aws.Bool(false),
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
