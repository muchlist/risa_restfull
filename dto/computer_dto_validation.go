package dto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c ComputerRequest) Validate() error {
	var errorList []string

	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.Division, validation.Required),
		validation.Field(&c.Type, validation.Required),
	); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate type
	if err := computerTypeValidation(c.Type); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate location
	if err := locationValidation(c.Location); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate division
	if err := divisionValidation(c.Division); err != nil {
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

func (c ComputerEditRequest) Validate() error {
	var errorList []string

	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.Division, validation.Required),
		validation.Field(&c.Type, validation.Required),
		validation.Field(&c.FilterTimestamp, validation.Required),
	); err != nil {
		errorList = append(errorList, err.Error())
	}

	// validate type
	if err := computerTypeValidation(c.Type); err != nil {
		errorList = append(errorList, err.Error())
	}
	// validate division
	if err := divisionValidation(c.Division); err != nil {
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
