// Code generated by go-swagger; DO NOT EDIT.

package disk

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDiskByIDParams creates a new DiskByIDParams object
// with the default values initialized.
func NewDiskByIDParams() DiskByIDParams {
	var ()
	return DiskByIDParams{}
}

// DiskByIDParams contains all the bound params for the disk by Id operation
// typically these are obtained from a http.Request
//
// swagger:parameters diskById
type DiskByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The ID of the requested disk
	  Required: true
	  In: path
	*/
	DiskID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *DiskByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rDiskID, rhkDiskID, _ := route.Params.GetOK("disk_id")
	if err := o.bindDiskID(rDiskID, rhkDiskID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DiskByIDParams) bindDiskID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.DiskID = raw

	return nil
}
