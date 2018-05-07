package consulutil

import (
	"github.com/pkg/errors"

	"github.com/Stratoscale/golib/consul"
)

// Lock lock a key in consul and returns an unlock function that should be deferred.
// Since it is not using any health-checks, it holds the lock until the unlock is called or
// until the process dies (there are renewals every 7 seconds, and the lock is released 15 seconds
// after the last renewal.
// the returned unlock function fails only if the lock is not held, which in this case won't happened
// since no health-check was defiend.
func Lock(locker consul.Locker, key string) (unlock func() error, err error) {
	lock, err := locker.LockKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "creating lock")

	}
	if _, err = lock.Lock(nil); err != nil {
		return nil, errors.Wrap(err, "locking")
	}

	return lock.Unlock, nil
}
