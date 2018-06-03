package diskops

import (
	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/sirupsen/logrus"
	"github.com/Stratoscale/golib/httputil"
	"github.com/Stratoscale/disk-manager-exercise/internal/osops"
	"encoding/json"
)

type Config struct {
	Log logrus.FieldLogger
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

func (o *OsDiskMgr) ListDisks(hostName *string) (models.ListDisksOKBody , error) {
	out, err1 := o.OsOps.ExecCommand("python", "lib/pytools/storage.py")
	if err1 != nil {
		return nil, err1
	}

	result := models.ListDisksOKBody{}

	lst := models.ListDisksOKBody{}
	err2 := json.Unmarshal([]byte(out), &lst)
	if err2 != nil {
		return nil, err2
	}
	hostname, err3 := o.OsOps.Hostname()
	if err3 != nil {
		hostname = "UNKNOWN"
	}
	for _, val := range lst {
		id := "1234" //TODO generate UUID with seed
		disk := models.Disk{ID: &id,
			Hostname: hostname,
			MediaType: val.MediaType,
			Model: val.Model,
			Path: val.Path,
			Serial: val.Serial,
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
	return nil, httputil.NewErrBadRequest("Invalid ID %d", id)
}