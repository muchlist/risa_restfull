package dto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c OtherRequest) Validate() error {
	var errorList []string

	if err := validation.ValidateStruct(&c,
		validation.Field(&c.SubCategory, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
	); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate category
	if err := subCategoryValidation(c.SubCategory); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate location
	if err := locationValidation(c.Location); err != nil {
		errorList = append(errorList, err.Error())
	}

	if len(errorList) != 0 {
		errorString := ""
		for _, v := range errorList {
			errorString = errorString + v + ". "
		}
		return errors.New(errorString)
	}

	return nil
}

func (c OtherEditRequest) Validate() error {
	var errorList []string

	if err := validation.ValidateStruct(&c,
		validation.Field(&c.FilterSubCategory, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.FilterTimestamp, validation.Required),
	); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate subcategory
	if err := subCategoryValidation(c.FilterSubCategory); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate location
	if err := locationValidation(c.Location); err != nil {
		errorList = append(errorList, err.Error())
	}

	if len(errorList) != 0 {
		errorString := ""
		for _, v := range errorList {
			errorString = errorString + v + ". "
		}
		return errors.New(errorString)
	}

	return nil
}
