package dbutil

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/Stratoscale/golib/consul"
	"github.com/Stratoscale/golib/consul/consulutil"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	kvSuffix     = "/conn-string"
	kvLockSuffix = "/conn-string-lock"

	passwordLen = 32
)

func init() {
	rand.Seed(time.Now().Unix())
}

var reConnString = regexp.MustCompile(`([^:]+):([^@]*)@tcp\(([^:]+:\d+)\)/(.*)?`)

type Credentialer struct {
	KV     consul.KV
	Locker consul.Locker
	Log    logrus.FieldLogger

	kvKey string
}

// ConnectionString sets the ConnString in the Config
// It first searches in consul if the connection string is already stored there
// If it is there it returns it.
// If it does not exists in consul, it creates database and user in mysql,
// stores the connection string in consul and returns it.
func (c *Credentialer) ConnectionString(cfg *Config) error {
	c.kvKey = cfg.Name + kvSuffix
	kvLockKey := cfg.Name + kvLockSuffix

	unlock, err := consulutil.Lock(c.Locker, kvLockKey)
	if err != nil {
		return errors.Wrap(err, "lock consul")
	}
	defer unlock()

	connString, err := c.kvRead()
	if err != nil {
		return errors.Wrap(err, "reading from consul")
	}
	if connString != "" {
		if c.validate(connString) {
			cfg.ConnString = connString
			return nil
		}
		// connection string is invalid, recreate it
		c.kvReset()
	}

	connString, err = c.createDatabase(*cfg)
	if err != nil {
		return errors.Wrap(err, "creating service database")
	}

	err = c.kvPut(connString)
	if err != nil {
		return errors.Wrap(err, "storing connection string in database")
	}

	if !c.validate(connString) {
		return errors.Wrap(err, "failed validating connection string")
	}
	cfg.ConnString = connString
	return nil
}

func (c *Credentialer) createDatabase(cfg Config) (string, error) {
	var (
		user     = cfg.Name
		database = cfg.Name
	)
	password, err := genPass()
	if err != nil {
		return "", err
	}

	rootCredentials, err := parse(cfg.RootConnString)
	if err != nil {
		return "", errors.Wrap(err, "parse root connection string")
	}

	db, err := sql.Open("mysql", cfg.RootConnString)
	if err != nil {
		return "", errors.Wrap(err, "connecting to database as root user")
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET = 'utf8';", database))
	if err != nil {
		return "", errors.Wrap(err, "creating database")
	}
	_, err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'", database, user, password))
	if err != nil {
		return "", errors.Wrap(err, "granting privileges on database")
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, rootCredentials.Host, database), nil
}

// kvPut encodes the given credentials to a string, and stores it in the KV service.
func (c *Credentialer) kvPut(connString string) error {
	_, err := c.KV.Put(&api.KVPair{Key: c.kvKey, Value: []byte(connString)}, nil)
	return err
}

func (c *Credentialer) kvRead() (string, error) {
	pair, _, err := c.KV.Get(c.kvKey, nil)
	if err != nil || pair == nil {
		return "", err
	}
	if !reConnString.Match(pair.Value) {
		c.Log.Warn("credentials don't match the pattern")
		return "", c.kvReset()
	}
	return string(pair.Value), nil
}

// kvReset removes the credentials information from consul.
func (c *Credentialer) kvReset() error {
	c.Log.Warn("remove credentials from kv store")
	_, err := c.KV.Delete(c.kvKey, nil)
	return err
}

func (c *Credentialer) validate(connString string) bool {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return false
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == mysqlerr.ER_ACCESS_DENIED_ERROR {
			return false
		}
	}
	return true
}

type credentials struct {
	User     string
	Password string
	Host     string
	Database string
}

// Parse parses a connection string to a config and credentials structs
func parse(connString string) (*credentials, error) {
	match := reConnString.FindStringSubmatch(connString)
	if len(match) < 5 {
		return nil, fmt.Errorf("connection string %s did not match expected pattern", connString)
	}

	return &credentials{
		User:     match[1],
		Password: match[2],
		Host:     match[3],
		Database: match[4],
	}, nil
}

func genPass() (string, error) {
	b := make([]byte, passwordLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
