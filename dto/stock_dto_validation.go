package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (h StockRequest) Validate() error {
	if err := validation.ValidateStruct(&h,
		validation.Field(&h.Name, validation.Required),
		validation.Field(&h.StockCategory, validation.Required),
		validation.Field(&h.Unit, validation.Required),
		validation.Field(&h.Location, validation.Required),
	); err != nil {
		return err
	}

	// validate type
	if err := stockCategoryValidation(h.StockCategory); err != nil {
		return err
	}

	return nil
}

func (h StockEditRequest) Validate() error {
	if err := validation.ValidateStruct(&h,
		validation.Field(&h.Name, validation.Required),
		validation.Field(&h.FilterTimestamp, validation.Required),
		validation.Field(&h.StockCategory, validation.Required),
		validation.Field(&h.Unit, validation.Required),
		validation.Field(&h.Location, validation.Required),
	); err != nil {
		return err
	}

	// validate type
	if err := stockCategoryValidation(h.StockCategory); err != nil {
		return err
	}

	return nil
}

func (h StockChangeRequest) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Qty, validation.Required),
		validation.Field(&h.Note, validation.Required),
	)
}
