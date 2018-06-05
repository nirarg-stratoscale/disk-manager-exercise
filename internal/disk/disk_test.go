package disk

import (
	"context"
	"testing"

	"github.com/Stratoscale/disk-manager-exercise/internal/diskops"
	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/disk-manager-exercise/restapi/operations/disk"
	"github.com/Stratoscale/go-template/golib/testutil"
	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/assert"
)

var log = testutil.Log()

func TestListDisks(t *testing.T) {
	t.Parallel()

	id := "1234"
	hostName := "Test1"
	out := models.ListDisksOKBody{&models.Disk{ID: &id}}

	tests := []struct {
		name   string
		params disk.ListDisksParams
		want   middleware.Responder
	}{
		{
			name:   "TestListDisks-1",
			params: disk.ListDisksParams{Hostname: &hostName},
			want:   disk.NewListDisksOK().WithPayload(out),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				db = testutil.OpenDB(t)

				diskApiMock = new(diskops.MockDiskAPI)

				p = New(Config{DB: db, Log: log, DiskAPI: diskApiMock})
			)

			diskApiMock.On("ListDisks", &hostName).Return(out, nil)

			testutil.SyncAutoMigrate(t, p.AutoMigrate)

			got := p.ListDisks(context.Background(), tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDiskById(t *testing.T) {
	t.Parallel()

	id := "1234"
	out := models.Disk{ID: &id}

	tests := []struct {
		name   string
		params disk.DiskByIDParams
		want   middleware.Responder
	}{
		{
			name:   "TestDiskById-1",
			params: disk.DiskByIDParams{DiskID: id},
			want:   disk.NewDiskByIDOK().WithPayload(&out),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				db = testutil.OpenDB(t)

				diskApiMock = new(diskops.MockDiskAPI)

				p = New(Config{DB: db, Log: log, DiskAPI: diskApiMock})
			)

			diskApiMock.On("DiskByID", tt.params.DiskID).Return(&out, nil)

			testutil.SyncAutoMigrate(t, p.AutoMigrate)

			got := p.DiskByID(context.Background(), tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}
