package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/muchlist/risa_restfull/constants/enum"
)

// Validate input
func (h HistoryRequest) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.ParentID, validation.Required),
		validation.Field(&h.Status, validation.Required),
		validation.Field(&h.Problem, validation.Required),
		validation.Field(&h.CompleteStatus, validation.Max(enum.HComplete), validation.Min(0)),
	)
}

func (h HistoryEditRequest) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.FilterTimestamp, validation.Required),
		validation.Field(&h.Status, validation.Required),
		validation.Field(&h.Problem, validation.Required),
		validation.Field(&h.CompleteStatus, validation.Max(enum.HComplete), validation.Min(0)),
	)
}
