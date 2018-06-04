// Code generated by mockery v1.0.0. DO NOT EDIT.
package diskops

import mock "github.com/stretchr/testify/mock"
import models "github.com/Stratoscale/disk-manager-exercise/models"

// MockDiskAPI is an autogenerated mock type for the DiskAPI type
type MockDiskAPI struct {
	mock.Mock
}

// DiskByID provides a mock function with given fields: id
func (_m *MockDiskAPI) DiskByID(id string) (*models.Disk, error) {
	ret := _m.Called(id)

	var r0 *models.Disk
	if rf, ok := ret.Get(0).(func(string) *models.Disk); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Disk)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDisks provides a mock function with given fields: hostName
func (_m *MockDiskAPI) ListDisks(hostName *string) (models.ListDisksOKBody, error) {
	ret := _m.Called(hostName)

	var r0 models.ListDisksOKBody
	if rf, ok := ret.Get(0).(func(*string) models.ListDisksOKBody); ok {
		r0 = rf(hostName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.ListDisksOKBody)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(hostName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}