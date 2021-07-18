package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/dao/cctvdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/vendorcheckdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
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
	daoCTV      cctvdao.CctvDaoAssumer
	servHistory HistoryServiceAssumer
}
type VendorCheckServiceAssumer interface {
	InsertVendorCheck(user mjwt.CustomClaim, isVirtualCheck bool) (*string, rest_err.APIError)
	//EditVendorCheck(user mjwt.CustomClaim, vendorCheckID string, input dto.VendorCheckEdit) (*dto.VendorCheck, rest_err.APIError)
	//DeleteVendorCheck(user mjwt.CustomClaim, id string) rest_err.APIError
	//UpdateVendorCheckItem(user mjwt.CustomClaim, input dto.VendorCheckItemUpdate) (*dto.VendorCheck, rest_err.APIError)
	//PutChildImage(user mjwt.CustomClaim, parentID string, childID string, imagePath string) (*dto.VendorCheck, rest_err.APIError)
	//
	//GetVendorCheckByID(vendorCheckID string, branchIfSpecific string) (*dto.VendorCheck, rest_err.APIError)
	//FindVendorCheck(branch string, filter dto.FilterTimeRangeLimit) ([]dto.VendorCheck, rest_err.APIError)
}

func (c *vendorCheckService) InsertVendorCheck(user mjwt.CustomClaim, isVirtualCheck bool) (*string, rest_err.APIError) {
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
	cctvItems, err := c.daoCTV.FindCctv(dto.FilterBranchLocIPNameDisable{
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

		if isVirtualCheck {
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
		} else {
			vendorCheckItems = append(vendorCheckItems, dto.VendorCheckItemEmbed{
				ID:        v.ID.Hex(),
				Name:      v.Name,
				Location:  v.Location,
				CheckedAt: 0,
				CheckedBy: "",
				IsChecked: false,
				IsBlur:    isBlur,
				IsOffline: isOffline,
			})
		}
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
		IsVirtualCheck:   isVirtualCheck,
		IsFinish:         false,
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

//
//func (c *vendorCheckService) EditVendorCheck(user mjwt.CustomClaim, vendorCheckID string, input dto.VendorCheckEditRequest) (*dto.VendorCheck, rest_err.APIError) {
//	oid, errT := primitive.ObjectIDFromHex(vendorCheckID)
//	if errT != nil {
//		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
//	}
//
//	// Filling data
//	timeNow := time.Now().Unix()
//	data := dto.VendorCheckEdit{
//		FilterIDBranchAuthor: dto.FilterIDBranchAuthor{
//			FilterID:       oid,
//			FilterBranch:   user.Branch,
//			FilterAuthorID: user.Identity,
//		},
//		UpdatedAt:   timeNow,
//		UpdatedBy:   user.Name,
//		UpdatedByID: user.Identity,
//		IsFinish:    input.IsFinish,
//		Note:        input.Note,
//	}
//
//	// DB
//	vendorCheckEdited, err := c.daoC.EditVendorCheck(data)
//	if err != nil {
//		return nil, err
//	}
//
//	// IMPROVEMENT : make looping insert history to use goroutine
//	// Jika isFinish true , maka masukkan semua vendorCheckItem yang bertipe cctv
//	// looping ke insert history
//	var errorList []rest_err.APIError
//	if vendorCheckEdited.IsFinish {
//		for _, vendorCheckItem := range vendorCheckEdited.VendorCheckItems {
//			if vendorCheckItem.Type == category.Cctv {
//				// jika vendorCheckItemnya tidak di vendorCheck lewati
//				if !vendorCheckItem.IsVendorChecked {
//					continue
//				}
//				// cek complete status tidak boleh 0 atau 3, set default ke 1
//				if !(vendorCheckItem.CompleteStatus == enum.HComplete) {
//					vendorCheckItem.CompleteStatus = enum.HProgress
//				}
//				dataHistory := dto.HistoryRequest{
//					ParentID:       vendorCheckItem.ID,
//					Status:         "VendorChecklist",
//					Problem:        vendorCheckItem.VendorCheckedNote,
//					ProblemResolve: "",
//					CompleteStatus: vendorCheckItem.CompleteStatus,
//					DateStart:      timeNow,
//					DateEnd:        timeNow,
//					Tag:            []string{},
//				}
//				_, err := c.servHistory.InsertHistory(user, dataHistory)
//				if err != nil {
//					errorList = append(errorList, err)
//				}
//			}
//		}
//	}
//
//	// mengkoleksi semua error hasil looping insert history
//	errMessage := ""
//	if len(errorList) != 0 {
//		for _, err := range errorList {
//			errMessage = errMessage + ". " + err.Message()
//		}
//		logger.Error(fmt.Sprintf("VendorCheck berhasil diubah namun menambahkan history cctv gagal (EditVendorCheck isFinish) : %s", errMessage), errors.New("internal error"))
//		return nil, rest_err.NewInternalServerError("VendorCheck berhasil diubah namun menambahkan history cctv gagal", errors.New("internal error"))
//	}
//
//	return vendorCheckEdited, nil
//}
//
//func (c *vendorCheckService) DeleteVendorCheck(user mjwt.CustomClaim, id string) rest_err.APIError {
//	oid, errT := primitive.ObjectIDFromHex(id)
//	if errT != nil {
//		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
//	}
//
//	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
//	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
//	// DB
//	_, err := c.daoC.DeleteVendorCheck(dto.FilterIDBranchCreateGte{
//		FilterID:        oid,
//		FilterBranch:    user.Branch,
//		FilterCreateGTE: timeMinusOneDay.Unix(),
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// PutImage memasukkan lokasi file (path) ke dalam database vendorCheck dengan mengecek kesesuaian branch
//func (c *vendorCheckService) PutChildImage(user mjwt.CustomClaim, parentID string, childID string, imagePath string) (*dto.VendorCheck, rest_err.APIError) {
//	parentOid, errT := primitive.ObjectIDFromHex(parentID)
//	if errT != nil {
//		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
//	}
//
//	filter := dto.FilterParentIDChildIDAuthor{
//		FilterParentID: parentOid,
//		FilterChildID:  childID,
//		FilterAuthorID: user.Identity,
//	}
//
//	vendorCheck, err := c.daoC.UploadChildImage(filter, imagePath)
//	if err != nil {
//		return nil, err
//	}
//	return vendorCheck, nil
//}
//
//func (c *vendorCheckService) UpdateVendorCheckItem(user mjwt.CustomClaim, input dto.VendorCheckChildUpdateRequest) (*dto.VendorCheck, rest_err.APIError) {
//	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
//	if errT != nil {
//		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
//	}
//
//	timeNow := time.Now().Unix()
//
//	// DB
//	data := dto.VendorCheckChildUpdate{
//		FilterParentIDChildIDAuthor: dto.FilterParentIDChildIDAuthor{
//			FilterParentID: parentOid,
//			FilterChildID:  input.ChildID,
//			FilterAuthorID: user.Identity,
//		},
//		UpdatedAt:         timeNow,
//		VendorCheckedAt:   timeNow,
//		IsVendorChecked:   input.IsVendorChecked,
//		TagSelected:       input.TagSelected,
//		TagExtraSelected:  input.TagExtraSelected,
//		VendorCheckedNote: input.VendorCheckedNote,
//		HaveProblem:       input.HaveProblem,
//		CompleteStatus:    input.CompleteStatus,
//	}
//	vendorCheck, err := c.daoC.UpdateVendorCheckItem(data)
//	if err != nil {
//		return nil, err
//	}
//
//	// Cek index dan Type childID yang diupdate
//	var updatedType string
//	// var indexItems int //index digunakan untuk memepermudah mendapatkan nama VendorCheckitem yang diupdate
//	for _, v := range vendorCheck.VendorCheckItems {
//		if v.ID == input.ChildID {
//			updatedType = v.Type
//			// indexItems = i
//		}
//	}
//
//	// Mengupdate value di vendorCheck item, agar pada pembuatan vendorCheck berikutnya pesan tetap berlanjut
//	// Kecuali cctv yang mana tidak memiliki vendorCheck item
//	childOid, errT := primitive.ObjectIDFromHex(input.ChildID)
//	if errT != nil {
//		return nil, rest_err.NewBadRequestError("Child ObjectID yang dimasukkan salah")
//	}
//
//	if updatedType != category.Cctv {
//		_, err = c.daoCI.EditVendorCheckItemValue(
//			dto.VendorCheckItemEditBySys{
//				FilterID:          childOid,
//				UpdatedAt:         0,
//				VendorCheckedNote: input.VendorCheckedNote,
//				HaveProblem:       input.HaveProblem,
//				CompleteStatus:    input.CompleteStatus,
//			})
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return vendorCheck, nil
//}
//
//func (c *vendorCheckService) GetVendorCheckByID(vendorCheckID string, branchIfSpecific string) (*dto.VendorCheck, rest_err.APIError) {
//	oid, errT := primitive.ObjectIDFromHex(vendorCheckID)
//	if errT != nil {
//		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
//	}
//
//	vendorCheck, err := c.daoC.GetVendorCheckByID(oid, branchIfSpecific)
//	if err != nil {
//		return nil, err
//	}
//	return vendorCheck, nil
//}
//
//func (c *vendorCheckService) FindVendorCheck(branch string, filter dto.FilterTimeRangeLimit) (dto.VendorCheckResponseMinList, rest_err.APIError) {
//	vendorCheckList, err := c.daoC.FindVendorCheck(branch, filter)
//	if err != nil {
//		return nil, err
//	}
//	return vendorCheckList, nil
//}