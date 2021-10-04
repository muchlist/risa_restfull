package service

import (
	"context"
	"fmt"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/dao/cctvdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/vendorcheckdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func NewVendorCheckService(
	vendorCheckDao vendorcheckdao.CheckVendorDaoAssumer,
	genUnitDao genunitdao.GenUnitDaoAssumer,
	cctvDao cctvdao.CctvDaoAssumer,
	histService HistoryServiceAssumer,
) VendorCheckServiceAssumer {
	return &vendorCheckService{
		daoC:        vendorCheckDao,
		daoG:        genUnitDao,
		daoCTV:      cctvDao,
		servHistory: histService,
	}
}

type vendorCheckService struct {
	daoC        vendorcheckdao.CheckVendorDaoAssumer
	daoG        genunitdao.GenUnitDaoAssumer
	daoCTV      cctvdao.CctvLoader
	servHistory HistoryServiceAssumer
}
type VendorCheckServiceAssumer interface {
	InsertVendorCheck(user mjwt.CustomClaim) (*string, rest_err.APIError)
	DeleteVendorCheck(user mjwt.CustomClaim, id string) rest_err.APIError
	GetVendorCheckByID(vendorCheckID string, branchIfSpecific string) (*dto.VendorCheck, rest_err.APIError)
	FindVendorCheck(branch string, filter dto.FilterTimeRangeLimit) ([]dto.VendorCheck, rest_err.APIError)
	UpdateVendorCheckItem(user mjwt.CustomClaim, input dto.VendorCheckItemUpdateRequest) (*dto.VendorCheck, rest_err.APIError)
	BulkUpdateVendorItem(user mjwt.CustomClaim, inputs []dto.VendorCheckItemUpdateRequest) (string, rest_err.APIError)
	FinishCheck(user mjwt.CustomClaim, detailID string) (*dto.VendorCheck, rest_err.APIError)
}

func (c *vendorCheckService) InsertVendorCheck(user mjwt.CustomClaim) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	// ambil cctv genUnit item berdasarkan cabang yang di input
	// mendapatkan data cases
	genItems, err := c.daoG.FindUnit(dto.GenUnitFilter{
		Branch:   user.Branch,
		Category: category.Cctv,
		Pings:    false,
	})
	if err != nil {
		return nil, err
	}

	// ambil cctv untuk mendapatkan data lokasi
	// cctvItems sudah sorted berdasarkan lokasi sedangkan genItems tidak
	cctvItems, err := c.daoCTV.FindCctv(context.TODO(), dto.FilterBranchLocIPNameDisable{
		FilterBranch: user.Branch,
	})
	if err != nil {
		return nil, err
	}

	// ubah cctv genUnit menjadi map
	genItemsMap := make(map[string]dto.GenUnitResponse)
	for _, item := range genItems {
		genItemsMap[item.ID] = item
	}

	// kategorikan vendorCheck item menggunakan perulangan cctvItems
	// karena cctvItems sudah sorted
	var vendorCheckItems []dto.VendorCheckItemEmbed
	for _, v := range cctvItems {
		cctvInfoFromGenUnit := genItemsMap[v.ID.Hex()]

		// jika didalam semua case yang ada di cctv tersebut ada tag #isBlur maka kita anggap cctvnya blur
		// dan tidak mati
		isOffline := cctvInfoFromGenUnit.CasesSize != 0
		isBlur := strings.Contains(fmt.Sprintf("%v", cctvInfoFromGenUnit.Cases), "#isBlur")
		if isBlur {
			isOffline = false
		}

		// pengecekan secara virtual default sudah tercek semua di waktu pembuatan
		vendorCheckItems = append(vendorCheckItems, dto.VendorCheckItemEmbed{
			ID:        v.ID.Hex(),
			Name:      v.Name,
			Location:  v.Location,
			CheckedAt: timeNow,
			CheckedBy: user.Name,
			IsChecked: true,
			IsBlur:    isBlur,
			IsOffline: isOffline,
		})
	}

	data := dto.VendorCheck{
		CreatedAt:        timeNow,
		CreatedBy:        user.Name,
		CreatedByID:      user.Identity,
		UpdatedAt:        timeNow,
		UpdatedBy:        user.Name,
		UpdatedByID:      user.Identity,
		Branch:           user.Branch,
		TimeStarted:      timeNow,
		TimeEnded:        0,
		IsVirtualCheck:   true,
		IsFinish:         false,
		IsCheckByVendor:  sfunc.InSlice(roles.RoleVendor, user.Roles),
		VendorCheckItems: vendorCheckItems,
		Note:             "",
	}

	// DB
	insertedID, err := c.daoC.InsertCheck(data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (c *vendorCheckService) DeleteVendorCheck(user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := c.daoC.DeleteCheck(dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *vendorCheckService) UpdateVendorCheckItem(user mjwt.CustomClaim, input dto.VendorCheckItemUpdateRequest) (*dto.VendorCheck, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// DB
	data := dto.VendorCheckItemUpdate{
		FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
			FilterParentID: parentOid,
			FilterChildID:  input.ChildID,
			FilterBranch:   user.Branch,
		},
		CheckedAt: timeNow,
		CheckedBy: user.Name,
		IsChecked: input.IsChecked,
		IsBlur:    input.IsBlur,
		IsOffline: input.IsOffline,
	}
	vendorCheck, err := c.daoC.UpdateCheckItem(data)
	if err != nil {
		return nil, err
	}
	return vendorCheck, nil
}

func (c *vendorCheckService) BulkUpdateVendorItem(user mjwt.CustomClaim, inputs []dto.VendorCheckItemUpdateRequest) (string, rest_err.APIError) {
	if len(inputs) == 0 {
		return "", rest_err.NewBadRequestError("tidak ada perubahan")
	}

	parentOid, errT := primitive.ObjectIDFromHex(inputs[0].ParentID)
	if errT != nil {
		return "", rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	inputDatas := make([]dto.VendorCheckItemUpdate, len(inputs))
	for i, input := range inputs {
		inputDatas[i] = dto.VendorCheckItemUpdate{
			FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
				FilterParentID: parentOid,
				FilterChildID:  input.ChildID,
				FilterBranch:   user.Branch,
			},
			CheckedAt: timeNow,
			CheckedBy: user.Name,
			IsChecked: input.IsChecked,
			IsBlur:    input.IsBlur,
			IsOffline: input.IsOffline,
		}
	}
	result, err := c.daoC.BulkUpdateItem(inputDatas)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d data telah diubah", result), nil
}

func (c *vendorCheckService) FinishCheck(user mjwt.CustomClaim, detailID string) (*dto.VendorCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(detailID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// 1. cek cctv existing, untuk mendapatkan keterangan apakah ada case
	genItems, err := c.daoG.FindUnit(dto.GenUnitFilter{
		Branch:   user.Branch,
		Category: category.Cctv,
		Pings:    false,
	})
	if err != nil {
		return nil, err
	}

	// 2. cctv yang memiliki case tersebut di exclude dari update berikutnya
	var cctvExcludeID []string
	for _, cctv := range genItems {
		if cctv.CasesSize != 0 {
			cctvExcludeID = append(cctvExcludeID, cctv.ID)
		}
	}

	// 3. filter items check yang memiliki is_offline atau is_blur true
	var cctvBlurID []string
	var cctvOfflineID []string
	cctvChecklistDetail, err := c.daoC.GetCheckByID(oid, user.Branch)
	if err != nil {
		return nil, err
	}
	for _, cctv := range cctvChecklistDetail.VendorCheckItems {
		if sfunc.InSlice(cctv.ID, cctvExcludeID) {
			continue
		}
		if cctv.IsBlur {
			cctvBlurID = append(cctvBlurID, cctv.ID)
		}
		if cctv.IsOffline {
			cctvOfflineID = append(cctvOfflineID, cctv.ID)
		}
	}

	// send to background
	go func() {
		// Insert History isBlur
		if len(cctvBlurID) != 0 {
			for _, cctvID := range cctvBlurID {
				_, _ = c.servHistory.InsertHistory(user, dto.HistoryRequest{
					ParentID:       cctvID,
					Status:         "",
					Problem:        "Display CCTV buram #isBlur",
					ProblemResolve: "",
					CompleteStatus: enum.HPending,
					DateStart:      timeNow,
					DateEnd:        0,
					Tag:            []string{},
				})
			}
		}

		// Insert History isoffline
		if len(cctvOfflineID) != 0 {
			for _, cctvID := range cctvOfflineID {
				_, _ = c.servHistory.InsertHistory(user, dto.HistoryRequest{
					ParentID:       cctvID,
					Status:         "",
					Problem:        "CCTV Offline",
					ProblemResolve: "",
					CompleteStatus: enum.HProgress,
					DateStart:      timeNow,
					DateEnd:        0,
					Tag:            []string{},
				})
			}
		}
	}()

	// 7. tandai isFinish true dan end_date ke waktu sekarang
	cctvChecklistDetail, err = c.daoC.EditCheck(dto.VendorCheckEdit{
		FilterIDBranch: dto.FilterIDBranch{
			FilterID:     oid,
			FilterBranch: user.Branch,
		},
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		TimeStarted: cctvChecklistDetail.TimeStarted,
		TimeEnded:   timeNow,
		IsFinish:    true,
		Note:        "",
	})
	if err != nil {
		return nil, err
	}

	return cctvChecklistDetail, nil
}

func (c *vendorCheckService) GetVendorCheckByID(vendorCheckID string, branchIfSpecific string) (*dto.VendorCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(vendorCheckID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	vendorCheck, err := c.daoC.GetCheckByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return vendorCheck, nil
}

func (c *vendorCheckService) FindVendorCheck(branch string, filter dto.FilterTimeRangeLimit) ([]dto.VendorCheck, rest_err.APIError) {
	vendorCheckList, err := c.daoC.FindCheck(branch, filter, false)
	if err != nil {
		return []dto.VendorCheck{}, err
	}
	return vendorCheckList, nil
}
