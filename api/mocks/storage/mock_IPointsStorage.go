// Code generated by mockery v2.40.1. DO NOT EDIT.

package storage

import (
	context "context"

	models "github.com/sebboness/yektaspoints/models"
	mock "github.com/stretchr/testify/mock"
)

// MockIPointsStorage is an autogenerated mock type for the IPointsStorage type
type MockIPointsStorage struct {
	mock.Mock
}

type MockIPointsStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIPointsStorage) EXPECT() *MockIPointsStorage_Expecter {
	return &MockIPointsStorage_Expecter{mock: &_m.Mock}
}

// GetPointByID provides a mock function with given fields: ctx, userId, id
func (_m *MockIPointsStorage) GetPointByID(ctx context.Context, userId string, id string) (models.Point, error) {
	ret := _m.Called(ctx, userId, id)

	if len(ret) == 0 {
		panic("no return value specified for GetPointByID")
	}

	var r0 models.Point
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (models.Point, error)); ok {
		return rf(ctx, userId, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) models.Point); ok {
		r0 = rf(ctx, userId, id)
	} else {
		r0 = ret.Get(0).(models.Point)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userId, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIPointsStorage_GetPointByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPointByID'
type MockIPointsStorage_GetPointByID_Call struct {
	*mock.Call
}

// GetPointByID is a helper method to define mock.On call
//   - ctx context.Context
//   - userId string
//   - id string
func (_e *MockIPointsStorage_Expecter) GetPointByID(ctx interface{}, userId interface{}, id interface{}) *MockIPointsStorage_GetPointByID_Call {
	return &MockIPointsStorage_GetPointByID_Call{Call: _e.mock.On("GetPointByID", ctx, userId, id)}
}

func (_c *MockIPointsStorage_GetPointByID_Call) Run(run func(ctx context.Context, userId string, id string)) *MockIPointsStorage_GetPointByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockIPointsStorage_GetPointByID_Call) Return(_a0 models.Point, _a1 error) *MockIPointsStorage_GetPointByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIPointsStorage_GetPointByID_Call) RunAndReturn(run func(context.Context, string, string) (models.Point, error)) *MockIPointsStorage_GetPointByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetPointsByUserID provides a mock function with given fields: ctx, userId, filters
func (_m *MockIPointsStorage) GetPointsByUserID(ctx context.Context, userId string, filters models.QueryPointsFilter) ([]models.Point, error) {
	ret := _m.Called(ctx, userId, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetPointsByUserID")
	}

	var r0 []models.Point
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.QueryPointsFilter) ([]models.Point, error)); ok {
		return rf(ctx, userId, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, models.QueryPointsFilter) []models.Point); ok {
		r0 = rf(ctx, userId, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Point)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, models.QueryPointsFilter) error); ok {
		r1 = rf(ctx, userId, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIPointsStorage_GetPointsByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPointsByUserID'
type MockIPointsStorage_GetPointsByUserID_Call struct {
	*mock.Call
}

// GetPointsByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userId string
//   - filters models.QueryPointsFilter
func (_e *MockIPointsStorage_Expecter) GetPointsByUserID(ctx interface{}, userId interface{}, filters interface{}) *MockIPointsStorage_GetPointsByUserID_Call {
	return &MockIPointsStorage_GetPointsByUserID_Call{Call: _e.mock.On("GetPointsByUserID", ctx, userId, filters)}
}

func (_c *MockIPointsStorage_GetPointsByUserID_Call) Run(run func(ctx context.Context, userId string, filters models.QueryPointsFilter)) *MockIPointsStorage_GetPointsByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(models.QueryPointsFilter))
	})
	return _c
}

func (_c *MockIPointsStorage_GetPointsByUserID_Call) Return(_a0 []models.Point, _a1 error) *MockIPointsStorage_GetPointsByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIPointsStorage_GetPointsByUserID_Call) RunAndReturn(run func(context.Context, string, models.QueryPointsFilter) ([]models.Point, error)) *MockIPointsStorage_GetPointsByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// SavePoint provides a mock function with given fields: ctx, point
func (_m *MockIPointsStorage) SavePoint(ctx context.Context, point models.Point) error {
	ret := _m.Called(ctx, point)

	if len(ret) == 0 {
		panic("no return value specified for SavePoint")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Point) error); ok {
		r0 = rf(ctx, point)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIPointsStorage_SavePoint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SavePoint'
type MockIPointsStorage_SavePoint_Call struct {
	*mock.Call
}

// SavePoint is a helper method to define mock.On call
//   - ctx context.Context
//   - point models.Point
func (_e *MockIPointsStorage_Expecter) SavePoint(ctx interface{}, point interface{}) *MockIPointsStorage_SavePoint_Call {
	return &MockIPointsStorage_SavePoint_Call{Call: _e.mock.On("SavePoint", ctx, point)}
}

func (_c *MockIPointsStorage_SavePoint_Call) Run(run func(ctx context.Context, point models.Point)) *MockIPointsStorage_SavePoint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Point))
	})
	return _c
}

func (_c *MockIPointsStorage_SavePoint_Call) Return(_a0 error) *MockIPointsStorage_SavePoint_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIPointsStorage_SavePoint_Call) RunAndReturn(run func(context.Context, models.Point) error) *MockIPointsStorage_SavePoint_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIPointsStorage creates a new instance of MockIPointsStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIPointsStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIPointsStorage {
	mock := &MockIPointsStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
