package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//func stockCategoryValidation(stockCategory string) error {
//	if !sfunc.InSlice(stockCategory, stock_category.GetStockCategoryAvailable()) {
//		return errors.New(fmt.Sprintf("Kategory yang dimasukkan tidak tersedia. gunakan %s", hw_type.GetCctvTypeAvailable()))
//	}
//	return nil
//}

func (h StockRequest) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Name, validation.Required),
		validation.Field(&h.StockCategory, validation.Required),
		validation.Field(&h.Unit, validation.Required),
		validation.Field(&h.Location, validation.Required),
	)
}
