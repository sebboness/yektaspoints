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

	apierr "github.com/sebboness/yektaspoints/util/error"
	"github.com/sebboness/yektaspoints/util/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserRegisterConfirmHandler(t *testing.T) {
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

			req := &userRegisterConfirmRequest{
				Username: "john",
				Code:     "123456",
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockAuther := mocks.NewMockAuthController(t)

			if !c.state.invalidBody {
				mockAuther.EXPECT().ConfirmRegistration(mock.Anything, mock.Anything, mock.Anything).Return(c.state.errAuth).Once()
			} else {
				evtBodyStr = `{"":`
			}

			ctrl := UserController{
				auth: mockAuther,
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", bytes.NewReader([]byte(evtBodyStr)))

			ctrl.UserRegisterConfirmHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			mockAuther.AssertExpectations(t)
		})
	}
}

func Test_handleUserRegisterConfirm(t *testing.T) {
	type state struct {
		hasValidationErr bool
		regErr           error
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctx := context.Background()
			mockAuther := mocks.NewMockAuthController(t)

			ctrl := UserController{
				auth: mockAuther,
			}

			if !c.state.hasValidationErr {
				mockAuther.EXPECT().ConfirmRegistration(mock.Anything, mock.Anything, mock.Anything).Return(c.state.regErr).Once()
			}

			req := &userRegisterConfirmRequest{
				Username: "john",
				Code:     "123456",
			}

			if c.state.hasValidationErr {
				req.Code = ""
			}

			err := ctrl.handleUserRegisterConfirm(ctx, req)

			tests.AssertError(t, err, c.want.err)
			mockAuther.AssertExpectations(t)
		})
	}
}

func Test_validateUserRegisterConfirm(t *testing.T) {
	type state struct {
		code  string
		uname string
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
		{"happy path", state{uname: "john", code: "John"}, want{}},
		{"missing username", state{uname: "", code: "John"}, want{"missing username"}},
		{"missing code", state{uname: "john", code: ""}, want{"missing code"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &userRegisterConfirmRequest{
				Username: c.state.uname,
				Code:     c.state.code,
			}

			err := validateUserRegisterConfirm(req)
			tests.AssertError(t, err, c.want.err)
		})
	}
}
