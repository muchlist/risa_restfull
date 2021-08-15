package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (ac AltaiPhyCheckItemUpdateRequest) Validate() error {
	return validation.ValidateStruct(&ac,
		validation.Field(&ac.ParentID, validation.Required),
		validation.Field(&ac.ChildID, validation.Required),
	)
}
