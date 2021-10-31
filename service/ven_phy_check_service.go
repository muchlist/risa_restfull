package service

import (
	"context"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/dao/cctvdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/venphycheckdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func NewVenPhyCheckService(
	venPhyCheckDao venphycheckdao.CheckVenPhyDaoAssumer,
	genUnitDao genunitdao.GenUnitLoader,
	cctvDao cctvdao.CctvDaoAssumer,
	histService HistoryServiceAssumer,
) VenPhyCheckServiceAssumer {
	return &venPhyCheckService{
		daoC:        venPhyCheckDao,
		daoG:        genUnitDao,
		daoCTV:      cctvDao,
		servHistory: histService,
	}
}

type venPhyCheckService struct {
	daoC        venphycheckdao.CheckVenPhyDaoAssumer
	daoG        genunitdao.GenUnitLoader
	daoCTV      cctvdao.CctvLoader
	servHistory HistoryServiceAssumer
}
type VenPhyCheckServiceAssumer interface {
	InsertVenPhyCheck(ctx context.Context, user mjwt.CustomClaim, name string, isQuarterMode bool) (*string, rest_err.APIError)
	DeleteVenPhyCheck(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	GetVenPhyCheckByID(ctx context.Context, vendorCheckID string, branchIfSpecific string) (*dto.VenPhyCheck, rest_err.APIError)
	FindVenPhyCheck(ctx context.Context, branch string, filter dto.FilterTimeRangeLimit) ([]dto.VenPhyCheck, rest_err.APIError)
	UpdateVenPhyCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.VenPhyCheckItemUpdateRequest) (*dto.VenPhyCheck, rest_err.APIError)
	BulkUpdateVenPhyItem(ctx context.Context, user mjwt.CustomClaim, inputs []dto.VenPhyCheckItemUpdateRequest) (string, rest_err.APIError)
	FinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.VenPhyCheck, rest_err.APIError)
	UndoFinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.VenPhyCheck, rest_err.APIError)
	FreshUpdateNameCCTV(ctx context.Context, branch string) (string, rest_err.APIError)
}

func (vc *venPhyCheckService) InsertVenPhyCheck(ctx context.Context, user mjwt.CustomClaim, name string, isQuarterMode bool) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	// ambil cctv genUnit item berdasarkan cabang yang di input
	// mendapatkan data cases
	genItems, err := vc.daoG.FindUnit(ctx, dto.GenUnitFilter{
		Branch:   user.Branch,
		Category: category.Cctv,
		Pings:    false,
	})
	if err != nil {
		return nil, err
	}

	// ambil cctv untuk mendapatkan data lokasi
	// cctvItems sudah sorted berdasarkan lokasi sedangkan genItems tidak
	cctvItems, err := vc.daoCTV.FindCctv(ctx, dto.FilterBranchLocIPNameDisable{
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
	var vendorCheckItems []dto.VenPhyCheckItemEmbed
	for _, v := range cctvItems {
		cctvInfoFromGenUnit := genItemsMap[v.ID.Hex()]

		// jika didalam semua case yang ada di cctv tersebut ada tag #isBlur maka kita anggap cctvnya blur
		// dan tidak mati
		isOffline := cctvInfoFromGenUnit.CasesSize != 0
		isBlur := strings.Contains(fmt.Sprintf("%v", cctvInfoFromGenUnit.Cases), "#isBlur")
		if isBlur {
			isOffline = false
		}

		vendorCheckItems = append(vendorCheckItems, dto.VenPhyCheckItemEmbed{
			ID:           v.ID.Hex(),
			Name:         v.Name,
			Location:     v.Location,
			CheckedAt:    0,
			CheckedBy:    "",
			IsChecked:    false,
			IsMaintained: false,
			IsBlur:       isBlur,
			IsOffline:    isOffline,
			DisVendor:    v.DisVendor,
		})
	}

	data := dto.VenPhyCheck{
		QuarterlyMode:    isQuarterMode,
		Name:             name,
		CreatedAt:        timeNow,
		CreatedBy:        user.Name,
		CreatedByID:      user.Identity,
		UpdatedAt:        timeNow,
		UpdatedBy:        user.Name,
		UpdatedByID:      user.Identity,
		Branch:           user.Branch,
		TimeStarted:      timeNow,
		TimeEnded:        0,
		IsFinish:         false,
		VenPhyCheckItems: vendorCheckItems,
		Note:             "",
	}

	// DB
	insertedID, err := vc.daoC.InsertCheck(ctx, data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (vc *venPhyCheckService) DeleteVenPhyCheck(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := vc.daoC.DeleteCheck(ctx, dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateVenPhyCheckItem
// setiap melakukan update akan mengupdate cek fisik lainnya yang masih belum finish
func (vc *venPhyCheckService) UpdateVenPhyCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.VenPhyCheckItemUpdateRequest) (*dto.VenPhyCheck, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// DB
	data := dto.VenPhyCheckItemUpdate{
		FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
			FilterParentID: parentOid,
			FilterChildID:  input.ChildID,
			FilterBranch:   user.Branch,
		},
		CheckedAt:    timeNow,
		CheckedBy:    user.Name,
		IsChecked:    input.IsChecked,
		IsMaintained: input.IsMaintained,
		IsBlur:       input.IsBlur,
		IsOffline:    input.IsOffline,
	}
	// harus sukses mengupdate dirinya sendiri dulu karena ada validasi, baru update
	// check fisik lain yang belum finish
	vendorCheck, err := vc.daoC.UpdateCheckItem(ctx, data)
	if err != nil {
		return nil, err
	}

	go func() {
		checkStillOpens, err := vc.daoC.FindCheckStillOpen(ctx, user.Branch, false)
		if err != nil {
			return
		}
		var checkStillOpensValid []dto.VenPhyCheckItemUpdate
		for _, check := range checkStillOpens {
			// jika id nya sama lewati, karena sudah diubah pada kesempatan pertama diatas
			if check.ID.Hex() == input.ParentID {
				continue
			}
			checkStillOpensValid = append(checkStillOpensValid, dto.VenPhyCheckItemUpdate{
				FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
					FilterParentID: check.ID,
					FilterChildID:  input.ChildID,
					FilterBranch:   user.Branch,
				},
				CheckedAt:    timeNow,
				CheckedBy:    user.Name,
				IsChecked:    input.IsChecked,
				IsMaintained: input.IsMaintained,
				IsBlur:       input.IsBlur,
				IsOffline:    input.IsOffline,
			})
		}

		if len(checkStillOpensValid) == 0 {
			return
		}
		updatedCount, err := vc.daoC.BulkUpdateItemForUpdateCheckItem(ctx, checkStillOpensValid)
		if err != nil {
			logger.Error("gagal bulk update pada (UpdateVenPhyCheckItem)", err)
		}
		logger.Info(fmt.Sprintf("berhasil bulk update (UpdateVenPhyCheckItem) dengan %d perubahan", updatedCount))
	}()

	return vendorCheck, nil
}

func (vc *venPhyCheckService) BulkUpdateVenPhyItem(ctx context.Context, user mjwt.CustomClaim, inputs []dto.VenPhyCheckItemUpdateRequest) (string, rest_err.APIError) {
	if len(inputs) == 0 {
		return "", rest_err.NewBadRequestError("tidak ada perubahan")
	}

	parentOid, errT := primitive.ObjectIDFromHex(inputs[0].ParentID)
	if errT != nil {
		return "", rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	inputDatas := make([]dto.VenPhyCheckItemUpdate, len(inputs))
	for i, input := range inputs {
		timeInput := time.Now().Unix()
		if input.CheckedAt != 0 {
			timeInput = input.CheckedAt
		}

		inputDatas[i] = dto.VenPhyCheckItemUpdate{
			FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
				FilterParentID: parentOid,
				FilterChildID:  input.ChildID,
				FilterBranch:   user.Branch,
			},
			CheckedAt:    timeInput,
			CheckedBy:    user.Name,
			IsChecked:    input.IsChecked,
			IsMaintained: input.IsMaintained,
			IsBlur:       input.IsBlur,
			IsOffline:    input.IsOffline,
		}
	}
	result, err := vc.daoC.BulkUpdateItem(ctx, inputDatas)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d data telah diubah", result), nil
}

func (vc *venPhyCheckService) FinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.VenPhyCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(detailID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// tandai isFinish true dan end_date ke waktu sekarang
	cctvChecklistDetail, err := vc.daoC.EditCheck(ctx, dto.VenPhyCheckEdit{
		FilterIDBranch: dto.FilterIDBranch{
			FilterID:     oid,
			FilterBranch: user.Branch,
		},
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		TimeEnded:   timeNow,
		IsFinish:    true,
		Note:        "",
	})
	if err != nil {
		return nil, err
	}

	return cctvChecklistDetail, nil
}

func (vc *venPhyCheckService) UndoFinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.VenPhyCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(detailID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	cctvChecklistDetail, err := vc.daoC.UndoFinishCheck(ctx, oid, user.Branch)
	if err != nil {
		return nil, err
	}

	return cctvChecklistDetail, nil
}

func (vc *venPhyCheckService) FreshUpdateNameCCTV(ctx context.Context, branch string) (string, rest_err.APIError) {
	if branch == "" {
		return "", rest_err.NewBadRequestError("cabang harus di isi")
	}

	// mendapatkan list cek fisik cctv yang masih terbuka
	vendorCheckList, err := vc.daoC.FindCheckStillOpen(ctx, strings.ToUpper(branch), true)
	if err != nil {
		return "", err
	}

	// mendapatkan nama cctv eksisting
	cctvList, err := vc.daoCTV.FindCctv(ctx, dto.FilterBranchLocIPNameDisable{
		FilterBranch:  strings.ToUpper(branch),
		FilterDisable: false,
	})
	if err != nil {
		return "", err
	}

	// buat daftar nama cctv jadi map
	cctvNameMap := make(map[string]string, len(cctvList))
	cctvDisVendorMap := make(map[string]bool, len(cctvList))
	for _, cctv := range cctvList {
		cctvNameMap[cctv.ID.Hex()] = cctv.Name
		cctvDisVendorMap[cctv.ID.Hex()] = cctv.DisVendor
	}

	totalChange := 0

	// perulangan untuk setiap cek fisik cctv
	for _, v := range vendorCheckList {
		hasChanged := 0
		// replace nama cctv untuk setiap vendorchecklist dengan yang baru
		checklist := v.VenPhyCheckItems
		for i, c := range checklist {
			newCCTVName, exist := cctvNameMap[c.ID]
			if !exist {
				checklist[i].Name = "Deleted CCTV"
				continue
			} else {
				if checklist[i].Name != newCCTVName {
					checklist[i].Name = newCCTVName
					hasChanged++
				}
				checklist[i].DisVendor = cctvDisVendorMap[c.ID]
			}
		}

		// override vendorchecklist jika hasChanged != 0
		if hasChanged != 0 {
			totalChange += hasChanged
			_, err := vc.daoC.OverwriteChecklist(ctx, v.ID, checklist)
			if err != nil {
				// stop proses jika ada error
				return "", err
			}
		}
	}

	return fmt.Sprintf("%d item-check has been affected", totalChange), nil
}

func (vc *venPhyCheckService) GetVenPhyCheckByID(ctx context.Context, vendorCheckID string, branchIfSpecific string) (*dto.VenPhyCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(vendorCheckID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	vendorCheck, err := vc.daoC.GetCheckByID(ctx, oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return vendorCheck, nil
}

func (vc *venPhyCheckService) FindVenPhyCheck(ctx context.Context, branch string, filter dto.FilterTimeRangeLimit) ([]dto.VenPhyCheck, rest_err.APIError) {
	vendorCheckList, err := vc.daoC.FindCheck(ctx, branch, filter, false)
	if err != nil {
		return []dto.VenPhyCheck{}, err
	}
	return vendorCheckList, nil
}
