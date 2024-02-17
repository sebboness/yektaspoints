package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
)

type IPointsStorage interface {
	GetPointByID(ctx context.Context, userId, id string) (models.Point, error)
	GetPointsByUserID(ctx context.Context, userId string, filters models.QueryPointsFilter) ([]models.Point, error)
	SavePoint(ctx context.Context, point models.Point) error
}

func (s *DynamoDbStorage) SavePoint(ctx context.Context, point models.Point) error {

	item, err := attributevalue.MarshalMap(point)
	if err != nil {
		return fmt.Errorf("failed to marshal map from point: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tablePoints),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to save point: %w", err)
	}

	return nil
}

func (s *DynamoDbStorage) GetPointByID(ctx context.Context, userId, id string) (models.Point, error) {
	point := models.Point{}

	idKey, err := attributevalue.Marshal(id)
	if err != nil {
		return point, fmt.Errorf("failed to marshal id key: %w", err)
	}

	userIdKey, err := attributevalue.Marshal(userId)
	if err != nil {
		return point, fmt.Errorf("failed to marshal user_id key: %w", err)
	}

	key := map[string]types.AttributeValue{
		"id":      idKey,
		"user_id": userIdKey,
	}

	resp, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(s.tablePoints),
	})

	if err != nil {
		var provExceeded *types.ProvisionedThroughputExceededException
		if errors.As(err, &provExceeded) {
			return point, fmt.Errorf("failed to get point: %w", err)
		}
		return point, fmt.Errorf("failed to get point: %w", err)
	}

	if resp.Item == nil {
		logger.WithContext(ctx).WithFields(map[string]any{"userId": userId, "id": id}).Warnf("item (id:%s) not found", id)
		return point, apierr.New(apierr.NotFound).WithError(fmt.Sprintf("point (id=%s)", id))
	}

	err = attributevalue.UnmarshalMap(resp.Item, &point)
	if err != nil {
		return point, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return point, nil
}

func (s *DynamoDbStorage) GetPointsByUserID(ctx context.Context, userId string, filters models.QueryPointsFilter) ([]models.Point, error) {
	points := []models.Point{}

	keyEx := expression.Key("user_id").Equal(expression.Value(userId))
	exprBuilder := expression.NewBuilder().WithKeyCondition(keyEx)
	var filterExpr expression.ConditionBuilder

	// Date filters
	if filters.RequestedOn.IsSet() {
		filterExpr = dateFilterExpression("requested_on", filters.RequestedOn)
	}

	// Statuses
	if len(filters.Statuses) > 0 {
		statusFilter := valueInListExpression("status_id", filters.Statuses)
		if filterExpr.IsSet() {
			filterExpr = filterExpr.And(statusFilter)
		} else {
			filterExpr = statusFilter
		}
	}

	// Types
	if len(filters.Types) > 0 {
		typeFilter := valueInListExpression("type", filters.Types)
		if filterExpr.IsSet() {
			filterExpr = filterExpr.And(typeFilter)
		} else {
			filterExpr = typeFilter
		}
	}

	if filterExpr.IsSet() {
		exprBuilder = exprBuilder.WithFilter(filterExpr)
	}

	expr, err := exprBuilder.Build()
	if err != nil {
		return points, fmt.Errorf("failed to build expression for query: %w", err)
	}

	// setup query paginator
	queryPaginator := dynamodb.NewQueryPaginator(s.client, &dynamodb.QueryInput{
		TableName:                 aws.String(s.tablePoints),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
	})

	// fetch items from each page
	for queryPaginator.HasMorePages() {
		qctx := context.TODO()
		resp, err := queryPaginator.NextPage(qctx)

		if err != nil {
			return points, fmt.Errorf("failed to query next points page: %w", err)
		} else {
			var point []models.Point
			err = attributevalue.UnmarshalListOfMaps(resp.Items, &point)
			if err != nil {
				return points, fmt.Errorf("failed to unmarshal points from query response: %w", err)
			} else {
				points = append(points, point...)
			}
		}
	}

	return points, nil
}
