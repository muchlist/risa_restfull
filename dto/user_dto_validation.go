package dto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/muchlist/risa_restfull/roles_const"
	"github.com/muchlist/risa_restfull/utils"
)

func roleValidation(roles []string) error {
	if len(roles) > 0 {
		if !utils.ValueInSliceIsAvailable(roles, roles_const.GetRolesAvailable()) {
			return errors.New("the role entered is not available")
		}
	}
	return nil
}

//Validate input
func (u UserRequest) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	); err != nil {
		return err
	}

	if err := roleValidation(u.Roles); err != nil {
		return err
	}

	return nil
}

func (u UserEditRequest) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.TimestampFilter, validation.Required),
	); err != nil {
		return err
	}

	if err := roleValidation(u.Roles); err != nil {
		return err
	}

	return nil
}

//Validate input
func (u UserLoginRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	)
}

//Validate input
func (u UserChangePasswordRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
		validation.Field(&u.NewPassword, validation.Required, validation.Length(3, 20)),
	)
}

//Validate input
func (u UserRefreshTokenRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.RefreshToken, validation.Required),
	)
}
