package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util"
	"github.com/sebboness/yektaspoints/util/env"
	"github.com/sebboness/yektaspoints/util/log"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errFail = errors.New("fail")

func Test_DynamoDbStorage_GetPointByID(t *testing.T) {
	type state struct {
		errGetItem    error
		failUnmarshal bool
		itemNotFound  bool
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
		{"fail - get item", state{errGetItem: errFail}, want{"failed to get point: fail"}},
		{"fail - unmarshal", state{failUnmarshal: true}, want{"failed to unmarshal item: unmarshal failed"}},
		{"fail - not found", state{itemNotFound: true}, want{"resource not found: point (id=123)"}},
	}

	for _, c := range cases {

		output := &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"id":      &types.AttributeValueMemberS{Value: "123"},
				"user_id": &types.AttributeValueMemberS{Value: "456"},
				"points":  &types.AttributeValueMemberN{Value: "100"},
			},
		}

		if c.state.failUnmarshal {
			output.Item = map[string]types.AttributeValue{
				"points": &types.AttributeValueMemberS{Value: "abc"},
			}
		}

		if c.state.itemNotFound {
			output.Item = nil
		}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)
		mockDynamoClient.EXPECT().GetItem(mock.Anything, mock.Anything, mock.Anything).Return(output, c.state.errGetItem)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		res, err := s.GetPointByID(context.Background(), "456", "123")
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		} else {
			assert.Empty(t, c.want.err)
			assert.Equal(t, "123", res.ID)
			assert.Equal(t, "456", res.UserID)
			assert.Equal(t, 100, res.Points)
		}

		mockDynamoClient.AssertExpectations(t)
	}
}

func Test_DynamoDbStorage_GetPointsByUserID(t *testing.T) {
	type state struct {
		errQuery      error
		failUnmarshal bool
		itemNotFound  bool
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
		{"fail - query", state{errQuery: errFail}, want{"failed to query next points page: fail"}},
		{"fail - unmarshal", state{failUnmarshal: true}, want{"failed to unmarshal points from query response: unmarshal failed"}},
	}

	for _, c := range cases {

		output := &dynamodb.QueryOutput{
			Items: []map[string]types.AttributeValue{
				{
					"id":      &types.AttributeValueMemberS{Value: "1"},
					"user_id": &types.AttributeValueMemberS{Value: "a"},
					"points":  &types.AttributeValueMemberN{Value: "7"},
				},
				{
					"id":      &types.AttributeValueMemberS{Value: "2"},
					"user_id": &types.AttributeValueMemberS{Value: "b"},
					"points":  &types.AttributeValueMemberN{Value: "9"},
				},
			},
		}

		if c.state.failUnmarshal {
			output.Items = []map[string]types.AttributeValue{
				{
					"points": &types.AttributeValueMemberS{Value: "xyz"},
				},
			}
		}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)
		mockDynamoClient.EXPECT().Query(mock.Anything, mock.Anything, mock.Anything).Return(output, c.state.errQuery)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		filter := models.QueryPointsFilter{
			Statuses: []models.PointStatus{models.PointStatusApproved},
			Types:    []models.PointType{models.PointTypeAdd},
		}

		res, err := s.GetPointsByUserID(context.Background(), "456", filter)
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		} else {
			assert.Len(t, res, 2)
			assert.Equal(t, res[0].ID, "1")
			assert.Equal(t, res[0].UserID, "a")
			assert.Equal(t, res[0].Points, 7)
			assert.Equal(t, res[1].ID, "2")
			assert.Equal(t, res[1].UserID, "b")
			assert.Equal(t, res[1].Points, 9)
		}

		mockDynamoClient.AssertExpectations(t)
	}
}

func Test_DynamoDbStorage_SavePoint(t *testing.T) {
	type state struct {
		errSaveItem error
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
		{"fail - save point", state{errSaveItem: errFail}, want{"failed to save point: fail"}},
	}

	for _, c := range cases {

		output := &dynamodb.PutItemOutput{}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)
		mockDynamoClient.EXPECT().PutItem(mock.Anything, mock.Anything, mock.Anything).Return(output, c.state.errSaveItem)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		err := s.SavePoint(context.Background(), models.Point{})
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		}

		mockDynamoClient.AssertExpectations(t)
	}
}

// Unit tests against real dev environment
// These tests should be skipped unless debugging with real services

func TestReal_DynamoDbStorage_GetPointByID(t *testing.T) {
	t.Skip("Skip real test")

	type state struct {
		id     string
		userId string
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{id: "af9bfcd5-c708-4158-a2ef-f33b57a86fc9", userId: "d31b6627-cf66-4013-9e35-a46f0cb2e884"}, want{}},
	}

	for _, c := range cases {
		s, err := NewDynamoDbStorage(Config{Env: env.GetEnv("ENV")})
		assert.Nil(t, err)

		res, err := s.GetPointByID(context.Background(), c.state.userId, c.state.id)
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		} else {
			res.ParseTimes()
			assert.Empty(t, c.want.err)
			assert.Equal(t, c.state.id, res.ID)
			assert.Equal(t, c.state.userId, res.UserID)
			assert.Equal(t, 5, res.Points)
			assert.Equal(t, 5, res.Balance)
			assert.Equal(t, 50, res.BalancePoints)
			assert.Equal(t, models.PointTypeAdd, res.Type)
		}
	}
}

func TestReal_DynamoDbStorage_GetPointsByUserID(t *testing.T) {
	// t.Skip("Skip real test")

	type state struct {
		userId string
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{userId: "d31b6627-cf66-4013-9e35-a46f0cb2e884"}, want{}},
	}

	for _, c := range cases {
		s, err := NewDynamoDbStorage(Config{Env: env.GetEnv("ENV")})
		assert.Nil(t, err)

		from := time.Date(2024, 2, 8, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 2, 11, 59, 59, 0, 0, time.UTC)

		filter := models.QueryPointsFilter{
			RequestedOn: models.DateFilter{
				From: &from,
				To:   &to,
			},
			Statuses: []models.PointStatus{models.PointStatusApproved},
			Types:    []models.PointType{models.PointTypeCashout},
		}

		res, err := s.GetPointsByUserID(context.Background(), c.state.userId, filter)

		if err != nil {
			log.Get().Errorf("%v", err)
			assert.Contains(t, err.Error(), c.want.err)
		} else {
			for i, d := range res {
				log.Get().WithField("data", d).Infof("[%v]points", i)
			}
			assert.NotEmpty(t, res)
		}
	}
}

func TestReal_DynamoDbStorage_SavePoint(t *testing.T) {
	t.Skip("Skip real test")

	type state struct {
	}
	type want struct {
		err string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{}},
	}

	for _, c := range cases {
		s, err := NewDynamoDbStorage(Config{Env: env.GetEnv("ENV")})
		assert.Nil(t, err)

		p := models.Point{
			ID:             ksuid.New().String(),
			UserID:         "d31b6627-cf66-4013-9e35-a46f0cb2e884",
			RequestedOnStr: util.ToFormattedUTC(time.Now()),
			Points:         3,
			BalancePoints:  8,
			Balance:        8,
			StatusID:       models.PointStatusRequested,
			Type:           models.PointTypeAdd,
			Reason:         "I cleaned up my room",
		}

		err = s.SavePoint(context.Background(), p)
		if err != nil {
			assert.Contains(t, err.Error(), c.want.err)
		} else {
			assert.Nil(t, err)
		}
	}
}
