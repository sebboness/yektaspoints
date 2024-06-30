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

type IPointsStorage interface {
	GetLatestBalance(ctx context.Context, userId string) (models.PointBalance, error)
	GetPointByID(ctx context.Context, userId, id string) (models.Point, error)
	GetPointsByUserID(ctx context.Context, userId string, filters models.QueryPointsFilter) ([]models.Point, error)
	SavePoint(ctx context.Context, point models.Point) error
}

func (s *DynamoDbStorage) GetLatestBalance(ctx context.Context, userId string) (models.PointBalance, error) {
	balance := models.PointBalance{}

	keyEx := expression.Key("user_id").Equal(expression.Value(userId))
	statusExpr := expression.Name("status").Equal(expression.Value(models.PointStatusSettled))
	balanceExpr := expression.Name("balance").GreaterThan(expression.Value(aws.Int32(0)))

	exprBuilder := expression.NewBuilder().WithKeyCondition(keyEx)
	exprBuilder = exprBuilder.WithFilter(statusExpr)
	exprBuilder = exprBuilder.WithFilter(balanceExpr)
	exprBuilder = exprBuilder.WithProjection(selectAttributesExpression([]string{"user_id", "id", "balance"}))

	expr, err := exprBuilder.Build()
	if err != nil {
		return balance, fmt.Errorf("failed to build expression for query: %w", err)
	}

	resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(s.tablePoints),
		IndexName:                 aws.String("updated_on-index"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		Limit:                     aws.Int32(1), // Just need top 1 item
		ProjectionExpression:      expr.Projection(),
		ScanIndexForward:          aws.Bool(false), // Sort descending
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return balance, apiErr
	}

	// if there are no items, return a 0-balance record, and include the user id
	if len(resp.Items) == 0 {
		balance.UserID = userId
		return balance, nil
	}

	err = attributevalue.UnmarshalMap(resp.Items[0], &balance)
	if err != nil {
		return balance, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return balance, nil
}

func (s *DynamoDbStorage) GetPointByID(ctx context.Context, userId, pointId string) (models.Point, error) {
	point := models.Point{}

	userIdAttr, err := attributevalue.Marshal(userId)
	if err != nil {
		return point, fmt.Errorf("failed to marshal user_id key: %w", err)
	}
	pointIdAttr, err := attributevalue.Marshal(pointId)
	if err != nil {
		return point, fmt.Errorf("failed to marshal point_id key: %w", err)
	}

	key := map[string]types.AttributeValue{
		"user_id": userIdAttr,
		"id":      pointIdAttr,
	}

	resp, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tablePoints),
		Key:       key,
	})

	if err != nil {
		logger.WithContext(ctx).AddFields(map[string]any{"user_id": userId, "point_id": pointId})
		apiErr := apierr.GetAwsError(err)
		return point, apiErr
	}

	if resp.Item == nil {
		logger.WithContext(ctx).AddFields(map[string]any{"user_id": userId, "point_id": pointId}).Errorf("item (id:%s) not found", pointId)
		return point, apierr.New(apierr.NotFound).WithError(fmt.Sprintf("point (id=%s)", pointId))
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
	var filterExpr expression.ConditionBuilder

	// created_on filter
	if filters.CreatedOn.IsSet() {
		filterExpr = dateFilterExpression("created_on", filters.CreatedOn)
	}

	// updated_on filter
	if filters.UpdatedOn.IsSet() {
		updatedOnFilterExpr := dateFilterKeyExpression("updated_on", filters.UpdatedOn)
		keyEx = expression.KeyAnd(keyEx, updatedOnFilterExpr)
	}

	exprBuilder := expression.NewBuilder().WithKeyCondition(keyEx)

	// status filter
	if len(filters.Statuses) > 0 {
		statusFilter := valueInListExpression("status", filters.Statuses)
		if filterExpr.IsSet() {
			filterExpr = filterExpr.And(statusFilter)
		} else {
			filterExpr = statusFilter
		}
	}

	// type filter
	if len(filters.Types) > 0 {
		typeFilter := valueInListExpression("request.type", filters.Types)
		if filterExpr.IsSet() {
			filterExpr = filterExpr.And(typeFilter)
		} else {
			filterExpr = typeFilter
		}
	}

	if filterExpr.IsSet() {
		exprBuilder = exprBuilder.WithFilter(filterExpr)
	}

	if len(filters.Attributes) > 0 {
		exprBuilder.WithProjection(selectAttributesExpression(filters.Attributes))
	}

	expr, err := exprBuilder.Build()
	if err != nil {
		return points, fmt.Errorf("failed to build expression for query: %w", err)
	}

	// setup query paginator
	queryPaginator := dynamodb.NewQueryPaginator(s.client, &dynamodb.QueryInput{
		TableName:                 aws.String(s.tablePoints),
		IndexName:                 aws.String("updated_on-index"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ScanIndexForward:          aws.Bool(false), // order by updated_on descending (latest first)
		ProjectionExpression:      expr.Projection(),
	})

	// fetch items from each page
	for queryPaginator.HasMorePages() {
		qctx := context.TODO()
		resp, err := queryPaginator.NextPage(qctx)

		if err != nil {
			apiErr := apierr.GetAwsError(err)
			return points, fmt.Errorf("failed to query next points page: %w", apiErr)
		} else {
			var queriedPoints []models.Point
			err = attributevalue.UnmarshalListOfMaps(resp.Items, &queriedPoints)
			if err != nil {
				return points, fmt.Errorf("failed to unmarshal points from query response: %w", err)
			} else {
				for _, p := range queriedPoints {
					p.ParseTimes()
					points = append(points, p)
				}
			}
		}
	}

	return points, nil
}

func (s *DynamoDbStorage) SavePoint(ctx context.Context, point models.Point) error {

	if err := s.validateNewPoint(point); err != nil {
		return err
	}

	item, err := attributevalue.MarshalMap(point)
	if err != nil {
		return fmt.Errorf("failed to marshal map from point: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tablePoints),
		Item:      item,
	})

	if err != nil {
		apiErr := apierr.GetAwsError(err)
		return apiErr
	}

	return nil
}

func (s *DynamoDbStorage) validateNewPoint(point models.Point) error {
	apierr := apierr.New(fmt.Errorf("%w: failed to validate request", apierr.InvalidInput))

	if point.UserID == "" {
		apierr.AppendError("missing user_id")
	}

	if point.ID == "" {
		apierr.AppendError("missing id")
	}

	if point.UpdatedOnStr == "" {
		apierr.AppendError("missing updated_on")
	}

	if len(apierr.Errors()) > 0 {
		return apierr
	}

	return nil
}
