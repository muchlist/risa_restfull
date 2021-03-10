package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c CctvRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.Type, validation.Required),
	); err != nil {
		return err
	}

	// validate type
	if err := cctvTypeValidation(c.Type); err != nil {
		return err
	}

	// validate location
	return locationValidation(c.Location)
}

func (c CctvEditRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.Type, validation.Required),
		validation.Field(&c.FilterTimestamp, validation.Required),
	); err != nil {
		return err
	}

	// validate type
	if err := cctvTypeValidation(c.Type); err != nil {
		return err
	}

	// validate location
	return locationValidation(c.Location)
}
