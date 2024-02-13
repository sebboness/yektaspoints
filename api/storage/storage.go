package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sebboness/yektaspoints/util/log"
)

type DynamoDbClient interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type DynamoDbStorage struct {
	client      DynamoDbClient
	tablePoints string
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
		tablePoints: fmt.Sprintf("yektaspoints-%s-points", strings.ToLower(cfg.Env)),
	}, nil
}
