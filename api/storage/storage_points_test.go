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
	"github.com/sebboness/yektaspoints/util/tests"
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

	throughputErr := &types.ProvisionedThroughputExceededException{}

	cases := []test{
		{"happy path", state{}, want{}},
		{"fail - get item", state{errGetItem: errFail}, want{"fail"}},
		{"fail - get item - exceeded throughput", state{errGetItem: throughputErr}, want{"ProvisionedThroughputExceededException"}},
		{"fail - unmarshal", state{failUnmarshal: true}, want{"failed to unmarshal item: unmarshal failed"}},
		{"fail - not found", state{itemNotFound: true}, want{"resource not found: point (id=123)"}},
	}

	for _, c := range cases {

		output := &dynamodb.QueryOutput{
			Items: []map[string]types.AttributeValue{
				{
					"id":      &types.AttributeValueMemberS{Value: "123"},
					"user_id": &types.AttributeValueMemberS{Value: "456"},
					"points":  &types.AttributeValueMemberN{Value: "100"},
				},
			},
		}

		if c.state.failUnmarshal {
			output.Items = []map[string]types.AttributeValue{
				{
					"points": &types.AttributeValueMemberS{Value: "abc"},
				},
			}
		}

		if c.state.itemNotFound {
			output.Items = []map[string]types.AttributeValue{}
		}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)
		mockDynamoClient.EXPECT().Query(mock.Anything, mock.Anything).Return(output, c.state.errGetItem)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		res, err := s.GetPointByID(context.Background(), "456", "123")
		tests.AssertError(t, err, c.want.err)

		if err == nil {
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

		to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		from := time.Date(2030, 12, 31, 11, 59, 59, 0, time.UTC)

		filter := models.QueryPointsFilter{
			CreatedOn: models.DateFilter{
				To:   &to,
				From: &from,
			},
			UpdatedOn: models.DateFilter{
				To:   &to,
				From: &from,
			},
			Statuses:   []models.PointStatus{models.PointStatusSettled},
			Types:      []models.PointRequestType{models.PointRequestTypeAdd},
			Attributes: []string{"id", "user_id", "status", "updated_on"},
		}

		res, err := s.GetPointsByUserID(context.Background(), "456", filter)
		tests.AssertError(t, err, c.want.err)

		if err == nil {
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
		missingID        bool
		missingUserID    bool
		missingUpdatedOn bool
		errSaveItem      error
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
		{"fail - validation error - missing id", state{missingID: true}, want{"missing id"}},
		{"fail - validation error - missing user_id", state{missingUserID: true}, want{"missing user_id"}},
		{"fail - validation error - missing updated_on", state{missingUpdatedOn: true}, want{"missing updated_on"}},
		{"fail - save point", state{errSaveItem: errFail}, want{"fail"}},
	}

	for _, c := range cases {

		output := &dynamodb.PutItemOutput{}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		point := models.Point{
			UserID:       "a",
			ID:           "1",
			UpdatedOnStr: util.ToFormatted(time.Now()),
		}

		hasValidationErr := false

		if c.state.missingID {
			point.ID = ""
			hasValidationErr = true
		}
		if c.state.missingUserID {
			point.UserID = ""
			hasValidationErr = true
		}
		if c.state.missingUpdatedOn {
			point.UpdatedOnStr = ""
			hasValidationErr = true
		}

		if !hasValidationErr {
			mockDynamoClient.EXPECT().PutItem(mock.Anything, mock.Anything, mock.Anything).Return(output, c.state.errSaveItem)
		}

		err := s.SavePoint(context.Background(), point)
		tests.AssertError(t, err, c.want.err)
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
		tests.AssertError(t, err, c.want.err)

		if err == nil {
			res.ParseTimes()
			assert.Empty(t, c.want.err)
			assert.Equal(t, c.state.id, res.ID)
			assert.Equal(t, c.state.userId, res.UserID)
			assert.Equal(t, 5, res.Points)
			assert.Equal(t, 5, res.Balance)
			assert.Equal(t, models.PointRequestTypeAdd, res.Request.Type)
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
		{"happy path", state{userId: "98a1c330-4051-701a-57c1-e8debd152f2b"}, want{}},
	}

	for _, c := range cases {
		s, err := NewDynamoDbStorage(Config{Env: "dev"})
		assert.Nil(t, err)

		from := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 3, 31, 59, 59, 0, 0, time.UTC)

		filter := models.QueryPointsFilter{
			UpdatedOn: models.DateFilter{
				From: &from,
				To:   &to,
			},
			Statuses: []models.PointStatus{models.PointStatusWaiting},
		}

		// simulate getting most recent points with waiting status
		resWaitingPoints, err := s.GetPointsByUserID(context.Background(), c.state.userId, filter)
		tests.AssertError(t, err, c.want.err)
		assert.NotEmpty(t, resWaitingPoints)

		filter = models.QueryPointsFilter{
			UpdatedOn: models.DateFilter{
				From: &from,
				To:   &to,
			},
			Statuses: []models.PointStatus{models.PointStatusSettled},
			Types:    []models.PointRequestType{models.PointRequestTypeAdd, models.PointRequestTypeSubtract},
		}

		// simulate getting most recent points
		resRecentPoints, err := s.GetPointsByUserID(context.Background(), c.state.userId, filter)
		tests.AssertError(t, err, c.want.err)
		assert.NotEmpty(t, resRecentPoints)

		filter = models.QueryPointsFilter{
			UpdatedOn: models.DateFilter{
				From: &from,
				To:   &to,
			},
			Statuses: []models.PointStatus{models.PointStatusSettled},
			Types:    []models.PointRequestType{models.PointRequestTypeCashout},
		}

		// simulate getting most recent points
		resRecentCashouts, err := s.GetPointsByUserID(context.Background(), c.state.userId, filter)
		tests.AssertError(t, err, c.want.err)
		assert.NotEmpty(t, resRecentCashouts)
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
		s, err := NewDynamoDbStorage(Config{Env: "dev"})
		assert.Nil(t, err)

		points := int(8)

		p := models.Point{
			ID:           ksuid.New().String(),
			UserID:       "d31b6627-cf66-4013-9e35-a46f0cb2e884",
			CreatedOnStr: util.ToFormattedUTC(time.Now()),
			UpdatedOnStr: util.ToFormattedUTC(time.Now()),
			Points:       3,
			Balance:      &points,
			Status:       models.PointStatusWaiting,
			Request: models.PointRequest{
				Type:   models.PointRequestTypeAdd,
				Reason: "I cleaned up my room",
			},
		}

		err = s.SavePoint(context.Background(), p)
		tests.AssertError(t, err, c.want.err)
	}
}
