package testutil

import (
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func OpenDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open("sqlite3", ":memory:")
	require.Nil(t, err)

	if testing.Verbose() {
		db.LogMode(true)
		db.SetLogger(&testLogWriter{T: t})
	}

	return db
}

// SyncAutoMigrate is a function for syncing automigrate database.
// database automigrate is not safe for concurrent use.
// this function synchronize the calls
func SyncAutoMigrate(t *testing.T, autoMigrate func() error) {
	t.Helper()
	autoMigrateLock.Lock()
	defer autoMigrateLock.Unlock()

	require.Nil(t, autoMigrate())
}

var autoMigrateLock sync.Mutex

// testLogWriter is an adapter for gorm that implements the logger interface for db.SetLogger
type testLogWriter struct {
	*testing.T
}

func (t *testLogWriter) Print(args ...interface{}) {
	t.Log(args...)
}
