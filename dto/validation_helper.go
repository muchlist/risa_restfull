package dto

import (
	"errors"
	"fmt"
	"github.com/muchlist/risa_restfull/constants/branches"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/check_type"
	"github.com/muchlist/risa_restfull/constants/hw_type"
	"github.com/muchlist/risa_restfull/constants/location"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/constants/stock_category"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func locationValidation(loc string) error {
	if !sfunc.InSlice(loc, location.GetLocationAvailable()) {
		return errors.New(fmt.Sprintf("Lokasi yang dimasukkan tidak tersedia. gunakan %s", location.GetLocationAvailable()))
	}
	return nil
}

func categoryValidation(cat string) error {
	if !sfunc.InSlice(cat, category.GetCategoryAvailable()) {
		return errors.New(fmt.Sprintf("Category yang dimasukkan tidak tersedia. gunakan %s", category.GetCategoryAvailable()))
	}
	return nil
}

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

func stockCategoryValidation(stockCategory string) error {
	if !sfunc.InSlice(stockCategory, stock_category.GetStockCategoryAvailable()) {
		return errors.New(fmt.Sprintf("Category yang dimasukkan tidak tersedia. gunakan %s", stock_category.GetStockCategoryAvailable()))
	}
	return nil
}

func checkTypeValidation(checkType string) error {
	if !sfunc.InSlice(checkType, check_type.GetCheckTypeAvailable()) {
		return errors.New(fmt.Sprintf("Tipe yang dimasukkan tidak tersedia. gunakan %s", check_type.GetCheckTypeAvailable()))
	}
	return nil
}

func cctvTypeValidation(cctvType string) error {
	if !sfunc.InSlice(cctvType, hw_type.GetCctvTypeAvailable()) {
		return errors.New(fmt.Sprintf("Tipe yang dimasukkan tidak tersedia. gunakan %s", hw_type.GetCctvTypeAvailable()))
	}
	return nil
}
