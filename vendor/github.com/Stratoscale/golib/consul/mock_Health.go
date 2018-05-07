// Code generated by mockery v1.0.0
package consul

import api "github.com/hashicorp/consul/api"
import mock "github.com/stretchr/testify/mock"

// MockHealth is an autogenerated mock type for the Health type
type MockHealth struct {
	mock.Mock
}

// Service provides a mock function with given fields: service, tag, passingOnly, q
func (_m *MockHealth) Service(service string, tag string, passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {
	ret := _m.Called(service, tag, passingOnly, q)

	var r0 []*api.ServiceEntry
	if rf, ok := ret.Get(0).(func(string, string, bool, *api.QueryOptions) []*api.ServiceEntry); ok {
		r0 = rf(service, tag, passingOnly, q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*api.ServiceEntry)
		}
	}

	var r1 *api.QueryMeta
	if rf, ok := ret.Get(1).(func(string, string, bool, *api.QueryOptions) *api.QueryMeta); ok {
		r1 = rf(service, tag, passingOnly, q)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*api.QueryMeta)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, bool, *api.QueryOptions) error); ok {
		r2 = rf(service, tag, passingOnly, q)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}