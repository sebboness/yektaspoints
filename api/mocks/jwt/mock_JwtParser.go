// Code generated by mockery v2.40.1. DO NOT EDIT.

package jwt

import mock "github.com/stretchr/testify/mock"

// MockJwtParser is an autogenerated mock type for the JwtParser type
type MockJwtParser struct {
	mock.Mock
}

type MockJwtParser_Expecter struct {
	mock *mock.Mock
}

func (_m *MockJwtParser) EXPECT() *MockJwtParser_Expecter {
	return &MockJwtParser_Expecter{mock: &_m.Mock}
}

// GetJwtClaims provides a mock function with given fields: tokenValue
func (_m *MockJwtParser) GetJwtClaims(tokenValue string) (map[string]interface{}, error) {
	ret := _m.Called(tokenValue)

	if len(ret) == 0 {
		panic("no return value specified for GetJwtClaims")
	}

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (map[string]interface{}, error)); ok {
		return rf(tokenValue)
	}
	if rf, ok := ret.Get(0).(func(string) map[string]interface{}); ok {
		r0 = rf(tokenValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockJwtParser_GetJwtClaims_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetJwtClaims'
type MockJwtParser_GetJwtClaims_Call struct {
	*mock.Call
}

// GetJwtClaims is a helper method to define mock.On call
//   - tokenValue string
func (_e *MockJwtParser_Expecter) GetJwtClaims(tokenValue interface{}) *MockJwtParser_GetJwtClaims_Call {
	return &MockJwtParser_GetJwtClaims_Call{Call: _e.mock.On("GetJwtClaims", tokenValue)}
}

func (_c *MockJwtParser_GetJwtClaims_Call) Run(run func(tokenValue string)) *MockJwtParser_GetJwtClaims_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockJwtParser_GetJwtClaims_Call) Return(_a0 map[string]interface{}, _a1 error) *MockJwtParser_GetJwtClaims_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockJwtParser_GetJwtClaims_Call) RunAndReturn(run func(string) (map[string]interface{}, error)) *MockJwtParser_GetJwtClaims_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockJwtParser creates a new instance of MockJwtParser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockJwtParser(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJwtParser {
	mock := &MockJwtParser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}