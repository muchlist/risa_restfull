package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (g GenUnitPingStateRequest) Validate() error {
	if err := validation.ValidateStruct(&g,
		validation.Field(&g.Category, validation.Required),
		validation.Field(&g.PingCode, validation.Min(0), validation.Max(2)),
	); err != nil {
		return err
	}

	// validate type

	// validate location
	return categoryValidation(g.Category)
}
