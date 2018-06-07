package disk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/disk-manager-exercise/restapi/operations/disk"
	"github.com/Stratoscale/golib/httputil"
	"github.com/Stratoscale/go-template/golib/testutil"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var log = testutil.Log()





func TestListDisks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		params  disk.ListDisksParams
		want    middleware.Responder
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

				p  = New(Config{DB: db, Log: log})
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
		name    string
		params  disk.DiskByIdParams
		want    middleware.Responder
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

				p  = New(Config{DB: db, Log: log})
			)


			testutil.SyncAutoMigrate(t, p.AutoMigrate)


			got := p.DiskById(context.Background(), tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}
