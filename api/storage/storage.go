package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util"
	"github.com/sebboness/yektaspoints/util/log"
)

type DynamoDbClient interface {
	ExecuteStatement(ctx context.Context, params *dynamodb.ExecuteStatementInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ExecuteStatementOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type DynamoDbStorage struct {
	client      DynamoDbClient
	tablePoints string
	tableUser   string
}

type Config struct {
	Env string
}

var logger = log.Get()

func NewDynamoDbStorage(cfg Config) (*DynamoDbStorage, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	dynamoClient := dynamodb.NewFromConfig(sdkConfig)

	return &DynamoDbStorage{
		client:      dynamoClient,
		tablePoints: fmt.Sprintf("mypoints-%s-points", strings.ToLower(cfg.Env)),
		tableUser:   fmt.Sprintf("mypoints-%s-user", strings.ToLower(cfg.Env)),
	}, nil
}

// dateFilterExpression builds a date filter expression based on the given date filter.
// If both from and to dates are given, returns a "Between" filter.
// If only from date is given, returns a "greater than or equal to" filter.
// If only to date is given, returns a "less than or equal to" filter.
func dateFilterExpression(name string, f models.DateFilter) expression.ConditionBuilder {

	filterEx := expression.ConditionBuilder{}
	nameExpr := expression.Name(name)

	if f.From != nil && f.To != nil {
		filterEx = nameExpr.Between(
			expression.Value(util.ToFormatted(*f.From)), expression.Value(util.ToFormatted(*f.To)))
	} else if f.From != nil {
		filterEx = nameExpr.GreaterThanEqual(expression.Value(util.ToFormatted(*f.From)))
	} else if f.To != nil {
		filterEx = nameExpr.LessThanEqual(expression.Value(util.ToFormatted(*f.To)))
	}

	return filterEx
}

func valueInListExpression[K comparable](name string, values []K) expression.ConditionBuilder {
	if len(values) == 0 {
		return expression.ConditionBuilder{}
	}

	nameExpr := expression.Name(name)
	in := []expression.OperandBuilder{}
	for _, v := range values {
		in = append(in, expression.Value(v))
	}

	filterEx := expression.ConditionBuilder{}

	if len(in) == 1 {
		filterEx = nameExpr.In(in[0])
	} else {
		filterEx = nameExpr.In(in[0], in[1:]...)
	}

	return filterEx
}
