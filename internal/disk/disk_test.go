package disk

import (
	"context"
	"errors"
	"testing"

	"fmt"

	"github.com/Stratoscale/disk-manager-exercise/internal/diskops"
	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/disk-manager-exercise/restapi/operations/disk"
	"github.com/Stratoscale/go-template/golib/testutil"
	"github.com/Stratoscale/golib/httputil"
	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/assert"
)

var log = testutil.Log()

func TestListDisks(t *testing.T) {
	t.Parallel()

	id := "1234"
	d := models.Disk{ID: &id}
	hostName := "Test1"
	out := []*models.Disk{&d}

	tests := []struct {
		name    string
		params  disk.ListDisksParams
		want    middleware.Responder
		wantErr error
	}{
		{
			name:   "TestListDisks-ok",
			params: disk.ListDisksParams{Hostname: &hostName},
			want:   disk.NewListDisksOK().WithPayload(out),
		},
		{
			name:    "TestListDisks-httpErr",
			params:  disk.ListDisksParams{Hostname: &hostName},
			wantErr: httputil.NewErrBadRequest("ListDisk bad request error"),
			want:    httputil.NewErrBadRequest("ListDisk bad request error"),
		},
		{
			name:    "TestListDisks-err",
			params:  disk.ListDisksParams{Hostname: &hostName},
			wantErr: errors.New("ListDisk not http error"),
			want: disk.NewListDisksInternalServerError().WithPayload(models.Error500(
				fmt.Sprintf("ListDisks with hostname %s failed %s", hostName,
					errors.New("ListDisk not http error")))),
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

			diskApiMock.On("ListDisks", &hostName).Return(out, tt.wantErr)

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
		name    string
		params  disk.DiskByIDParams
		want    middleware.Responder
		wantErr error
	}{
		{
			name:   "TestDiskById-ok",
			params: disk.DiskByIDParams{DiskID: id},
			want:   disk.NewDiskByIDOK().WithPayload(&out),
		},
		{
			name:    "TestDiskById-httpErr",
			params:  disk.DiskByIDParams{DiskID: id},
			wantErr: httputil.NewErrNotFound("DiskById not found error"),
			want:    httputil.NewErrNotFound("DiskById not found error"),
		},
		{
			name:    "TestDiskById-err",
			params:  disk.DiskByIDParams{DiskID: id},
			wantErr: errors.New("DiskById not http error"),
			want: disk.NewDiskByIDInternalServerError().WithPayload(models.Error500(
				fmt.Sprintf("DiskByID with id %s failed %s", id, errors.New("DiskById not http error")))),
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

			diskApiMock.On("DiskByID", tt.params.DiskID).Return(&out, tt.wantErr)

			testutil.SyncAutoMigrate(t, p.AutoMigrate)

			got := p.DiskByID(context.Background(), tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}
