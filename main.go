package main

import (
	"flag"

	"github.com/Stratoscale/disk-manager-exercise/restapi"
	"github.com/Stratoscale/go-template/golib/app"
	"github.com/Stratoscale/go-template/golib/consulutil"
	"github.com/Stratoscale/go-template/golib/dbutil"
	"github.com/Stratoscale/go-template/golib/middleware"
	//"github.com/Stratoscale/golib/consul"
	"github.com/kelseyhightower/envconfig"

	"github.com/Stratoscale/disk-manager-exercise/internal/disk"
)

var options struct {
	App    app.Config
	DB     dbutil.Config
	Consul consulutil.Config
}

func init() {
	flag.Usage = func() { envconfig.Usage("", &options) }
	flag.Parse()
}

func main() {

	err := envconfig.Process("", &options)
	if err != nil {
		panic("processing environment variables")
	}

	a := app.New(options.App, "disk-manager-exercise")

	a.Log.Infof("Options: %+v", options)

	//consulClient, err := consulutil.Client(options.Consul)
	//a.FailOnError(err, "initialize consul client")

	//credentialer := dbutil.Credentialer{
	//	KV:     consulClient.KV(),
	//	Locker: consul.NewLocker(consulClient),
	//	Log:    a.Log.WithField("pkg", "credentials"),
	//}

	//err = credentialer.ConnectionString(&options.DB)
	//a.FailOnError(err, "connection string")

	//db, err := dbutil.Open(options.DB, a.Log.WithField("pkg", "db"))
	//a.FailOnError(err, "initializing database")
	//defer db.Close()

	disk := disk.New(disk.Config{
		//DB:  db,
		Log: a.Log.WithField("pkg", "disk"),
	})

	//err = disk.AutoMigrate()
	//a.FailOnError(err, "migrating disk database")

	h, err := restapi.Handler(restapi.Config{
		DiskAPI:        disk,
		Logger:         a.Log.WithField("pkg", "restapi").Debugf,
		InnerMiddleware: middleware.Policy,
	})
	a.FailOnError(err, "initializing handler")

	a.RunHTTP(h)
	a.WaitForSignal()
}
