package disk

import (
    "context"

	"github.com/Stratoscale/disk-manager-exercise/restapi/operations/disk"
	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/Stratoscale/golib/httputil"
	"github.com/go-openapi/runtime/middleware"
	"github.com/Stratoscale/disk-manager-exercise/internal/diskops"
	"fmt"
)

type Config struct {
	DB  *gorm.DB
	Log logrus.FieldLogger
	DiskAPI diskops.DiskAPI
}

func New(c Config) *manager {
	return &manager{
		Config:  c,
	}
}

type manager struct {
	Config
}

func (m *manager) AutoMigrate() error {
	return m.DB.AutoMigrate().Error
}

func (m *manager) ListDisks(ctx context.Context, params disk.ListDisksParams) middleware.Responder {
	out, err := m.DiskAPI.ListDisks(params.Hostname)
	if err != nil {
		httpErr, ok := err.(httputil.Error)
		if ok {
			return httpErr
		}
		return disk.NewListDisksInternalServerError().WithPayload(models.Error500(
			fmt.Sprintf("ListDisks with hostname %s failed %s", *params.Hostname, err)))
	}
	return disk.NewListDisksOK().WithPayload(out)

}

func (m *manager) DiskByID(ctx context.Context, params disk.DiskByIDParams) middleware.Responder {
	out, err := m.DiskAPI.DiskByID(params.DiskID)
	if err != nil {
		httpErr, ok := err.(httputil.Error)
		if ok {
			return httpErr
		}
		return disk.NewDiskByIDInternalServerError().WithPayload(models.Error500(
			fmt.Sprintf("DiskByID with id %s failed %s", params.DiskID, err)))
	}
	return disk.NewDiskByIDOK().WithPayload(out)

}
