package consul

import (
	"errors"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/mock"
)

// MockKV is mock implementation of consul KV service.
type MockKV struct {
	mock.Mock
}

// MockGet mocks get key
func (m *MockKV) MockGet(key string, value []byte) *mock.Call {
	return m.On("Get", key, (*api.QueryOptions)(nil)).
		Return(&api.KVPair{Key: key, Value: value}, nil, nil)
}

// MockGetMissing mocks get key which is missing
func (m *MockKV) MockGetMissing(key string) *mock.Call {
	return m.On("Get", key, (*api.QueryOptions)(nil)).Return(nil, nil, nil)
}

// MockGetFailure mocks get key failure
func (m *MockKV) MockGetFailure(key string) *mock.Call {
	return m.On("Get", key, (*api.QueryOptions)(nil)).Return(nil, nil, errors.New("failed"))
}

// MockPut mocks get key
func (m *MockKV) MockPut(key string, value string) *mock.Call {
	return m.On("Put", &api.KVPair{Key: key, Value: []byte(value)}).
		Return((*api.WriteMeta)(nil), nil)
}

// MockPutFailure mocks get key
func (m *MockKV) MockPutFailure(key string, value string) *mock.Call {
	return m.On("Put", &api.KVPair{Key: key, Value: []byte(value)}).
		Return((*api.WriteMeta)(nil), errors.New("failed"))
}

// Get is used to lookup a single key. The returned pointer
// to the KVPair will be nil if the key does not exist.
func (m *MockKV) Get(key string, q *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error) {
	args := m.Called(key, q)
	var pair *api.KVPair
	if args.Get(0) != nil {
		pair = args.Get(0).(*api.KVPair)
	}
	return pair, nil, args.Error(2)
}

// Put is used to write a new value. Only the
// Key, Flags and Value is respected.
func (m *MockKV) Put(p *api.KVPair, q *api.WriteOptions) (*api.WriteMeta, error) {
	args := m.Called(p)
	return nil, args.Error(1)
}

// Keys implements the KV.Keys method.
func (m *MockKV) Keys(p1, p2 string, q *api.QueryOptions) ([]string, *api.QueryMeta, error) {
	args := m.Called(p1, p2, q)
	return args.Get(0).([]string), nil, args.Error(2)
}

// List implements the KV.List method.
func (m *MockKV) List(p string, q *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error) {
	args := m.Called(p, q)
	return args.Get(0).(api.KVPairs), nil, args.Error(2)
}

// Delete implements the KV.Delete method.
func (m *MockKV) Delete(key string, w *api.WriteOptions) (*api.WriteMeta, error) {
	args := m.Called(key)
	return nil, args.Error(1)
}

// CAS Check-And-Set operation
func (m *MockKV) CAS(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error) {
	args := m.Called(p, q)
	return args.Get(0).(bool), nil, args.Error(2)
}
