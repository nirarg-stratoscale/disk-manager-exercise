package disk

import (
    "context"
    "net/http"
    
	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/disk-manager-exercise/restapi/operations/disk"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/Stratoscale/golib/httputil"
	"github.com/go-openapi/runtime/middleware"
)

type Config struct {
DB  *gorm.DB
Log logrus.FieldLogger
}

func New(c Config) *manager {
	return &manager{
		Config: c,
	}
}

type manager struct {
	Config
}


func (m *manager) AutoMigrate() error {
	return m.DB.AutoMigrate(
	).Error
}


func (m *manager) ListDisks(ctx context.Context, params disk.ListDisksParams) middleware.Responder {
    return httputil.NewError(http.StatusNotImplemented, "ListDisks not implemented yet")
}


func (m *manager) DiskById(ctx context.Context, params disk.DiskByIdParams) middleware.Responder {
    return httputil.NewError(http.StatusNotImplemented, "DiskById not implemented yet")
}
