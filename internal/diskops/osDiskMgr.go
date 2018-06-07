package diskops

import (
	"encoding/json"

	"github.com/Stratoscale/disk-manager-exercise/internal/osops"
	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/golib/httputil"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const BaseUuid = "a99aac7c-688b-11e8-b6ea-c85b7692659d"

type pyStorageResponse struct {
	Path            string `json:"path"`
	Serial          string `json:"serial"`
	MediaType       string `json:"mediaType"`
	Model           string `json:"model"`
	TotalCapacityMB int64  `json:"totalCapacityMB"`
}

type ListResponse []*pyStorageResponse

type Config struct {
	Log   logrus.FieldLogger
	OsOps osops.OSOperations
}

func NewOsDiskMgr(c Config) *OsDiskMgr {
	return &OsDiskMgr{
		Config: c,
	}
}

type OsDiskMgr struct {
	Config
}

func (o *OsDiskMgr) ListDisks(hostName *string) ([]*models.Disk, error) {
	out, err1 := o.OsOps.ExecCommand("python", "lib/pytools/storage.py")
	if err1 != nil {
		return nil, httputil.NewErrInternalServer("ListDisks failed to get the disks info with error %s", err1)
	}

	result := []*models.Disk{}

	lst := ListResponse{}
	err2 := json.Unmarshal([]byte(out), &lst)
	if err2 != nil {
		return nil, httputil.NewErrInternalServer("ListDisks failed to unmarshal os json response with error %s", err2)
	}
	hostname, err3 := o.OsOps.ExecCommand("cat", "/etc/hostname")
	if err3 != nil {
		return nil, httputil.NewErrInternalServer("ListDisks failed to get the hostname with error %s", err3)
	}
	for _, val := range lst {
		baseUuid, err4 := uuid.FromString(BaseUuid)
		if err4 != nil {
			return nil, httputil.NewErrInternalServer("ListDisks failed to generate uuid from serial %s with error %s", val.Serial, err4)
		}
		uuidFromSerial := uuid.NewV3(baseUuid, val.Serial)
		id := uuidFromSerial.String()
		disk := models.Disk{ID: &id,
			Hostname:        hostname,
			MediaType:       val.MediaType,
			Model:           val.Model,
			Path:            &val.Path,
			Serial:          val.Serial,
			TotalCapacityMB: val.TotalCapacityMB}
		result = append(result, &disk)
	}

	return result, nil
}

func (o *OsDiskMgr) DiskByID(id string) (*models.Disk, error) {
	lst, err := o.ListDisks(nil)
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, httputil.NewErrInternalServer("DiskByID with id %s failed because ListDisks returns empty list", id)
	}
	for _, disk := range lst {
		if disk != nil && id == *disk.ID {
			return disk, nil
		}
	}
	return nil, httputil.NewErrBadRequest("Invalid ID %s", id)
}
