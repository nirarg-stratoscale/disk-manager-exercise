package dbutil

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

// Config is configuration for database from environment variables
type Config struct {
	// ConnString is for direct use of Open
	ConnString string `envconfig:"DB_CONN_STRING"`

	//RootConnString and Name are to be used by Credentialer
	RootConnString string `envconfig:"DB_ROOT_CONN_STRING"`
	Name           string `envconfig:"DB_NAME"`

	Debug bool `envconfig:"DEBUG"`
}

// Open opens a new configured DB connection
func Open(cfg Config, logger logrus.FieldLogger) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", cfg.ConnString)
	if err != nil {
		return nil, err
	}
	db.LogMode(cfg.Debug)
	db.SetLogger(gorm.Logger{LogWriter: logger})
	return db, nil
}
