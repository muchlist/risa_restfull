package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (i ImproveRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Title, validation.Required),
	)
}

func (i ImproveEditRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.FilterTimeStamp, validation.Required),
		validation.Field(&i.Title, validation.Required),
	)
}
