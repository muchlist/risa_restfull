package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/muchlist/risa_restfull/constants/branches"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func roleValidation(rolesIn []string) error {
	if len(rolesIn) > 0 {
		if !sfunc.ValueInSliceIsAvailable(rolesIn, roles.GetRolesAvailable()) {
			return errors.New(fmt.Sprintf("role yang dimasukkan tidak tersedia. gunakan %s", roles.GetRolesAvailable()))
		}
	}
	return nil
}

func branchValidation(branch string) error {

	if !sfunc.InSlice(branch, branches.GetBranchesAvailable()) {
		return errors.New(fmt.Sprintf("branch yang dimasukkan tidak tersedia. gunakan %s", branches.GetBranchesAvailable()))
	}

	return nil
}

//Validate input
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
	if err := branchValidation(u.Branch); err != nil {
		return err
	}

	return nil
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
	if err := branchValidation(u.Branch); err != nil {
		return err
	}

	return nil
}

//Validate input
func (u UserLoginRequest) Validate() error {
	return validation.ValidateStruct(&u,
		//validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.ID, validation.Required),
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
