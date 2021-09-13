package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Validate memvalidasi input altai saat update child item
func (cc ConfigCheckItemUpdateRequest) Validate() error {
	return validation.ValidateStruct(&cc,
		validation.Field(&cc.ParentID, validation.Required),
		validation.Field(&cc.ChildID, validation.Required),
	)
}
