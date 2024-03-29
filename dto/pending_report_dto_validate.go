package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/muchlist/risa_restfull/constants/ba"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func (pr PendingReportRequest) Validate() error {
	err := validation.ValidateStruct(&pr,
		validation.Field(&pr.Number, validation.Required),
		validation.Field(&pr.Title, validation.Required),
		validation.Field(&pr.Location, validation.Required),
	)
	if err != nil {
		return err
	}

	if len(pr.Descriptions) > 0 {
		for _, desc := range pr.Descriptions {
			err := validation.ValidateStruct(&desc,
				validation.Field(&desc.Description, validation.Required),
				validation.Field(&desc.DescriptionType, validation.Required),
				validation.Field(&desc.Position, validation.Required),
			)
			if err != nil {
				return err
			}
			if !sfunc.InSlice(desc.DescriptionType, ba.GetDescTypeAvailable()) {
				return errors.New(fmt.Sprintf("Desc type yang dimasukkan tidak tersedia, gunakan %v", ba.GetDescTypeAvailable()))
			}
		}
	}

	if len(pr.Equipments) > 0 {
		for _, equip := range pr.Equipments {
			err := validation.ValidateStruct(&equip,
				validation.Field(&equip.ID, validation.Required),
				validation.Field(&equip.Description, validation.Required),
				validation.Field(&equip.EquipmentName, validation.Required),
				validation.Field(&equip.Qty, validation.Required),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pr PendingReportEditRequest) Validate() error {
	err := validation.ValidateStruct(&pr,
		validation.Field(&pr.FilterTimestamp, validation.Required),
		validation.Field(&pr.Number, validation.Required),
		validation.Field(&pr.Title, validation.Required),
		validation.Field(&pr.Location, validation.Required),
	)
	if err != nil {
		return err
	}

	if len(pr.Descriptions) > 0 {
		for _, desc := range pr.Descriptions {
			err := validation.ValidateStruct(&desc,
				validation.Field(&desc.Description, validation.Required),
				validation.Field(&desc.DescriptionType, validation.Required),
				validation.Field(&desc.Position, validation.Required),
			)
			if err != nil {
				return err
			}
			if !sfunc.InSlice(desc.DescriptionType, ba.GetDescTypeAvailable()) {
				return errors.New(fmt.Sprintf("Desc type yang dimasukkan tidak tersedia, gunakan %v", ba.GetDescTypeAvailable()))
			}
		}
	}

	if len(pr.Equipments) > 0 {
		for _, equip := range pr.Equipments {
			err := validation.ValidateStruct(&equip,
				validation.Field(&equip.ID, validation.Required),
				validation.Field(&equip.Description, validation.Required),
				validation.Field(&equip.EquipmentName, validation.Required),
				validation.Field(&equip.Qty, validation.Required),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pr PendingReportTempOneRequest) Validate() error {
	err := validation.ValidateStruct(&pr,
		validation.Field(&pr.Number, validation.Required),
		validation.Field(&pr.Title, validation.Required),
		validation.Field(&pr.Location, validation.Required),
	)
	if err != nil {
		return err
	}

	if len(pr.Equipments) > 0 {
		for _, equip := range pr.Equipments {
			err := validation.ValidateStruct(&equip,
				validation.Field(&equip.ID, validation.Required),
				validation.Field(&equip.Description, validation.Required),
				validation.Field(&equip.EquipmentName, validation.Required),
				validation.Field(&equip.Qty, validation.Required),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
