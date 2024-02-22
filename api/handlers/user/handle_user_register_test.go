package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mocks "github.com/sebboness/yektaspoints/mocks/auth"

	"github.com/sebboness/yektaspoints/util/auth"
	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errFail = errors.New("fail")

func Test_UserRegisterHandler(t *testing.T) {
	type state struct {
		invalidBody bool
		errAuth     error
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
		{"fail - validation error", state{errAuth: apierr.New(apierr.InvalidInput)}, want{"invalid input", 400}},
		{"fail - internal server error", state{errAuth: errors.New("fail")}, want{"fail", 500}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			req := &userRegisterRequest{
				Username: "john",
				Password: "123",
				Email:    "john@site.com",
				Name:     "John",
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockAuther := mocks.NewMockAuthController(t)

			authRes := auth.AuthResult{AccessToken: "abc"}

			if !c.state.invalidBody {
				mockAuther.EXPECT().Authenticate(mock.Anything, mock.Anything, mock.Anything).Return(authRes, c.state.errAuth).Once()
			} else {
				evtBodyStr = `{"user_id":`
			}

			ctrl := UserController{
				auth: mockAuther,
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", bytes.NewReader([]byte(evtBodyStr)))

			ctrl.UserRegisterHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)

			if result.IsSuccess() {
				assert.Empty(t, c.want.err)
			} else {
				assert.Contains(t, result.Errors, c.want.err)
			}

			mockAuther.AssertExpectations(t)
		})
	}
}

func Test_handleUserRegister(t *testing.T) {
	type state struct {
		isPwFlow         bool
		isRtFlow         bool
		hasValidationErr bool
		authErr          error
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
		{"happy path - password flow", state{isPwFlow: true}, want{}},
		{"happy path - refresh token flow", state{isRtFlow: true}, want{}},
		{"fail - invalid input", state{isPwFlow: true, hasValidationErr: true}, want{"failed to validate request"}},
		{"fail - password flow", state{isPwFlow: true, authErr: errFail}, want{"failed to authenticate"}},
		{"fail - refresh token flow", state{isRtFlow: true, authErr: errFail}, want{"failed to refresh token"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			mockAuther := mocks.NewMockAuthController(t)

			ctrl := UserController{
				auth: mockAuther,
			}

			req := &userRegisterRequest{
				Username: "123",
			}

			ctx := context.Background()
			res, err := ctrl.handleUserRegister(ctx, req)

			tests.AssertError(t, err, c.want.err)
			if err == nil {
				assert.Equal(t, "123", res.Username)
			}

			mockAuther.AssertExpectations(t)
		})
	}
}

func Test_validateUserRegister(t *testing.T) {
	type state struct {
		name      string
		email     string
		username  string
		password  string
		cpassword string
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
		{"happy path granttype password", state{email: auth.GrantTypePassword, username: "123", password: "456"}, want{}},
		{"happy path granttype refreshtoken", state{email: auth.GrantTypeRefreshToken, username: "123", name: "456"}, want{}},
		{"fail granttype unsupported", state{email: "client_credentials"}, want{"unsupported grant_type \"client_credentials\""}},
		{"fail granttype password - missing username", state{email: auth.GrantTypePassword, password: "456"}, want{"missing username"}},
		{"fail granttype password - missing password", state{email: auth.GrantTypePassword, username: "123"}, want{"missing password"}},
		{"fail granttype refreshtoken - missing username", state{email: auth.GrantTypeRefreshToken, name: "456"}, want{"missing username"}},
		{"fail granttype refreshtoken - missing refreshtoken", state{email: auth.GrantTypeRefreshToken, username: "123"}, want{"missing refresh_token"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &userRegisterRequest{
				Username:        c.state.username,
				Password:        c.state.password,
				ConfirmPassword: c.state.cpassword,
				Email:           c.state.email,
				Name:            c.state.name,
			}

			err := validateUserRegister(req)
			tests.AssertError(t, err, c.want.err)
		})
	}
}
