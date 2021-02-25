package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/muchlist/risa_restfull/constants/hw_type"
	"github.com/muchlist/risa_restfull/constants/location"
	"github.com/muchlist/risa_restfull/utils"
)

func cctvTypeValidation(cctvType string) error {
	if !utils.InSlice(cctvType, hw_type.GetCctvTypeAvailable()) {
		return errors.New(fmt.Sprintf("Tipe yang dimasukkan tidak tersedia. gunakan %s", hw_type.GetCctvTypeAvailable()))
	}
	return nil
}

func locationValidation(loc string) error {
	if !utils.InSlice(loc, location.GetLocationAvailable()) {
		return errors.New(fmt.Sprintf("Lokasi yang dimasukkan tidak tersedia. gunakan %s", location.GetLocationAvailable()))
	}
	return nil
}

func (c CctvRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Location, validation.Required),
		validation.Field(&c.Type, validation.Required),
	); err != nil {
		return err
	}

	// validate role
	if err := cctvTypeValidation(c.Type); err != nil {
		return err
	}

	// validate location
	if err := locationValidation(c.Type); err != nil {
		return err
	}

	return nil
}
