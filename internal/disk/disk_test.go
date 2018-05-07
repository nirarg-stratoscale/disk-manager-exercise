package disk

import (
	"context"
	"net/http"
	"testing"

	"github.com/Stratoscale/disk-manager-exercise/restapi/operations/disk"
	"github.com/Stratoscale/go-template/golib/testutil"
	"github.com/Stratoscale/golib/httputil"
	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/assert"
)

var log = testutil.Log()

func TestListDisks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params disk.ListDisksParams
		want   middleware.Responder
	}{
		{
			name: "example",
			want: httputil.NewError(http.StatusNotImplemented, "ListDisks not implemented yet"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				db = testutil.OpenDB(t)

				p = New(Config{DB: db, Log: log})
			)

			testutil.SyncAutoMigrate(t, p.AutoMigrate)

			got := p.ListDisks(context.Background(), tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDiskById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params disk.DiskByIDParams
		want   middleware.Responder
	}{
		{
			name: "example",
			want: httputil.NewError(http.StatusNotImplemented, "DiskById not implemented yet"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				db = testutil.OpenDB(t)

				p = New(Config{DB: db, Log: log})
			)

			testutil.SyncAutoMigrate(t, p.AutoMigrate)

			got := p.DiskByID(context.Background(), tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}
