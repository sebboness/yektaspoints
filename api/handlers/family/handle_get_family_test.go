package family

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	"github.com/sebboness/yektaspoints/models"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errFail = errors.New("fail")

func Test_Controller_GetFamilyHandler(t *testing.T) {
	type state struct {
		missingFamilyID bool
		invalidUser     bool
		err             error
	}
	type want struct {
		err  string
		code int
	}
	type test struct {
		name string
		state
		want
	}

	cases := []test{
		{"happy path", state{}, want{"", 200}},
		{"fail - missing family_id", state{missingFamilyID: true}, want{"invalid input", http.StatusBadRequest}},
		{"fail - invalid user", state{invalidUser: true}, want{"access denied", http.StatusForbidden}},
		{"fail - internal server error", state{err: errFail}, want{"fail", http.StatusInternalServerError}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			familyDB := mocks.NewMockIFamilyStorage(t)

			ctrl := FamilyController{
				familyDB: familyDB,
			}

			family := models.Family{
				FamilyID: "456",
			}

			familyUsers := []models.FamilyUser{
				{
					FamilyID: "456",
					UserID:   "1",
				},
			}

			if c.state.invalidUser {
				familyUsers[0].UserID = "2"
			}

			if !c.state.missingFamilyID && !c.state.invalidUser {
				familyDB.EXPECT().GetFamilyUsers(mock.Anything, mock.Anything).Return(familyUsers, nil).Once()
				familyDB.EXPECT().GetFamilyMembersByUserIDs(mock.Anything, mock.Anything, mock.Anything).Return(family, c.state.err).Once()
			} else if !c.state.missingFamilyID {
				familyDB.EXPECT().GetFamilyUsers(mock.Anything, mock.Anything).Return(familyUsers, nil).Once()
			}

			ctx := context.Background()

			evt := events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{
						"claims": map[string]interface{}{
							"sub": "1",
						},
					},
				},
			}

			if c.state.missingFamilyID {
				evt.RequestContext.Authorizer = nil
			}

			ctx = handlers.PrepareAuthorizedContext(ctx, evt)

			endpoint := "/v1/family"
			if !c.state.missingFamilyID {
				endpoint += "?family_id=456"
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", endpoint, nil).WithContext(ctx)

			ctrl.GetFamilyHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			if c.want.code == 200 {
				assert.NotNil(t, result.Data)
				if result.Data != nil {
					famRes := result.Data.(map[string]any)["family"]
					assert.Equal(t, "456", famRes.(map[string]any)["family_id"])
				}
			}

			familyDB.AssertExpectations(t)
		})
	}
}

func Test_handleUserRegister(t *testing.T) {
	type state struct {
		isInvalidUser  bool
		getFamErr      error
		getFamUsersErr error
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
		{"fail - user mismatch", state{isInvalidUser: true}, want{"user is not part of family"}},
		{"fail - get family users error", state{getFamUsersErr: errFail}, want{"failed to get family users"}},
		{"fail - get family error", state{getFamErr: errFail}, want{"failed to get family"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()
			familyDB := mocks.NewMockIFamilyStorage(t)

			ctrl := FamilyController{
				familyDB: familyDB,
			}

			family := models.Family{
				FamilyID: "456",
			}

			familyUsers := []models.FamilyUser{
				{
					FamilyID: "456",
					UserID:   "1",
				},
				{
					FamilyID: "456",
					UserID:   "2",
				},
			}

			familyDB.EXPECT().GetFamilyUsers(mock.Anything, mock.Anything).Return(familyUsers, c.state.getFamUsersErr).Once()

			if c.state.getFamUsersErr == nil && !c.state.isInvalidUser {
				familyDB.EXPECT().GetFamilyMembersByUserIDs(mock.Anything, mock.Anything, mock.Anything).Return(family, c.state.getFamErr).Once()
			}

			req := &getFamilyHandlerRequest{
				FamilyID: "456",
				UserID:   "1",
			}

			if c.state.isInvalidUser {
				req.UserID = "3"
			}

			res, err := ctrl.handleGetFamily(ctx, req)

			tests.AssertError(t, err, c.want.err)
			if c.want.err == "" {
				assert.Equal(t, "456", res.Family.FamilyID)
			}

			familyDB.AssertExpectations(t)
		})
	}
}
