package storage

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_IFamilyStorage_GetFamilyMembersByUserIDs(t *testing.T) {
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
		{"fail - unmarshal", state{failUnmarshal: true}, want{"failed parsing users: user[0] unmarshal error"}},
		{"fail - not found", state{itemNotFound: true}, want{"resource not found: family (family_id=456)"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output := &dynamodb.ExecuteStatementOutput{
				Items: []map[string]types.AttributeValue{
					{
						"family_id": &types.AttributeValueMemberS{Value: "123"},
						"user_id":   &types.AttributeValueMemberS{Value: "1"},
						"email":     &types.AttributeValueMemberN{Value: "john@info.co"},
						"name":      &types.AttributeValueMemberN{Value: "John"},
						"roles":     &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "parent"}}},
					},
					{
						"family_id": &types.AttributeValueMemberS{Value: "123"},
						"user_id":   &types.AttributeValueMemberS{Value: "2"},
						"email":     &types.AttributeValueMemberN{Value: "kid@info.co"},
						"name":      &types.AttributeValueMemberN{Value: "Kid"},
						"roles":     &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "child"}}},
					},
				},
			}

			if c.state.failUnmarshal {
				output.Items = []map[string]types.AttributeValue{
					{
						"roles": &types.AttributeValueMemberS{Value: "abc"},
					},
				}
			}

			if c.state.itemNotFound {
				output.Items = []map[string]types.AttributeValue{}
			}

			mockDynamoClient := mocks.NewMockDynamoDbClient(t)
			mockDynamoClient.EXPECT().ExecuteStatement(mock.Anything, mock.Anything).Return(output, c.state.errGetItem)

			s := DynamoDbStorage{
				client: mockDynamoClient,
			}

			res, err := s.GetFamilyMembersByUserIDs(context.Background(), "456", []string{"1", "2"})
			tests.AssertError(t, err, c.want.err)

			if c.want.err == "" {
				assert.Equal(t, "456", res.FamilyID)
				assert.Len(t, res.Parents, 1)
				assert.Len(t, res.Children, 1)
				assert.Equal(t, "John", res.Parents["1"].Name)
				assert.Equal(t, "john@info.co", res.Parents["1"].Email)
				assert.Equal(t, "1", res.Parents["1"].UserID)
				assert.Equal(t, "Kid", res.Children["2"].Name)
				assert.Equal(t, "kid@info.co", res.Children["2"].Email)
				assert.Equal(t, "2", res.Children["2"].UserID)
			}

			mockDynamoClient.AssertExpectations(t)
		})
	}
}

func Test_IFamilyStorage_GetFamilyUsers(t *testing.T) {
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
		{"fail - unmarshal", state{failUnmarshal: true}, want{"failed to unmarshal family users from query response"}},
		{"fail - not found", state{itemNotFound: true}, want{"resource not found: family users (family_id=456)"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output := &dynamodb.ExecuteStatementOutput{
				Items: []map[string]types.AttributeValue{
					{
						"family_id": &types.AttributeValueMemberS{Value: "456"},
						"user_id":   &types.AttributeValueMemberS{Value: "1"},
					},
					{
						"family_id": &types.AttributeValueMemberS{Value: "456"},
						"user_id":   &types.AttributeValueMemberS{Value: "2"},
					},
				},
			}

			if c.state.failUnmarshal {
				output.Items = []map[string]types.AttributeValue{
					{
						"family_id": &types.AttributeValueMemberB{Value: []byte("")},
					},
				}
			}

			if c.state.itemNotFound {
				output.Items = []map[string]types.AttributeValue{}
			}

			mockDynamoClient := mocks.NewMockDynamoDbClient(t)
			mockDynamoClient.EXPECT().ExecuteStatement(mock.Anything, mock.Anything).Return(output, c.state.errGetItem)

			s := DynamoDbStorage{
				client: mockDynamoClient,
			}

			res, err := s.GetFamilyUsers(context.Background(), "456")
			tests.AssertError(t, err, c.want.err)

			if c.want.err == "" {
				assert.Len(t, res, 2)
				assert.Equal(t, "456", res[0].FamilyID)
				assert.Equal(t, "1", res[0].UserID)
				assert.Equal(t, "456", res[1].FamilyID)
				assert.Equal(t, "2", res[1].UserID)
			}

			mockDynamoClient.AssertExpectations(t)
		})
	}
}

// Tests against real db

func TestReal_IFamilyStorage_GetFamilyUsers(t *testing.T) {
	t.Skip("Skip real test")

	type state struct {
		familyId string
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
		{"happy path", state{familyId: "d578baa6-535c-4587-b177-e9ddcabb6688"}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := NewDynamoDbStorage(Config{Env: "dev"})
			assert.Nil(t, err)

			res, err := s.GetFamilyUsers(context.Background(), c.state.familyId)
			tests.AssertError(t, err, c.want.err)

			if c.want.err == "" {
				assert.Len(t, res, 2)
			}
		})
	}
}

func TestReal_IFamilyStorage_GetFamilyByFamilyID(t *testing.T) {
	t.Skip("Skip real test")

	type state struct {
		familyId string
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
		{"happy path", state{familyId: "d578baa6-535c-4587-b177-e9ddcabb6688"}, want{}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := NewDynamoDbStorage(Config{Env: "dev"})
			assert.Nil(t, err)

			res, err := s.GetFamilyMembersByUserIDs(context.Background(), c.state.familyId, []string{
				"58214350-c071-704a-82b0-1b83f248d1bd",
				"98a1c330-4051-701a-57c1-e8debd152f2b",
			})
			tests.AssertError(t, err, c.want.err)

			if c.want.err == "" {
				assert.Equal(t, c.state.familyId, res.FamilyID)
			}
		})
	}
}
