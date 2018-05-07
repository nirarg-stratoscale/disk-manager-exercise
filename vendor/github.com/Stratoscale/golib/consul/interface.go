package consul

import "github.com/hashicorp/consul/api"

// KV defines a simple consul KV interface.
type KV interface {
	Get(string, *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)

	Put(*api.KVPair, *api.WriteOptions) (*api.WriteMeta, error)

	Keys(string, string, *api.QueryOptions) ([]string, *api.QueryMeta, error)

	List(string, *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error)

	Delete(key string, w *api.WriteOptions) (*api.WriteMeta, error)

	CAS(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error)
}

//go:generate mockery -name Health -inpkg

// Health defines a simple consul Health interface
type Health interface {
	Service(service, tag string, passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error)
}

//go:generate mockery -name Agent -inpkg

// Agent defines a simple consul Agent interface
type Agent interface {
	// Members returns the known gossip members. The WAN
	Members(wan bool) ([]*api.AgentMember, error)
}

// NewLocker wraps consul client to implement the Locker interface.
// This is needed because consul client returns an api.Lock object which
// can't be mocked in tests
func NewLocker(c *api.Client) Locker {
	return &locker{Client: c}
}

type locker struct {
	*api.Client
}

func (l *locker) LockKey(key string) (Lock, error) {
	return l.Client.LockKey(key)
}

func (l *locker) LockOpts(opts *api.LockOptions) (Lock, error) {
	return l.Client.LockOpts(opts)
}

//go:generate mockery -name Locker -inpkg

// Locker is the interface of generating a consul lock
type Locker interface {
	// LockKey returns a handle to a lock struct which can be used
	// to acquire and release the mutex. The key used must have
	// write permissions.
	LockKey(key string) (Lock, error)
	LockOpts(opts *api.LockOptions) (Lock, error)
}

//go:generate mockery -name Lock -inpkg

// Lock is the interface of a consul lock
type Lock interface {
	Lock(stopCh <-chan struct{}) (<-chan struct{}, error)
	Unlock() error
	Destroy() error
}

var (
	_ KV = new(api.Client).KV()
	_    = NewLocker(new(api.Client))
)
