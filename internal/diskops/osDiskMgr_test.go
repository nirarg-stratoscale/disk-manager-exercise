package diskops

import (
	"testing"

	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/go-template/golib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/Stratoscale/disk-manager-exercise/internal/osops"
)

var log = testutil.Log()

func TestListDisks(t *testing.T) {
	t.Parallel()

	id := "1234"
	path := "/dev/sda"

	tests := []struct {
		name  string
		arg1  string
		want  models.ListDisksOKBody
		pyResult string
		hostResult string
		wantErr bool
	}{
		{
			name: "TestListDisks-1",
			arg1: "Test1",
			want: models.ListDisksOKBody{&models.Disk{ID: &id,
				Hostname: "host-test-1",
				MediaType: "SSD",
				Model: "SAMSUNG MZ7TN512",
				Path: &path,
				Serial: "S35NNY0HA05094",
				TotalCapacityMB: 512110}},
			pyResult: "[{\"path\": \"/dev/sda\", \"serial\": \"S35NNY0HA05094\", \"model\": \"SAMSUNG MZ7TN512\", \"totalCapacityMB\": 512110, \"mediaType\": \"SSD\"}]",
			hostResult: "host-test-1",
			wantErr: false,
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
			osOpsMock.On("ExecCommand", "python", "lib/pytools/storage.py").Return(tt.pyResult, nil)
			osOpsMock.On("Hostname").Return(tt.hostResult, nil)

			got, err := p.ListDisks(&tt.arg1)

			// assert results expectations
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got) // notice the go convention: want is first argument, got is the second argument
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
