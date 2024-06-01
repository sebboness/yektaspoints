package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type IFamilyStorage interface {
	GetFamilyMembersByUserIDs(ctx context.Context, family_id string, user_ids []string) (models.Family, error)
	GetFamilyUsers(ctx context.Context, family_id string) ([]models.FamilyUser, error)
	UserHasAccessToChild(ctx context.Context, family_id string, user_id string, child_id string) (bool, error)
	UserBelongsToFamily(ctx context.Context, user_id string, family_id string) (bool, error)
}

func (s *DynamoDbStorage) GetFamilyMembersByUserIDs(ctx context.Context, family_id string, user_ids []string) (models.Family, error) {
	family := models.Family{
		FamilyID: family_id,
		Parents:  map[string]models.FamilyMember{},
		Children: map[string]models.FamilyMember{},
	}

	paramValues := []interface{}{}
	paramMarks := []string{}
	for _, uid := range user_ids {
		paramMarks = append(paramMarks, "?")
		paramValues = append(paramValues, uid)
	}

	params, err := attributevalue.MarshalList(paramValues)
	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return family, fmt.Errorf("failed to marshal params: %w", apiErr)
	}

	resp, err := s.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE user_id IN [%s]", s.tableUser, strings.Join(paramMarks, ", "))),
		Parameters: params,
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return family, fmt.Errorf("failed to execute statement: %w", apiErr)
	}

	if len(resp.Items) == 0 {
		logger.WithContext(ctx).WithField("family_id", user_ids).Warnf("no users found (family_id:%s)", family_id)
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

func (s *DynamoDbStorage) GetFamilyUsers(ctx context.Context, family_id string) ([]models.FamilyUser, error) {
	familyUsers := []models.FamilyUser{}

	params, err := attributevalue.MarshalList([]interface{}{family_id})
	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return familyUsers, fmt.Errorf("failed to marshal params: %w", apiErr)
	}

	resp, err := s.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE family_id=?", s.tableFamilyUser)),
		Parameters: params,
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return familyUsers, fmt.Errorf("failed to execute statement: %w", apiErr)
	}

	if len(resp.Items) == 0 {
		logger.WithContext(ctx).WithField("family_id", family_id).Warnf("no family users found (family_id:%s)", family_id)
		return familyUsers, apierr.New(apierr.NotFound).WithError(fmt.Sprintf("family users (family_id=%s)", family_id))
	}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &familyUsers)
	if err != nil {
		return familyUsers, fmt.Errorf("failed to unmarshal family users from query response: %w", err)
	}

	return familyUsers, nil
}

func (s *DynamoDbStorage) UserBelongsToFamily(ctx context.Context, userId string, familyId string) (bool, error) {

	familyIdAttr, err := attributevalue.Marshal(familyId)
	if err != nil {
		return false, fmt.Errorf("failed to marshal family_id key: %w", err)
	}
	userIdAttr, err := attributevalue.Marshal(userId)
	if err != nil {
		return false, fmt.Errorf("failed to marshal user_id key: %w", err)
	}

	key := map[string]types.AttributeValue{
		"family_id": familyIdAttr,
		"user_id":   userIdAttr,
	}

	resp, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableFamilyUser),
		Key:       key,
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return false, apiErr
	}

	return resp.Item != nil, nil
}

func (s *DynamoDbStorage) UserHasAccessToChild(ctx context.Context, familyId string, userId string, childId string) (bool, error) {
	if userId == "" || childId == "" {
		return false, apierr.New(apierr.BadRequest).WithError("one or all user ids are empty")
	}
	if userId == childId {
		return true, nil
	}

	belongs, err := s.UserBelongsToFamily(ctx, userId, familyId)
	if err != nil {
		return false, err
	}
	if !belongs {
		return false, apierr.New(apierr.BadRequest).WithError(fmt.Sprintf("user id %v does not belong to family_id %v", userId, familyId))
	}

	belongs, err = s.UserBelongsToFamily(ctx, childId, familyId)
	if err != nil {
		return false, err
	}
	if !belongs {
		return false, apierr.New(apierr.BadRequest).WithError(fmt.Sprintf("child id %v does not belong to family_id %v", childId, familyId))
	}

	return true, nil
}
