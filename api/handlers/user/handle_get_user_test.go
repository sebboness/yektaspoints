package user

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	handlerMocks "github.com/sebboness/yektaspoints/mocks/handlers"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"
	"github.com/sebboness/yektaspoints/models"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Controller_GetUserHandler(t *testing.T) {
	type state struct {
		hasNoAuth  bool
		getUserErr error
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
		{"fail - unauthorized", state{hasNoAuth: true}, want{"unauthorized", 401}},
		{"fail - not found", state{getUserErr: apierr.New(apierr.NotFound)}, want{"not found", 404}},
		{"fail - internal server error", state{getUserErr: errFail}, want{"failed to get user: fail", 500}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			mockAuthContext := handlerMocks.NewMockAuthContext(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			ctx := context.Background()

			authInfo := handlers.AuthorizerInfo{
				Claims: map[string]interface{}{
					"cognito:username": "john",
					"email":            "john@info.co",
					"email_verified":   "true",
					"name":             "John",
					"sub":              "1",
				},
			}

			if c.state.hasNoAuth {
				authInfo = handlers.AuthorizerInfo{}
			} else {
				mockUserDB.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(models.User{UserID: "1"}, c.state.getUserErr).Once()
			}

			mockAuthContext.EXPECT().GetAuthorizerInfo(mock.Anything).Return(authInfo)

			ctrl := UserController{
				BaseController: handlers.BaseController{
					AuthContext: mockAuthContext,
				},
				userDB: mockUserDB,
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", nil).WithContext(ctx)

			ctrl.GetUserHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			if c.want.code == 200 {
				assert.NotNil(t, result.Data)
				if result.Data != nil {
					assert.Equal(t, "1", result.Data.(map[string]any)["user_id"])
				}
			}

			mockAuthContext.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_Controller_handleGetUser(t *testing.T) {
	type state struct {
		getUserErr error
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
		{"fail - get user error", state{getUserErr: errFail}, want{"failed to get user"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()
			mockUserDB := mocks.NewMockIUserStorage(t)

			ctrl := UserController{
				userDB: mockUserDB,
			}

			user := models.User{
				UserID: "1",
			}

			mockUserDB.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(user, c.state.getUserErr).Once()

			res, err := ctrl.handleGetUser(ctx, "1")

			tests.AssertError(t, err, c.want.err)
			if c.want.err == "" {
				assert.Equal(t, "1", res.UserID)
			}

			mockUserDB.AssertExpectations(t)
		})
	}
}
