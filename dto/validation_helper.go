package dto

import (
	"fmt"
	"github.com/muchlist/risa_restfull/constants/branches"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/checktype"
	"github.com/muchlist/risa_restfull/constants/hwlist"
	"github.com/muchlist/risa_restfull/constants/location"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/constants/stocklist"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func locationValidation(loc string) error {
	if !sfunc.InSlice(loc, location.GetLocationAvailable()) {
		return fmt.Errorf("lokasi yang dimasukkan tidak tersedia. gunakan %s", location.GetLocationAvailable())
	}
	return nil
}

func categoryValidation(cat string) error {
	if !sfunc.InSlice(cat, category.GetCategoryAvailable()) {
		return fmt.Errorf("category yang dimasukkan tidak tersedia. gunakan %s", category.GetCategoryAvailable())
	}
	return nil
}

func roleValidation(rolesIn []string) error {
	if len(rolesIn) > 0 {
		if !sfunc.ValueInSliceIsAvailable(rolesIn, roles.GetRolesAvailable()) {
			return fmt.Errorf("role yang dimasukkan tidak tersedia. gunakan %s", roles.GetRolesAvailable())
		}
	}
	return nil
}

func branchValidation(branch string) error {
	if !sfunc.InSlice(branch, branches.GetBranchesAvailable()) {
		return fmt.Errorf("branch yang dimasukkan tidak tersedia. gunakan %s", branches.GetBranchesAvailable())
	}

	return nil
}

func stockCategoryValidation(stockCategory string) error {
	if !sfunc.InSlice(stockCategory, stocklist.GetStockCategoryAvailable()) {
		return fmt.Errorf("category yang dimasukkan tidak tersedia. gunakan %s", stocklist.GetStockCategoryAvailable())
	}
	return nil
}

func checkTypeValidation(checkType string) error {
	if !sfunc.InSlice(checkType, checktype.GetCheckTypeAvailable()) {
		return fmt.Errorf("tipe yang dimasukkan tidak tersedia. gunakan %s", checktype.GetCheckTypeAvailable())
	}
	return nil
}

func cctvTypeValidation(cctvType string) error {
	if !sfunc.InSlice(cctvType, hwlist.GetCctvTypeAvailable()) {
		return fmt.Errorf("tipe yang dimasukkan tidak tersedia. gunakan %s", hwlist.GetCctvTypeAvailable())
	}
	return nil
}
