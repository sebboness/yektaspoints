package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.NewLogger("storage_points")

type IPointsStorage interface {
	GetPointByID(ctx context.Context, userId, id string) (models.Point, error)
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
		return point, fmt.Errorf("failed to get point: %w", err)
	}

	if resp.Item == nil {
		logger.WithContext(ctx).WithFields(map[string]any{"userId": userId, "id": id}).Warnf("item (id:%s) not found", id)
		return point, apierr.New(apierr.NotFound, fmt.Sprintf("point (id=%s)", id))
	}

	err = attributevalue.UnmarshalMap(resp.Item, &point)
	if err != nil {
		return point, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return point, nil
}
