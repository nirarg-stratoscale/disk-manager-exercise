package diskops

import "github.com/Stratoscale/disk-manager-exercise/models"

//go:generate mockery -name DiskAPI -inpkg

// API is the disk managers (os ans sql) common api
type DiskAPI interface {
	ListDisks(hostName *string) (models.ListDisksOKBody, error)
	DiskByID(id string) (*models.Disk, error)
}
