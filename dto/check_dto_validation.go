package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (c CheckRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Shift, validation.Required, validation.Min(0), validation.Max(3)),
	)
}

func (c CheckChildUpdateRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ParentID, validation.Required),
		validation.Field(&c.ChildID, validation.Required),
		validation.Field(&c.CompleteStatus, validation.Min(0), validation.Max(4)),
	)
}
