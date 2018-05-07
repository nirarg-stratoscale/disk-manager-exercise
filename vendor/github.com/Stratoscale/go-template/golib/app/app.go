package app

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Config is application configuration
type Config struct {
	ListenAddr string `envconfig:"LISTEN_ADDR" default:":80"`
	PProf      PProfConfig
	Log        LogConfig
}

type PProfConfig struct {
	// Enabled indicates if profiling data should be enabled.
	Enabled bool `envconfig:"PPROF_ENABLED"`
	// Addr is the address that pprof will listen on
	Addr   string `envconfig:"PPROF_ADDR" default:"0.0.0.0:6060"`
	Prefix string `envconfig:"PPROF_PREFIX" default:"/debug/pprof/"`
}

// App is a building block for microservice applications.
// It performs service registration, logging configuration and runs HTTP server.
type App struct {
	Config
	Log *logrus.Entry

	ctx    context.Context
	cancel context.CancelFunc
}

// New returns a new application
func New(cfg Config, name string) *App {
	app := new(App)
	app.Config = cfg
	app.Log = newLog(app.Config.Log, name)
	app.ctx, app.cancel = context.WithCancel(context.Background())

	if app.PProf.Enabled {
		app.runDebugTools()
	}

	return app
}

// Context returns the application's context
func (a *App) Context() context.Context {
	return a.ctx
}

// RunHTTP starts asynchronously server in HTTP.
func (a *App) RunHTTP(h http.Handler) {
	//h = a.withLogger(h)
	h = a.withHealthCheck(h)

	server := &http.Server{
		Addr:     a.ListenAddr,
		Handler:  h,
		ErrorLog: log.New(a.Log.Writer(), "", 0),
	}

	a.Log.Info("Starting HTTP")

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.Log.WithError(err).Fatal("HTTP server failed to start")
		}
	}()

	go server.Shutdown(a.ctx)
}

func (a *App) withHealthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// runDebugTools runs tools like pprof based on configuration.
func (a *App) runDebugTools() {
	h := http.NewServeMux()
	h.Handle(a.PProf.Prefix, http.HandlerFunc(pprof.Index))
	h.Handle(a.PProf.Prefix+"cmdline", http.HandlerFunc(pprof.Cmdline))
	h.Handle(a.PProf.Prefix+"profile", http.HandlerFunc(pprof.Profile))
	h.Handle(a.PProf.Prefix+"symbol", http.HandlerFunc(pprof.Symbol))
	h.Handle(a.PProf.Prefix+"trace", http.HandlerFunc(pprof.Trace))

	a.Log.Infof("Starting PPROF on %s", a.PProf.Addr)

	server := &http.Server{Addr: a.ListenAddr, Handler: h}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.FailOnError(err, "HTTP server failed to start")
		}
	}()

	go server.Shutdown(a.ctx)
}

// WaitForSignal blocks until SIGINT or SIGTERM is received, then exits process.
func (a *App) WaitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigs:
		a.Log.Info("Exiting on signal", "signal", sig)
		a.cancel()
	case <-a.ctx.Done():
	}
}

func (a *App) Close() {
	a.cancel()
}

// FailOnError fails the application with reason if err is not nil.
func (a *App) FailOnError(err error, reason string) {
	if err != nil {
		a.Log.WithError(err).Fatal(reason)
	}
}
