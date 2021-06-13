package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Validate input
func (u UserRequest) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Branch, validation.Required),
		validation.Field(&u.Roles, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	); err != nil {
		return err
	}

	// validate role
	if err := roleValidation(u.Roles); err != nil {
		return err
	}
	// validate branch
	return branchValidation(u.Branch)
}

func (u UserEditRequest) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.TimestampFilter, validation.Required),
		validation.Field(&u.Branch, validation.Required),
		validation.Field(&u.Roles, validation.Required),
	); err != nil {
		return err
	}
	if err := roleValidation(u.Roles); err != nil {
		return err
	}
	return branchValidation(u.Branch)
}

// Validate input
func (u UserLoginRequest) Validate() error {
	return validation.ValidateStruct(&u,
		// validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.ID, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	)
}

// Validate input
func (u UserChangePasswordRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
		validation.Field(&u.NewPassword, validation.Required, validation.Length(3, 20)),
	)
}

// Validate input
func (u UserRefreshTokenRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.RefreshToken, validation.Required),
	)
}

// Validate input
func (u UserUpdateFcmRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FcmToken, validation.Required),
	)
}
