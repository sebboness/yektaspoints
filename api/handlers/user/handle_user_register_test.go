package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	authmocks "github.com/sebboness/yektaspoints/mocks/auth"
	mocks "github.com/sebboness/yektaspoints/mocks/storage"

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
		errSave     error
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
		{"happy path", state{}, want{"", 201}},
		{"fail - invalid body", state{invalidBody: true}, want{"failed to unmarshal json body", 400}},
		{"fail - validation error", state{errSave: apierr.New(apierr.InvalidInput)}, want{"invalid input", 400}},
		{"fail - internal server error", state{errSave: errors.New("fail")}, want{"fail", 500}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			req := &userRegisterRequest{
				Username:        "john",
				Password:        "Test123!",
				ConfirmPassword: "Test123!",
				Email:           "john@info.co",
				Name:            "John",
			}

			evtBody, _ := json.Marshal(req)
			evtBodyStr := string(evtBody)

			mockAuther := authmocks.NewMockAuthController(t)
			mockUserDB := mocks.NewMockIUserStorage(t)

			authRes := auth.UserRegisterResult{Username: "john"}

			if !c.state.invalidBody {
				mockAuther.EXPECT().Register(mock.Anything, mock.Anything).Return(authRes, nil).Once()
				mockUserDB.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(c.state.errSave).Once()
			} else {
				evtBodyStr = `{"user_id":`
			}

			ctrl := UserController{
				auth:   mockAuther,
				userDB: mockUserDB,
			}

			w := httptest.NewRecorder()
			cgin, _ := gin.CreateTestContext(w)
			cgin.Request = httptest.NewRequest("GET", "/", bytes.NewReader([]byte(evtBodyStr)))

			ctrl.UserRegisterHandler(cgin)

			assert.Equal(t, c.want.code, w.Code)
			result := tests.AssertResult(t, w.Body)
			tests.AssertResultError(t, result, c.want.err)

			mockAuther.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_handleUserRegister(t *testing.T) {
	type state struct {
		hasValidationErr bool
		regErr           error
		saveErr          error
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
		{"fail - register error", state{regErr: errFail}, want{"failed to register user 'john'"}},
		{"fail - save error", state{saveErr: errFail}, want{"failed to save new user 'john'"}},
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

			regResult := auth.UserRegisterResult{Username: "john"}
			if c.state.regErr != nil {
				regResult = auth.UserRegisterResult{}
			}

			if !c.state.hasValidationErr {
				mockAuther.EXPECT().Register(mock.Anything, mock.Anything).Return(regResult, c.state.regErr).Once()
			}
			if !c.state.hasValidationErr && c.state.regErr == nil {
				mockUserDB.EXPECT().SaveUser(mock.Anything, mock.Anything).Return(c.state.saveErr).Once()
			}

			req := &userRegisterRequest{
				Username:        "john",
				Password:        "Test123!",
				ConfirmPassword: "Test123!",
				Email:           "john@info.co",
				Name:            "John",
			}

			if c.state.hasValidationErr {
				req.Email = "blah"
			}

			res, err := ctrl.handleUserRegister(ctx, req)

			tests.AssertError(t, err, c.want.err)
			if err == nil {
				assert.Equal(t, req.Username, res.Username)
			}

			mockAuther.AssertExpectations(t)
			mockUserDB.AssertExpectations(t)
		})
	}
}

func Test_validateUserRegister(t *testing.T) {
	type state struct {
		name  string
		email string
		uname string
		pass  string
		cpass string
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
		{"happy path", state{email: "john@info.co", uname: "john", pass: "Test123!", cpass: "Test123!", name: "John"}, want{}},
		{"missing email", state{email: "", uname: "john", pass: "Test123!", cpass: "Test123!", name: "John"}, want{"email must be a valid email address"}},
		{"invalid email", state{email: "blah", uname: "john", pass: "Test123!", cpass: "Test123!", name: "John"}, want{"email must be a valid email address"}},
		{"missing username", state{email: "john@info.co", uname: "", pass: "Test123!", cpass: "Test123!", name: "John"}, want{"username must be at least 4 characters long"}},
		{"missing name", state{email: "john@info.co", uname: "john", pass: "Test123!", cpass: "Test123!", name: ""}, want{"name must be at least 2 characters long"}},
		{"invalid pw length", state{email: "john@info.co", uname: "john", pass: "Tt1!", cpass: "Tt1!", name: "John"}, want{"password must be within 8 and 256 characters in length"}},
		{"invalid pw upper", state{email: "john@info.co", uname: "john", pass: "test123!", cpass: "test123!", name: "John"}, want{"password must have at least one upper case letter"}},
		{"invalid pw lower", state{email: "john@info.co", uname: "john", pass: "TEST123!", cpass: "TEST123!", name: "John"}, want{"password must have at least one lower case letter"}},
		{"invalid pw special", state{email: "john@info.co", uname: "john", pass: "Test1231", cpass: "Test1231", name: "John"}, want{"password must have at least one special character"}},
		{"invalid pw digit", state{email: "john@info.co", uname: "john", pass: "Testtes!", cpass: "Testtes!", name: "John"}, want{"password must have at least one digit"}},
		{"invalid confirm password", state{email: "john@info.co", uname: "john", pass: "Test123!", cpass: "T", name: "John"}, want{"confirm password does not match password"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &userRegisterRequest{
				Username:        c.state.uname,
				Password:        c.state.pass,
				ConfirmPassword: c.state.cpass,
				Email:           c.state.email,
				Name:            c.state.name,
			}

			err := validateUserRegister(req)
			tests.AssertError(t, err, c.want.err)
		})
	}
}
