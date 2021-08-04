package dto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/muchlist/risa_restfull/constants/enum"
	"path/filepath"
	"strings"
)

const (
	jpgExtension  = ".jpg"
	pngExtension  = ".png"
	jpegExtension = ".jpeg"
)

// Validate input
func (h HistoryRequest) Validate() error {
	var errorList []string

	if err := validation.ValidateStruct(&h,
		validation.Field(&h.ParentID, validation.Required),
		validation.Field(&h.Status, validation.Required),
		validation.Field(&h.Problem, validation.Required),
		validation.Field(&h.CompleteStatus, validation.Max(enum.HComplete), validation.Min(0)),
	); err != nil {
		errorList = append(errorList, err.Error())
	}

	if h.Image != "" {
		// cek ekstensi
		fileExtension := strings.ToLower(filepath.Ext(h.Image))
		if !(fileExtension == jpgExtension || fileExtension == pngExtension || fileExtension == jpegExtension) {
			errorList = append(errorList, "ekstensi file image salah")
		}
	}

	if len(errorList) != 0 {
		var errorString strings.Builder
		for _, v := range errorList {
			errorString.WriteString(v + ". ")
		}
		return errors.New(errorString.String())
	}

	return nil
}

func (h HistoryEditRequest) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.FilterTimestamp, validation.Required),
		validation.Field(&h.Status, validation.Required),
		validation.Field(&h.Problem, validation.Required),
		validation.Field(&h.CompleteStatus, validation.Max(enum.HComplete), validation.Min(0)),
	)
}
