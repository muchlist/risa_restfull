package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/muchlist/risa_restfull/constants/location"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func (c CheckItemRequest) Validate() error {
	var errorList []string

	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.Type, validation.Required),
	); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate type
	if err := checkTypeValidation(c.Type); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate location
	if err := locationValidation(c.Location); err != nil {
		errorList = append(errorList, err.Error())
	}

	if c.Shifts != nil {
		shiftsAvailable := []int{1, 2, 3}
		if !sfunc.ValueIntInSliceIsAvailable(c.Shifts, shiftsAvailable) {
			errorList = append(errorList, fmt.Sprintf("Lokasi yang dimasukkan tidak tersedia. Gunakan %s", location.GetLocationAvailable()))
		}
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

func (c CheckItemEditRequest) Validate() error {
	var errorList []string
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.FilterTimestamp, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.Type, validation.Required),
	); err != nil {
		return err
	}

	// validate type
	if err := checkTypeValidation(c.Type); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate location
	if err := locationValidation(c.Location); err != nil {
		errorList = append(errorList, err.Error())
	}

	if c.Shifts != nil {
		shiftsAvailable := []int{1, 2, 3}
		if !sfunc.ValueIntInSliceIsAvailable(c.Shifts, shiftsAvailable) {
			errorList = append(errorList, fmt.Sprintf("Lokasi yang dimasukkan tidak tersedia. Gunakan %s", location.GetLocationAvailable()))
		}
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
