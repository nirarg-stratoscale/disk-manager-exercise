// Code generated by mockery v1.0.0. DO NOT EDIT.
package disk

import context "context"
import mock "github.com/stretchr/testify/mock"

// MockAPI is an autogenerated mock type for the API type
type MockAPI struct {
	mock.Mock
}

// DiskByID provides a mock function with given fields: ctx, params
func (_m *MockAPI) DiskByID(ctx context.Context, params *DiskByIDParams) (*DiskByIDOK, error) {
	ret := _m.Called(ctx, params)

	var r0 *DiskByIDOK
	if rf, ok := ret.Get(0).(func(context.Context, *DiskByIDParams) *DiskByIDOK); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*DiskByIDOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *DiskByIDParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDisks provides a mock function with given fields: ctx, params
func (_m *MockAPI) ListDisks(ctx context.Context, params *ListDisksParams) (*ListDisksOK, error) {
	ret := _m.Called(ctx, params)

	var r0 *ListDisksOK
	if rf, ok := ret.Get(0).(func(context.Context, *ListDisksParams) *ListDisksOK); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ListDisksOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ListDisksParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
