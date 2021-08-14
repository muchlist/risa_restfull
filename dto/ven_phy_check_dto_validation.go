package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (vpc VenPhyCheckItemUpdateRequest) Validate() error {
	return validation.ValidateStruct(&vpc,
		validation.Field(&vpc.ParentID, validation.Required),
		validation.Field(&vpc.ChildID, validation.Required),
	)
}
