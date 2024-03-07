// Code generated by mockery v2.40.1. DO NOT EDIT.

package storage

import (
	context "context"

	models "github.com/sebboness/yektaspoints/models"
	mock "github.com/stretchr/testify/mock"
)

// MockIUserStorage is an autogenerated mock type for the IUserStorage type
type MockIUserStorage struct {
	mock.Mock
}

type MockIUserStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIUserStorage) EXPECT() *MockIUserStorage_Expecter {
	return &MockIUserStorage_Expecter{mock: &_m.Mock}
}

// GetUserByID provides a mock function with given fields: ctx, userId
func (_m *MockIUserStorage) GetUserByID(ctx context.Context, userId string) (models.User, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByID")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.User, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.User); ok {
		r0 = rf(ctx, userId)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIUserStorage_GetUserByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByID'
type MockIUserStorage_GetUserByID_Call struct {
	*mock.Call
}

// GetUserByID is a helper method to define mock.On call
//   - ctx context.Context
//   - userId string
func (_e *MockIUserStorage_Expecter) GetUserByID(ctx interface{}, userId interface{}) *MockIUserStorage_GetUserByID_Call {
	return &MockIUserStorage_GetUserByID_Call{Call: _e.mock.On("GetUserByID", ctx, userId)}
}

func (_c *MockIUserStorage_GetUserByID_Call) Run(run func(ctx context.Context, userId string)) *MockIUserStorage_GetUserByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIUserStorage_GetUserByID_Call) Return(_a0 models.User, _a1 error) *MockIUserStorage_GetUserByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIUserStorage_GetUserByID_Call) RunAndReturn(run func(context.Context, string) (models.User, error)) *MockIUserStorage_GetUserByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIUserStorage creates a new instance of MockIUserStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIUserStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIUserStorage {
	mock := &MockIUserStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
