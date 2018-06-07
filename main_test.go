package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Stratoscale/disk-manager-exercise/restapi"
	"github.com/Stratoscale/go-template/golib/middleware"
	"github.com/Stratoscale/go-template/golib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var log = testutil.Log()

func TestHTTPHandler(t *testing.T) {
	t.Parallel()

	// declare the test cases
	tests := []struct {
		// name is the name of the test
		name string

		// a set of parameters to run the test with
		req  *http.Request
		role string

		// prepare prepares the mock objects before performing the tested function
		prepare func(*testing.T, *restapi.MockDiskAPI)

		// a set of results that we expect
		wantCode int
		wantBody []byte
	}{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// create all mocks and variables for the test
			var (
				diskMock restapi.MockDiskAPI
				resp     = httptest.NewRecorder()
			)

			h, err := restapi.Handler(restapi.Config{
				DiskAPI:         &diskMock,
				InnerMiddleware: middleware.Policy,
				Logger:          log.Debugf,
			})
			require.Nil(t, err)

			// prepare mocks
			if tt.prepare != nil {
				tt.prepare(t, &diskMock)
			}

			// prepare the request for sending
			tt.req.Header.Set("Content-Type", "application/json")
			tt.req.Header.Set(middleware.KeyRoles, tt.role)

			// run the http routing with the produced request
			h.ServeHTTP(resp, tt.req)

			t.Logf("Got response for request %s %s: %d %s", tt.req.Method, tt.req.URL, resp.Code, resp.Body.String())

			// assert the response expectations
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantBody != nil {
				assert.JSONEq(t, string(tt.wantBody), resp.Body.String())
			}

			diskMock.AssertExpectations(t)
		})
	}
}
