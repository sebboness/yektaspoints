package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/sebboness/yektaspoints/handlers"
	authmocks "github.com/sebboness/yektaspoints/mocks/auth"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"

	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserRegisterConfirmHandler(t *testing.T) {
	type state struct {
		invalidBody bool
		updateErr   error
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
		{"fail - invalid body", state{invalidBody: true}, want{"failed to unmarshal json body", 400}},
		{"fail - validation error", state{updateErr: apierr.New(apierr.InvalidInput)}, want{"invalid input", 400}},
		{"fail - internal server error", state{updateErr: errors.New("fail")}, want{"fail", 500}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

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

			ctx = handlers.PrepareAuthorizedContext(ctx, evt)

			req := &userRegisterConfirmRequest{
				Code:     "123456",
				UserID:   "1",
				Username: "john",
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockAuther := authmocks.NewMockAuthController(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			if !c.state.invalidBody {
				mockAuther.EXPECT().ConfirmRegistration(mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mockUserDB.EXPECT().UpdateUserStatus(mock.Anything, mock.Anything, mock.Anything).Return(c.state.updateErr).Once()
			} else {
				evtBodyStr = `{"":`
			}

			ctrl := UserController{
				auth:   mockAuther,
				userDB: mockUserDB,
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", bytes.NewReader([]byte(evtBodyStr))).WithContext(ctx)

			ctrl.UserRegisterConfirmHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			mockAuther.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_handleUserRegisterConfirm(t *testing.T) {
	type state struct {
		hasValidationErr bool
		regErr           error
		updateErr        error
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
		{"happy path - password flow", state{}, want{}},
		{"fail - invalid input", state{hasValidationErr: true}, want{"failed to validate request"}},
		{"fail - register error", state{regErr: errFail}, want{"failed to confirm user registration for 'john'"}},
		{"fail - update error", state{updateErr: errFail}, want{"failed to update user status to active for 'john'"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()
			mockAuther := authmocks.NewMockAuthController(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			ctrl := UserController{
				auth:   mockAuther,
				userDB: mockUserDB,
			}

			if !c.state.hasValidationErr {
				mockAuther.EXPECT().ConfirmRegistration(mock.Anything, mock.Anything, mock.Anything).Return(c.state.regErr).Once()
			}
			if !c.state.hasValidationErr && c.state.regErr == nil {
				mockUserDB.EXPECT().UpdateUserStatus(mock.Anything, mock.Anything, mock.Anything).Return(c.state.updateErr).Once()
			}

			req := &userRegisterConfirmRequest{
				Code:     "123456",
				UserID:   "1",
				Username: "john",
			}

			if c.state.hasValidationErr {
				req.Code = ""
			}

			err := ctrl.handleUserRegisterConfirm(ctx, req)

			tests.AssertError(t, err, c.want.err)

			mockAuther.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_validateUserRegisterConfirm(t *testing.T) {
	type state struct {
		code   string
		uname  string
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
		{"happy path", state{uname: "john", userId: "1", code: "John"}, want{}},
		{"missing user_id", state{uname: "1", userId: "", code: "John"}, want{"missing user_id"}},
		{"missing username", state{uname: "", userId: "1", code: "John"}, want{"missing username"}},
		{"missing code", state{uname: "john", userId: "1", code: ""}, want{"missing code"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &userRegisterConfirmRequest{
				Code:     c.state.code,
				UserID:   c.state.userId,
				Username: c.state.uname,
			}

			err := validateUserRegisterConfirm(req)
			tests.AssertError(t, err, c.want.err)
		})
	}
}
