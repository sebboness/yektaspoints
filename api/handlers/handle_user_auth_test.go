package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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

func Test_UserAuthHandler(t *testing.T) {
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

			req := &userAuthRequest{
				GrantType: auth.GrantTypePassword,
				Username:  "john",
				Password:  "123",
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockAuther := mocks.NewMockAuthController(t)

			authRes := auth.AuthResult{Token: "abc"}

			if !c.state.invalidBody {
				mockAuther.EXPECT().Authenticate(mock.Anything, mock.Anything, mock.Anything).Return(authRes, c.state.errAuth).Once()
			} else {
				evtBodyStr = `{"user_id":`
			}

			ctrl := LambdaController{
				auth: mockAuther,
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/auth/token", bytes.NewReader([]byte(evtBodyStr)))

			ctrl.UserAuthHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			body, err := io.ReadAll(w.Body)
			assert.Nil(t, err, "reading response body should have no error")

			if c.want.err == "" {
				assert.Contains(t, string(body), `token":"abc"`)
			} else {
				assert.Contains(t, string(body), c.want.err)
			}

			mockAuther.AssertExpectations(t)
		})
	}
}

func Test_handleUserAuth(t *testing.T) {
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

			ctrl := LambdaController{
				auth: mockAuther,
			}

			req := &userAuthRequest{
				Username: "123",
			}

			authRes := auth.AuthResult{Token: "abc"}

			if c.state.isPwFlow {
				req.GrantType = auth.GrantTypePassword

				if !c.state.hasValidationErr {
					req.Password = "456"
					mockAuther.EXPECT().Authenticate(
						mock.Anything, mock.Anything, mock.Anything).Return(
						authRes, c.state.authErr)
				}
			}
			if c.state.isRtFlow {
				req.GrantType = auth.GrantTypeRefreshToken

				if !c.state.hasValidationErr {
					req.RefreshToken = "456"
					mockAuther.EXPECT().RefreshToken(
						mock.Anything, mock.Anything, mock.Anything).Return(
						authRes, c.state.authErr)
				}
			}

			ctx := context.Background()
			res, err := ctrl.handleUserAuth(ctx, req)

			tests.AssertError(t, err, c.want.err)
			if err == nil {
				assert.Equal(t, res.Token, "abc")
			}

			mockAuther.AssertExpectations(t)
		})
	}
}

func Test_validateUserAuth(t *testing.T) {
	type state struct {
		grantType    string
		username     string
		password     string
		refreshToken string
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
		{"happy path granttype password", state{grantType: auth.GrantTypePassword, username: "123", password: "456"}, want{}},
		{"happy path granttype refreshtoken", state{grantType: auth.GrantTypeRefreshToken, username: "123", refreshToken: "456"}, want{}},
		{"fail granttype unsupported", state{grantType: "client_credentials"}, want{"unsupported grant_type \"client_credentials\""}},
		{"fail granttype password - missing username", state{grantType: auth.GrantTypePassword, password: "456"}, want{"missing username"}},
		{"fail granttype password - missing password", state{grantType: auth.GrantTypePassword, username: "123"}, want{"missing password"}},
		{"fail granttype refreshtoken - missing username", state{grantType: auth.GrantTypeRefreshToken, refreshToken: "456"}, want{"missing username"}},
		{"fail granttype refreshtoken - missing refreshtoken", state{grantType: auth.GrantTypeRefreshToken, username: "123"}, want{"missing refresh_token"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &userAuthRequest{
				GrantType:    c.state.grantType,
				Username:     c.state.username,
				Password:     c.state.password,
				RefreshToken: c.state.refreshToken,
			}

			err := validateUserAuth(req)
			tests.AssertError(t, err, c.want.err)
		})
	}
}
