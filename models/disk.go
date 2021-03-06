// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Disk disk
// swagger:model Disk
type Disk struct {

	// created at
	CreatedAt string `json:"createdAt,omitempty"`

	// free capacity m b
	FreeCapacityMB int64 `json:"freeCapacityMB,omitempty"`

	// hostname
	Hostname string `json:"hostname,omitempty" gorm:"primary_key" query:"filter,sort"`

	// id
	// Required: true
	ID *int64 `json:"id"`

	// media type
	MediaType string `json:"mediaType,omitempty"`

	// model
	Model string `json:"model,omitempty"`

	// path
	// Required: true
	Path *string `json:"path"`

	// serial
	Serial string `json:"serial,omitempty"`

	// total capacity m b
	TotalCapacityMB int64 `json:"totalCapacityMB,omitempty"`

	// updated at
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// Validate validates this disk
func (m *Disk) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateMediaType(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validatePath(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Disk) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

var diskTypeMediaTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["SSD","HDD","NVMe"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		diskTypeMediaTypePropEnum = append(diskTypeMediaTypePropEnum, v)
	}
}

const (
	// DiskMediaTypeSSD captures enum value "SSD"
	DiskMediaTypeSSD string = "SSD"
	// DiskMediaTypeHDD captures enum value "HDD"
	DiskMediaTypeHDD string = "HDD"
	// DiskMediaTypeNVme captures enum value "NVMe"
	DiskMediaTypeNVme string = "NVMe"
)

// prop value enum
func (m *Disk) validateMediaTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, diskTypeMediaTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Disk) validateMediaType(formats strfmt.Registry) error {

	if swag.IsZero(m.MediaType) { // not required
		return nil
	}

	// value enum
	if err := m.validateMediaTypeEnum("mediaType", "body", m.MediaType); err != nil {
		return err
	}

	return nil
}

func (m *Disk) validatePath(formats strfmt.Registry) error {

	if err := validate.Required("path", "body", m.Path); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Disk) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Disk) UnmarshalBinary(b []byte) error {
	var res Disk
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
