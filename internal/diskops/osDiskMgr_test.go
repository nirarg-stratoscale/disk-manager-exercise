package diskops

import (
	"testing"

	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/go-template/golib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/Stratoscale/disk-manager-exercise/internal/osops"
	"fmt"
	"github.com/Stratoscale/golib/httputil"
)

var log = testutil.Log()

func TestListDisks(t *testing.T) {
	t.Parallel()

	id1 := "3ee19fad-c471-3adb-bce1-7c6c025cea7a"
	id2 := "619b5756-d6bd-3208-88f0-e97326f5fd5f"
	path1 := "/dev/sda1"
	path2 := "/dev/sda2"

	tests := []struct {
		name  string
		arg1  string
		want  models.ListDisksOKBody
		pyResult string
		pyErr error
		hostResult string
		hostErr error
		wantErr error
	}{
		{
			name: "TestListDisks-ok",
			arg1: "Test1",
			want: models.ListDisksOKBody{&models.Disk{ID: &id1,
											Hostname: "host-test-1",
											MediaType: "SSD",
											Model: "SAMSUNG1 MZ7TN512",
											Path: &path1,
											Serial: "S35NNY0HA05094111",
											TotalCapacityMB: 5121101},
										&models.Disk{ID: &id2,
											Hostname: "host-test-1",
											MediaType: "HDD",
											Model: "SAMSUNG2 MZ7TN512",
											Path: &path2,
											Serial: "S35NNY0HA05094222",
											TotalCapacityMB: 5121102}},
			pyResult: "[{\"path\": \"/dev/sda1\", \"serial\": \"S35NNY0HA05094111\", \"model\": \"SAMSUNG1 MZ7TN512\", \"totalCapacityMB\": 5121101, \"mediaType\": \"SSD\"}," +
				"{\"path\": \"/dev/sda2\", \"serial\": \"S35NNY0HA05094222\", \"model\": \"SAMSUNG2 MZ7TN512\", \"totalCapacityMB\": 5121102, \"mediaType\": \"HDD\"}]",
			hostResult: "host-test-1",
			wantErr: nil,
		},
		{
			name: "TestListDisks-unmarshal error",
			arg1: "Test2",
			pyResult: "",
			hostResult: "host-test-1",
			wantErr: httputil.NewErrInternalServer("ListDisks failed to unmarshal os json response with error unexpected end of JSON input"),
		},
		{
			name: "TestListDisks-python error",
			arg1: "Test3",
			pyErr: fmt.Errorf("test error"),
			hostResult: "host-test-1",
			wantErr: httputil.NewErrInternalServer("ListDisks failed to get the disks info with error test error"),
		},
		{
			name: "TestListDisks-hostname error",
			arg1: "Test4",
			pyResult: "[{\"path\": \"/dev/sda1\", \"serial\": \"S35NNY0HA05094111\", \"model\": \"SAMSUNG1 MZ7TN512\", \"totalCapacityMB\": 5121101, \"mediaType\": \"SSD\"}]",
			hostErr: fmt.Errorf("hostname test error"),
			wantErr: httputil.NewErrInternalServer("ListDisks failed to get the hostname with error hostname test error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				osOpsMock = new(osops.MockOSOperations)

				p  = NewOsDiskMgr(Config{Log: log, OsOps: osOpsMock})
			)
			osOpsMock.On("ExecCommand", "python", "lib/pytools/storage.py").Return(tt.pyResult, tt.pyErr)
			osOpsMock.On("Hostname").Return(tt.hostResult, tt.hostErr)

			got, err := p.ListDisks(&tt.arg1)

			// assert results expectations
			if tt.wantErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got) // notice the go convention: want is first argument, got is the second argument
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDiskByID(t *testing.T) {
	t.Parallel()

	id := "3ee19fad-c471-3adb-bce1-7c6c025cea7a"
	path1 := "/dev/sda1"

	tests := []struct {
		name  string
		arg1  string
		want  *models.Disk
		pyResult string
		pyErr error
		hostResult string
		hostErr error
		wantErr error
	}{
		{
			name: "TestDiskByID-ok",
			arg1: id,
			want: &models.Disk{ID: &id,
				Hostname: "host-test-1",
				MediaType: "SSD",
				Model: "SAMSUNG1 MZ7TN512",
				Path: &path1,
				Serial: "S35NNY0HA05094111",
				TotalCapacityMB: 5121101},
			pyResult: "[{\"path\": \"/dev/sda1\", \"serial\": \"S35NNY0HA05094111\", \"model\": \"SAMSUNG1 MZ7TN512\", \"totalCapacityMB\": 5121101, \"mediaType\": \"SSD\"}]",
			hostResult: "host-test-1",
		},
		{
			name: "TestDiskByID-Bad request",
			arg1: "notId",
			pyResult: "[{\"path\": \"/dev/sda1\", \"serial\": \"S35NNY0HA05094111\", \"model\": \"SAMSUNG1 MZ7TN512\", \"totalCapacityMB\": 5121101, \"mediaType\": \"SSD\"}]",
			hostResult: "host-test-1",
			wantErr: httputil.NewErrBadRequest("Invalid ID notId"),
		},
		{
			name: "TestDiskByID-empty disks",
			arg1: id,
			pyResult: "[]",
			hostResult: "host-test-1",
			wantErr: httputil.NewErrInternalServer("DiskByID with id %s failed because ListDisks returns empty list", id),
		},
		{
			name: "TestDiskByID-unmarshal error",
			arg1: id,
			pyResult: "",
			hostResult: "host-test-1",
			wantErr: httputil.NewErrInternalServer("ListDisks failed to unmarshal os json response with error unexpected end of JSON input"),
		},
		{
			name: "TestDiskByID-python error",
			arg1: id,
			pyErr: fmt.Errorf("test error"),
			hostResult: "host-test-1",
			wantErr: httputil.NewErrInternalServer("ListDisks failed to get the disks info with error test error"),
		},
		{
			name: "TestDiskByID-hostname error",
			arg1: id,
			pyResult: "[{\"path\": \"/dev/sda1\", \"serial\": \"S35NNY0HA05094111\", \"model\": \"SAMSUNG1 MZ7TN512\", \"totalCapacityMB\": 5121101, \"mediaType\": \"SSD\"}]",
			hostErr: fmt.Errorf("hostname test error"),
			wantErr: httputil.NewErrInternalServer("ListDisks failed to get the hostname with error hostname test error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				osOpsMock = new(osops.MockOSOperations)

				p  = NewOsDiskMgr(Config{Log: log, OsOps: osOpsMock})
			)
			osOpsMock.On("ExecCommand", "python", "lib/pytools/storage.py").Return(tt.pyResult, tt.pyErr)
			osOpsMock.On("Hostname").Return(tt.hostResult, tt.hostErr)

			got, err := p.DiskByID(tt.arg1)

			// assert results expectations
			if tt.wantErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got) // notice the go convention: want is first argument, got is the second argument
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
