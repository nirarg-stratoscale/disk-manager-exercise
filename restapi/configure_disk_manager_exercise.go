// Code generated by go-swagger; DO NOT EDIT.

package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Stratoscale/swagger/auth"
	"github.com/Stratoscale/swagger/query"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Stratoscale/disk-manager-exercise/models"
	"github.com/Stratoscale/disk-manager-exercise/restapi/operations"
	"github.com/Stratoscale/disk-manager-exercise/restapi/operations/disk"
)

//go:generate mockery -name DiskAPI -inpkg

// DiskAPI
type DiskAPI interface {
	DiskByID(ctx context.Context, params disk.DiskByIDParams) middleware.Responder
	ListDisks(ctx context.Context, params disk.ListDisksParams) middleware.Responder
}

// Config is configuration for Handler
type Config struct {
	DiskAPI
	Logger func(string, ...interface{})
	// InnerMiddleware is for the handler executors. These do not apply to the swagger.json document.
	// The middleware executes after routing but before authentication, binding and validation
	InnerMiddleware func(http.Handler) http.Handler
	Auth            auth.Auth
}

// Handler returns an http.Handler given the handler configuration
// It mounts all the business logic implementers in the right routing.
func Handler(c Config) (http.Handler, error) {
	spec, err := loads.Analyzed(swaggerCopy(SwaggerJSON), "")
	if err != nil {
		return nil, fmt.Errorf("analyze swagger: %v", err)
	}
	api := operations.NewDiskManagerExerciseAPI(spec)
	api.ServeError = errors.ServeError
	api.Logger = c.Logger

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.DiskDiskByIDHandler = disk.DiskByIDHandlerFunc(func(params disk.DiskByIDParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = c.Auth.Store(ctx, principal)
		return c.DiskAPI.DiskByID(ctx, params)
	})
	api.DiskListDisksHandler = disk.ListDisksHandlerFunc(func(params disk.ListDisksParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = c.Auth.Store(ctx, principal)
		return c.DiskAPI.ListDisks(ctx, params)
	})
	api.ServerShutdown = func() {}
	return api.Serve(c.InnerMiddleware), nil
}

// Query parse functions for all the models
// Those can be used to extract database query from the http path's query string
var (
	DiskQueryParse = query.MustNewBuilder(&query.Config{Model: models.Disk{}}).ParseRequest
)

// swaggerCopy copies the swagger json to prevent data races in runtime
func swaggerCopy(orig json.RawMessage) json.RawMessage {
	c := make(json.RawMessage, len(orig))
	copy(c, orig)
	return c
}
