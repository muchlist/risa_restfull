package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func categoryValidation(cat string) error {
	if !sfunc.InSlice(cat, category.GetCategoryAvailable()) {
		return errors.New(fmt.Sprintf("Category yang dimasukkan tidak tersedia. gunakan %s", category.GetCategoryAvailable()))
	}
	return nil
}

func (g GenUnitPingStateRequest) Validate() error {
	if err := validation.ValidateStruct(&g,
		validation.Field(&g.Category, validation.Required),
		validation.Field(&g.PingCode, validation.Min(0), validation.Max(2)),
	); err != nil {
		return err
	}

	// validate type

	// validate location
	if err := categoryValidation(g.Category); err != nil {
		return err
	}

	return nil
}
