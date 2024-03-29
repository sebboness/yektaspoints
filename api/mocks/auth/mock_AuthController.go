// Code generated by mockery v2.40.1. DO NOT EDIT.

package auth

import (
	context "context"

	auth "github.com/sebboness/yektaspoints/util/auth"

	mock "github.com/stretchr/testify/mock"
)

// MockAuthController is an autogenerated mock type for the AuthController type
type MockAuthController struct {
	mock.Mock
}

type MockAuthController_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuthController) EXPECT() *MockAuthController_Expecter {
	return &MockAuthController_Expecter{mock: &_m.Mock}
}

// AssignUserToRole provides a mock function with given fields: ctx, username, role
func (_m *MockAuthController) AssignUserToRole(ctx context.Context, username string, role string) error {
	ret := _m.Called(ctx, username, role)

	if len(ret) == 0 {
		panic("no return value specified for AssignUserToRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, username, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthController_AssignUserToRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AssignUserToRole'
type MockAuthController_AssignUserToRole_Call struct {
	*mock.Call
}

// AssignUserToRole is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - role string
func (_e *MockAuthController_Expecter) AssignUserToRole(ctx interface{}, username interface{}, role interface{}) *MockAuthController_AssignUserToRole_Call {
	return &MockAuthController_AssignUserToRole_Call{Call: _e.mock.On("AssignUserToRole", ctx, username, role)}
}

func (_c *MockAuthController_AssignUserToRole_Call) Run(run func(ctx context.Context, username string, role string)) *MockAuthController_AssignUserToRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAuthController_AssignUserToRole_Call) Return(_a0 error) *MockAuthController_AssignUserToRole_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthController_AssignUserToRole_Call) RunAndReturn(run func(context.Context, string, string) error) *MockAuthController_AssignUserToRole_Call {
	_c.Call.Return(run)
	return _c
}

// Authenticate provides a mock function with given fields: ctx, username, password
func (_m *MockAuthController) Authenticate(ctx context.Context, username string, password string) (auth.AuthResult, error) {
	ret := _m.Called(ctx, username, password)

	if len(ret) == 0 {
		panic("no return value specified for Authenticate")
	}

	var r0 auth.AuthResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (auth.AuthResult, error)); ok {
		return rf(ctx, username, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) auth.AuthResult); ok {
		r0 = rf(ctx, username, password)
	} else {
		r0 = ret.Get(0).(auth.AuthResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthController_Authenticate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authenticate'
type MockAuthController_Authenticate_Call struct {
	*mock.Call
}

// Authenticate is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - password string
func (_e *MockAuthController_Expecter) Authenticate(ctx interface{}, username interface{}, password interface{}) *MockAuthController_Authenticate_Call {
	return &MockAuthController_Authenticate_Call{Call: _e.mock.On("Authenticate", ctx, username, password)}
}

func (_c *MockAuthController_Authenticate_Call) Run(run func(ctx context.Context, username string, password string)) *MockAuthController_Authenticate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAuthController_Authenticate_Call) Return(_a0 auth.AuthResult, _a1 error) *MockAuthController_Authenticate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthController_Authenticate_Call) RunAndReturn(run func(context.Context, string, string) (auth.AuthResult, error)) *MockAuthController_Authenticate_Call {
	_c.Call.Return(run)
	return _c
}

// ConfirmRegistration provides a mock function with given fields: ctx, username, code
func (_m *MockAuthController) ConfirmRegistration(ctx context.Context, username string, code string) error {
	ret := _m.Called(ctx, username, code)

	if len(ret) == 0 {
		panic("no return value specified for ConfirmRegistration")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, username, code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthController_ConfirmRegistration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfirmRegistration'
type MockAuthController_ConfirmRegistration_Call struct {
	*mock.Call
}

// ConfirmRegistration is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - code string
func (_e *MockAuthController_Expecter) ConfirmRegistration(ctx interface{}, username interface{}, code interface{}) *MockAuthController_ConfirmRegistration_Call {
	return &MockAuthController_ConfirmRegistration_Call{Call: _e.mock.On("ConfirmRegistration", ctx, username, code)}
}

func (_c *MockAuthController_ConfirmRegistration_Call) Run(run func(ctx context.Context, username string, code string)) *MockAuthController_ConfirmRegistration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAuthController_ConfirmRegistration_Call) Return(_a0 error) *MockAuthController_ConfirmRegistration_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthController_ConfirmRegistration_Call) RunAndReturn(run func(context.Context, string, string) error) *MockAuthController_ConfirmRegistration_Call {
	_c.Call.Return(run)
	return _c
}

// RefreshToken provides a mock function with given fields: ctx, username, token
func (_m *MockAuthController) RefreshToken(ctx context.Context, username string, token string) (auth.AuthResult, error) {
	ret := _m.Called(ctx, username, token)

	if len(ret) == 0 {
		panic("no return value specified for RefreshToken")
	}

	var r0 auth.AuthResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (auth.AuthResult, error)); ok {
		return rf(ctx, username, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) auth.AuthResult); ok {
		r0 = rf(ctx, username, token)
	} else {
		r0 = ret.Get(0).(auth.AuthResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthController_RefreshToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RefreshToken'
type MockAuthController_RefreshToken_Call struct {
	*mock.Call
}

// RefreshToken is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - token string
func (_e *MockAuthController_Expecter) RefreshToken(ctx interface{}, username interface{}, token interface{}) *MockAuthController_RefreshToken_Call {
	return &MockAuthController_RefreshToken_Call{Call: _e.mock.On("RefreshToken", ctx, username, token)}
}

func (_c *MockAuthController_RefreshToken_Call) Run(run func(ctx context.Context, username string, token string)) *MockAuthController_RefreshToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAuthController_RefreshToken_Call) Return(_a0 auth.AuthResult, _a1 error) *MockAuthController_RefreshToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthController_RefreshToken_Call) RunAndReturn(run func(context.Context, string, string) (auth.AuthResult, error)) *MockAuthController_RefreshToken_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: ctx, ur
func (_m *MockAuthController) Register(ctx context.Context, ur auth.UserRegisterRequest) (auth.UserRegisterResult, error) {
	ret := _m.Called(ctx, ur)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 auth.UserRegisterResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.UserRegisterRequest) (auth.UserRegisterResult, error)); ok {
		return rf(ctx, ur)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auth.UserRegisterRequest) auth.UserRegisterResult); ok {
		r0 = rf(ctx, ur)
	} else {
		r0 = ret.Get(0).(auth.UserRegisterResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, auth.UserRegisterRequest) error); ok {
		r1 = rf(ctx, ur)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthController_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type MockAuthController_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - ctx context.Context
//   - ur auth.UserRegisterRequest
func (_e *MockAuthController_Expecter) Register(ctx interface{}, ur interface{}) *MockAuthController_Register_Call {
	return &MockAuthController_Register_Call{Call: _e.mock.On("Register", ctx, ur)}
}

func (_c *MockAuthController_Register_Call) Run(run func(ctx context.Context, ur auth.UserRegisterRequest)) *MockAuthController_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(auth.UserRegisterRequest))
	})
	return _c
}

func (_c *MockAuthController_Register_Call) Return(_a0 auth.UserRegisterResult, _a1 error) *MockAuthController_Register_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthController_Register_Call) RunAndReturn(run func(context.Context, auth.UserRegisterRequest) (auth.UserRegisterResult, error)) *MockAuthController_Register_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePassword provides a mock function with given fields: ctx, session, username, password
func (_m *MockAuthController) UpdatePassword(ctx context.Context, session string, username string, password string) error {
	ret := _m.Called(ctx, session, username, password)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, session, username, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthController_UpdatePassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePassword'
type MockAuthController_UpdatePassword_Call struct {
	*mock.Call
}

// UpdatePassword is a helper method to define mock.On call
//   - ctx context.Context
//   - session string
//   - username string
//   - password string
func (_e *MockAuthController_Expecter) UpdatePassword(ctx interface{}, session interface{}, username interface{}, password interface{}) *MockAuthController_UpdatePassword_Call {
	return &MockAuthController_UpdatePassword_Call{Call: _e.mock.On("UpdatePassword", ctx, session, username, password)}
}

func (_c *MockAuthController_UpdatePassword_Call) Run(run func(ctx context.Context, session string, username string, password string)) *MockAuthController_UpdatePassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockAuthController_UpdatePassword_Call) Return(_a0 error) *MockAuthController_UpdatePassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthController_UpdatePassword_Call) RunAndReturn(run func(context.Context, string, string, string) error) *MockAuthController_UpdatePassword_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAuthController creates a new instance of MockAuthController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthController {
	mock := &MockAuthController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
