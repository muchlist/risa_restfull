package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (v VendorCheckItemUpdateRequest) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.ParentID, validation.Required),
		validation.Field(&v.ChildID, validation.Required),
	)
}
