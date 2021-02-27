package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/muchlist/risa_restfull/constants/stock_category"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func stockCategoryValidation(stockCategory string) error {
	if !sfunc.InSlice(stockCategory, stock_category.GetStockCategoryAvailable()) {
		return errors.New(fmt.Sprintf("Kategory yang dimasukkan tidak tersedia. gunakan %s", stock_category.GetStockCategoryAvailable()))
	}
	return nil
}

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
