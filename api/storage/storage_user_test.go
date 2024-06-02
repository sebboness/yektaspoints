package storage

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_IUserStorage_GetUserByID(t *testing.T) {
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
		{"fail - not found", state{itemNotFound: true}, want{"resource not found: user (user_id=456)"}},
	}

	for _, c := range cases {

		output := &dynamodb.QueryOutput{
			Items: []map[string]types.AttributeValue{
				{
					"family_id":  &types.AttributeValueMemberS{Value: "123"},
					"user_id":    &types.AttributeValueMemberS{Value: "456"},
					"username":   &types.AttributeValueMemberN{Value: "john"},
					"name":       &types.AttributeValueMemberN{Value: "John"},
					"created_on": &types.AttributeValueMemberN{Value: "2024-03-10T20:00:00.0000000Z"},
					"updated_on": &types.AttributeValueMemberN{Value: "2024-03-31T20:00:00.0000000Z"},
					"roles":      &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "parent"}}},
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
		mockDynamoClient.EXPECT().Query(mock.Anything, mock.Anything).Return(output, c.state.errGetItem)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		res, err := s.GetUserByID(context.Background(), "456")
		tests.AssertError(t, err, c.want.err)

		if err == nil {
			assert.Empty(t, c.want.err)
			assert.Equal(t, "456", res.UserID)
			assert.Equal(t, "john", res.Username)
			assert.Equal(t, "John", res.Name)
			assert.Equal(t, time.Date(2024, 3, 10, 20, 0, 0, 0, time.UTC), res.CreatedOn)
			assert.Equal(t, time.Date(2024, 3, 31, 20, 0, 0, 0, time.UTC), res.UpdatedOn)
			assert.Equal(t, []string{"parent"}, res.Roles)
		}

		mockDynamoClient.AssertExpectations(t)
	}
}

func Test_IUserStorage_SaveUser(t *testing.T) {
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
		{"fail - save user", state{errSaveItem: errFail}, want{"fail"}},
	}

	for _, c := range cases {

		output := &dynamodb.PutItemOutput{}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)
		mockDynamoClient.EXPECT().PutItem(mock.Anything, mock.Anything, mock.Anything).Return(output, c.state.errSaveItem)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		err := s.SaveUser(context.Background(), models.User{})
		tests.AssertError(t, err, c.want.err)
		mockDynamoClient.AssertExpectations(t)
	}
}

func Test_IUserStorage_UpdateUserFamily(t *testing.T) {
	type state struct {
		errUpdate error
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
		{"fail - update family", state{errUpdate: errFail}, want{"fail"}},
	}

	for _, c := range cases {

		output := &dynamodb.UpdateItemOutput{}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)
		mockDynamoClient.EXPECT().UpdateItem(mock.Anything, mock.Anything, mock.Anything).Return(output, c.state.errUpdate)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		err := s.UpdateUserFamily(context.Background(), UpdateUserFamilyRequest{})
		tests.AssertError(t, err, c.want.err)
		mockDynamoClient.AssertExpectations(t)
	}
}

func Test_IUserStorage_UpdateUserStatus(t *testing.T) {
	type state struct {
		errUpdate error
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
		{"fail - update status", state{errUpdate: errFail}, want{"fail"}},
	}

	for _, c := range cases {

		output := &dynamodb.UpdateItemOutput{}

		mockDynamoClient := mocks.NewMockDynamoDbClient(t)
		mockDynamoClient.EXPECT().UpdateItem(mock.Anything, mock.Anything, mock.Anything).Return(output, c.state.errUpdate)

		s := DynamoDbStorage{
			client: mockDynamoClient,
		}

		err := s.UpdateUserStatus(context.Background(), "1", models.UserStatusDeleted)
		tests.AssertError(t, err, c.want.err)
		mockDynamoClient.AssertExpectations(t)
	}
}

func Test_IUserStorage_ParentHasAccessToChild(t *testing.T) {
	type state struct {
		errGetParent      error
		errGetChild       error
		parentIdEmpty     bool
		childIdEmpty      bool
		parentOtherFamily bool
		childOtherFamily  bool
		parentNotParent   bool
		childNotChild     bool
	}
	type want struct {
		belongs bool
		err     string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{true, ""}},
		{"fail - parent not same family", state{parentOtherFamily: true}, want{false, ""}},
		{"fail - child not same family", state{childOtherFamily: true}, want{false, ""}},
		{"fail - parent id empty", state{parentIdEmpty: true}, want{false, "one or all user ids are empty"}},
		{"fail - child id empty", state{childIdEmpty: true}, want{false, "one or all user ids are empty"}},
		{"fail - get parent err", state{errGetParent: errFail}, want{false, "failed to get parent user record for parent-1"}},
		{"fail - get child err", state{errGetChild: errFail}, want{false, "failed to get child user record for child-1"}},
		{"fail - parent is not parent", state{parentNotParent: true}, want{false, "user parent-1 is not in role parent"}},
		{"fail - child is not child", state{childNotChild: true}, want{false, "user child-1 is not in role child"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			parentId := "parent-1"
			parentFamilyId := "family-1"
			childFamilyId := "family-1"

			childId := "child-1"
			parentsRole := "parent"
			childsRole := "child"

			if c.state.parentIdEmpty {
				parentId = ""
			}
			if c.state.parentOtherFamily {
				parentFamilyId = "blah"
			}
			if c.state.parentNotParent {
				parentsRole = "blah"
			}
			if c.state.childIdEmpty {
				childId = ""
			}
			if c.state.childOtherFamily {
				childFamilyId = "blah"
			}
			if c.state.childNotChild {
				childsRole = "blah"
			}

			outputGetParent := &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{
						"family_ids": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: parentFamilyId}}},
						"user_id":    &types.AttributeValueMemberS{Value: "parent-1"},
						"roles":      &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: parentsRole}}},
					},
				},
			}

			outputGetChild := &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{
						"family_ids": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: childFamilyId}}},
						"user_id":    &types.AttributeValueMemberS{Value: "child-1"},
						"roles":      &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: childsRole}}},
					},
				},
			}

			mockDynamoClient := mocks.NewMockDynamoDbClient(t)

			if !c.state.parentIdEmpty && !c.state.childIdEmpty {
				mockDynamoClient.EXPECT().Query(mock.Anything, mock.Anything).Return(outputGetParent, c.state.errGetParent).Once()

				if c.state.errGetParent == nil {
					mockDynamoClient.EXPECT().Query(mock.Anything, mock.Anything).Return(outputGetChild, c.state.errGetChild).Once()
				}
			}

			s := DynamoDbStorage{
				client: mockDynamoClient,
			}

			res, err := s.ParentHasAccessToChild(context.Background(), parentId, childId)
			tests.AssertError(t, err, c.want.err)

			if c.want.err == "" {
				assert.Equal(t, c.want.belongs, res)
			}

			mockDynamoClient.AssertExpectations(t)
		})
	}
}

// Tests below test against real db
// Comment out t.Skip() to run

func TestReal_IUserStorage_GetUserByID(t *testing.T) {
	t.Skip("Skip real test")

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
		{"happy path - exists", state{userId: "98a1c330-4051-701a-57c1-e8debd152f2b"}, want{""}},
		{"happy path - does not exist", state{userId: "123"}, want{"resource not found: user"}},
	}

	for _, c := range cases {
		s, err := NewDynamoDbStorage(Config{Env: "dev"})
		assert.Nil(t, err)

		user, err := s.GetUserByID(context.Background(), c.state.userId)
		tests.AssertError(t, err, c.want.err)

		if err == nil {
			assert.NotNil(t, user)
		} else {
			assert.Empty(t, user.UserID)
		}
	}
}

func TestReal_IUserStorage_ParentHasAccessToChild(t *testing.T) {
	t.Skip("Skip real test")

	type state struct {
		parentId string
		childId  string
	}
	type want struct {
		hasAccess bool
		err       string
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{
			"happy path - has access",
			state{
				childId:  "98a1c330-4051-701a-57c1-e8debd152f2b",
				parentId: "58214350-c071-704a-82b0-1b83f248d1bd",
			},
			want{true, ""},
		},
	}

	for _, c := range cases {
		s, err := NewDynamoDbStorage(Config{Env: "dev"})
		assert.Nil(t, err)

		hasAccess, err := s.ParentHasAccessToChild(context.Background(), c.state.parentId, c.state.childId)
		tests.AssertError(t, err, c.want.err)

		if err == nil {
			assert.Equal(t, c.want.hasAccess, hasAccess)
		}
	}
}

func TestReal_IUserStorage_UpdateUserStatus(t *testing.T) {
	t.Skip("Skip real test")

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
		{"happy path", state{userId: "58214350-c071-704a-82b0-1b83f248d1bd"}, want{}},
	}

	for _, c := range cases {
		s, err := NewDynamoDbStorage(Config{Env: "dev"})
		assert.Nil(t, err)

		err = s.UpdateUserStatus(context.Background(), c.state.userId, models.UserStatusActive)
		tests.AssertError(t, err, c.want.err)
	}
}
