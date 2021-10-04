package service

import (
	"context"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/dao/altaiphycheckdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/otherdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func NewAltaiPhyCheckService(
	altaiPhyCheckDao altaiphycheckdao.CheckAltaiPhyDaoAssumer,
	genUnitDao genunitdao.GenUnitDaoAssumer,
	otherDao otherdao.OtherLoader,
	histService HistoryServiceAssumer,
) AltaiPhyCheckServiceAssumer {
	return &altaiPhyCheckService{
		daoC:        altaiPhyCheckDao,
		daoG:        genUnitDao,
		daoAltai:    otherDao,
		servHistory: histService,
	}
}

type altaiPhyCheckService struct {
	daoC        altaiphycheckdao.CheckAltaiPhyDaoAssumer
	daoG        genunitdao.GenUnitDaoAssumer
	daoAltai    otherdao.OtherLoader
	servHistory HistoryServiceAssumer
}
type AltaiPhyCheckServiceAssumer interface {
	InsertAltaiPhyCheck(ctx context.Context, user mjwt.CustomClaim, name string, isQuarterMode bool) (*string, rest_err.APIError)
	DeleteAltaiPhyCheck(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	GetAltaiPhyCheckByID(ctx context.Context, altaiCheckID string, branchIfSpecific string) (*dto.AltaiPhyCheck, rest_err.APIError)
	FindAltaiPhyCheck(ctx context.Context, branch string, filter dto.FilterTimeRangeLimit) ([]dto.AltaiPhyCheck, rest_err.APIError)
	UpdateAltaiPhyCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.AltaiPhyCheckItemUpdateRequest) (*dto.AltaiPhyCheck, rest_err.APIError)
	BulkUpdateAltaiPhyItem(ctx context.Context, user mjwt.CustomClaim, inputs []dto.AltaiPhyCheckItemUpdateRequest) (string, rest_err.APIError)
	FinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.AltaiPhyCheck, rest_err.APIError)
}

func (vc *altaiPhyCheckService) InsertAltaiPhyCheck(ctx context.Context, user mjwt.CustomClaim, name string, isQuarterMode bool) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	// ambil altai genUnit item berdasarkan cabang yang di input
	// mendapatkan data cases
	genItems, err := vc.daoG.FindUnit(dto.GenUnitFilter{
		Branch:   user.Branch,
		Category: category.Altai,
		Pings:    false,
	})
	if err != nil {
		return nil, err
	}

	// ambil altai untuk mendapatkan data lokasi
	// altaiItems sudah sorted berdasarkan lokasi sedangkan genItems tidak
	altaiItems, err := vc.daoAltai.FindOther(ctx, dto.FilterOther{
		FilterBranch:      user.Branch,
		FilterSubCategory: category.Altai,
	})
	if err != nil {
		return nil, err
	}

	// ubah altai genUnit menjadi map
	genItemsMap := make(map[string]dto.GenUnitResponse)
	for _, item := range genItems {
		genItemsMap[item.ID] = item
	}

	// kategorikan altaiCheck item menggunakan perulangan altaiItems
	// karena altaiItems sudah sorted
	var altaiCheckItems []dto.AltaiPhyCheckItemEmbed
	for _, v := range altaiItems {
		altaiInfoFromGenUnit := genItemsMap[v.ID.Hex()]

		// jika didalam semua case yang ada di altai tersebut ada tag #isBlur maka kita anggap altainya blur
		// dan tidak mati
		isOffline := altaiInfoFromGenUnit.CasesSize != 0
		isBlur := strings.Contains(fmt.Sprintf("%v", altaiInfoFromGenUnit.Cases), "#isBlur")
		if isBlur {
			isOffline = false
		}

		altaiCheckItems = append(altaiCheckItems, dto.AltaiPhyCheckItemEmbed{
			ID:           v.ID.Hex(),
			Name:         v.Name,
			Location:     v.Location,
			CheckedAt:    0,
			CheckedBy:    "",
			IsChecked:    false,
			IsMaintained: false,
			IsOffline:    isOffline,
		})
	}

	data := dto.AltaiPhyCheck{
		QuarterlyMode:      isQuarterMode,
		Name:               name,
		CreatedAt:          timeNow,
		CreatedBy:          user.Name,
		CreatedByID:        user.Identity,
		UpdatedAt:          timeNow,
		UpdatedBy:          user.Name,
		UpdatedByID:        user.Identity,
		Branch:             user.Branch,
		TimeStarted:        timeNow,
		TimeEnded:          0,
		IsFinish:           false,
		AltaiPhyCheckItems: altaiCheckItems,
		Note:               "",
	}

	// DB
	insertedID, err := vc.daoC.InsertCheck(ctx, data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (vc *altaiPhyCheckService) DeleteAltaiPhyCheck(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
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

// UpdateAltaiPhyCheckItem
// setiap melakukan update akan mengupdate cek fisik lainnya yang masih belum finish
func (vc *altaiPhyCheckService) UpdateAltaiPhyCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.AltaiPhyCheckItemUpdateRequest) (*dto.AltaiPhyCheck, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// DB
	data := dto.AltaiPhyCheckItemUpdate{
		FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
			FilterParentID: parentOid,
			FilterChildID:  input.ChildID,
			FilterBranch:   user.Branch,
		},
		CheckedAt:    timeNow,
		CheckedBy:    user.Name,
		IsChecked:    input.IsChecked,
		IsMaintained: input.IsMaintained,
		IsOffline:    input.IsOffline,
	}
	// harus sukses mengupdate dirinya sendiri dulu karena ada validasi, baru update
	// check fisik lain yang belum finish
	altaiCheck, err := vc.daoC.UpdateCheckItem(ctx, data)
	if err != nil {
		return nil, err
	}

	go func() {
		checkStillOpens, err := vc.daoC.FindCheckStillOpen(ctx, user.Branch, false)
		if err != nil {
			return
		}
		var checkStillOpensValid []dto.AltaiPhyCheckItemUpdate
		for _, check := range checkStillOpens {
			// jika id nya sama lewati, karena sudah diubah pada kesempatan pertama diatas
			if check.ID.Hex() == input.ParentID {
				continue
			}
			checkStillOpensValid = append(checkStillOpensValid, dto.AltaiPhyCheckItemUpdate{
				FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
					FilterParentID: check.ID,
					FilterChildID:  input.ChildID,
					FilterBranch:   user.Branch,
				},
				CheckedAt:    timeNow,
				CheckedBy:    user.Name,
				IsChecked:    input.IsChecked,
				IsMaintained: input.IsMaintained,
				IsOffline:    input.IsOffline,
			})
		}

		if len(checkStillOpensValid) == 0 {
			return
		}
		updatedCount, err := vc.daoC.BulkUpdateItem(ctx, checkStillOpensValid)
		if err != nil {
			logger.Error("gagal bulk update pada (UpdateAltaiPhyCheckItem)", err)
		}
		logger.Info(fmt.Sprintf("berhasil bulk update (UpdateAltaiPhyCheckItem) dengan %d perubahan", updatedCount))
	}()

	return altaiCheck, nil
}

func (vc *altaiPhyCheckService) BulkUpdateAltaiPhyItem(ctx context.Context, user mjwt.CustomClaim, inputs []dto.AltaiPhyCheckItemUpdateRequest) (string, rest_err.APIError) {
	if len(inputs) == 0 {
		return "", rest_err.NewBadRequestError("tidak ada perubahan")
	}

	parentOid, errT := primitive.ObjectIDFromHex(inputs[0].ParentID)
	if errT != nil {
		return "", rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	inputDatas := make([]dto.AltaiPhyCheckItemUpdate, len(inputs))
	for i, input := range inputs {
		inputDatas[i] = dto.AltaiPhyCheckItemUpdate{
			FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
				FilterParentID: parentOid,
				FilterChildID:  input.ChildID,
				FilterBranch:   user.Branch,
			},
			CheckedAt:    timeNow,
			CheckedBy:    user.Name,
			IsChecked:    input.IsChecked,
			IsMaintained: input.IsMaintained,
			IsOffline:    input.IsOffline,
		}
	}
	result, err := vc.daoC.BulkUpdateItem(ctx, inputDatas)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d data telah diubah", result), nil
}

func (vc *altaiPhyCheckService) FinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.AltaiPhyCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(detailID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// tandai isFinish true dan end_date ke waktu sekarang
	altaiChecklistDetail, err := vc.daoC.EditCheck(ctx, dto.AltaiPhyCheckEdit{
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

	return altaiChecklistDetail, nil
}

func (vc *altaiPhyCheckService) GetAltaiPhyCheckByID(ctx context.Context, altaiCheckID string, branchIfSpecific string) (*dto.AltaiPhyCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(altaiCheckID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	altaiCheck, err := vc.daoC.GetCheckByID(ctx, oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return altaiCheck, nil
}

func (vc *altaiPhyCheckService) FindAltaiPhyCheck(ctx context.Context, branch string, filter dto.FilterTimeRangeLimit) ([]dto.AltaiPhyCheck, rest_err.APIError) {
	altaiCheckList, err := vc.daoC.FindCheck(ctx, branch, filter, false)
	if err != nil {
		return []dto.AltaiPhyCheck{}, err
	}
	return altaiCheckList, nil
}
