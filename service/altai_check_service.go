package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dao/altaicheckdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/otherdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewAltaiCheckService(
	altaiCheckDao altaicheckdao.CheckAltaiDaoAssumer,
	genUnitDao genunitdao.GenUnitDaoAssumer,
	altaiDao otherdao.OtherDaoAssumer,
	histService HistoryServiceAssumer,
) AltaiCheckServiceAssumer {
	return &altaiCheckService{
		daoC:        altaiCheckDao,
		daoG:        genUnitDao,
		daoAltai:    altaiDao,
		servHistory: histService,
	}
}

type altaiCheckService struct {
	daoC        altaicheckdao.CheckAltaiDaoAssumer
	daoG        genunitdao.GenUnitDaoAssumer
	daoAltai    otherdao.OtherDaoAssumer
	servHistory HistoryServiceAssumer
}
type AltaiCheckServiceAssumer interface {
	InsertAltaiCheck(user mjwt.CustomClaim) (*string, rest_err.APIError)
	DeleteAltaiCheck(user mjwt.CustomClaim, id string) rest_err.APIError
	GetAltaiCheckByID(altaiCheckID string, branchIfSpecific string) (*dto.AltaiCheck, rest_err.APIError)
	FindAltaiCheck(branch string, filter dto.FilterTimeRangeLimit) ([]dto.AltaiCheck, rest_err.APIError)
	UpdateAltaiCheckItem(user mjwt.CustomClaim, input dto.AltaiCheckItemUpdateRequest) (*dto.AltaiCheck, rest_err.APIError)
	BulkUpdateAltaiItem(user mjwt.CustomClaim, inputs []dto.AltaiCheckItemUpdateRequest) (string, rest_err.APIError)
	FinishCheck(user mjwt.CustomClaim, detailID string) (*dto.AltaiCheck, rest_err.APIError)
}

func (c *altaiCheckService) InsertAltaiCheck(user mjwt.CustomClaim) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	// ambil altai genUnit item berdasarkan cabang yang di input
	// mendapatkan data cases
	genItems, err := c.daoG.FindUnit(dto.GenUnitFilter{
		Branch:   user.Branch,
		Category: category.Altai,
		Pings:    false,
	})
	if err != nil {
		return nil, err
	}

	// ambil altai untuk mendapatkan data lokasi
	// altaiItems sudah sorted berdasarkan lokasi sedangkan genItems tidak
	altaiItems, err := c.daoAltai.FindOther(dto.FilterOther{
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
	var altaiCheckItems []dto.AltaiCheckItemEmbed
	for _, v := range altaiItems {
		altaiInfoFromGenUnit := genItemsMap[v.ID.Hex()]

		isOffline := altaiInfoFromGenUnit.CasesSize != 0

		// pengecekan secara virtual default sudah tercek semua di waktu pembuatan
		altaiCheckItems = append(altaiCheckItems, dto.AltaiCheckItemEmbed{
			ID:        v.ID.Hex(),
			Name:      v.Name,
			Location:  v.Location,
			CheckedAt: timeNow,
			CheckedBy: user.Name,
			IsChecked: true,
			IsOffline: isOffline,
		})
	}

	data := dto.AltaiCheck{
		CreatedAt:       timeNow,
		CreatedBy:       user.Name,
		CreatedByID:     user.Identity,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		Branch:          user.Branch,
		TimeStarted:     timeNow,
		TimeEnded:       0,
		IsFinish:        false,
		AltaiCheckItems: altaiCheckItems,
		Note:            "",
	}

	// DB
	insertedID, err := c.daoC.InsertCheck(data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (c *altaiCheckService) DeleteAltaiCheck(user mjwt.CustomClaim, id string) rest_err.APIError {
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

func (c *altaiCheckService) UpdateAltaiCheckItem(user mjwt.CustomClaim, input dto.AltaiCheckItemUpdateRequest) (*dto.AltaiCheck, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// DB
	data := dto.AltaiCheckItemUpdate{
		FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
			FilterParentID: parentOid,
			FilterChildID:  input.ChildID,
			FilterBranch:   user.Branch,
		},
		CheckedAt: timeNow,
		CheckedBy: user.Name,
		IsChecked: input.IsChecked,
		IsOffline: input.IsOffline,
	}
	altaiCheck, err := c.daoC.UpdateCheckItem(data)
	if err != nil {
		return nil, err
	}
	return altaiCheck, nil
}

func (c *altaiCheckService) BulkUpdateAltaiItem(user mjwt.CustomClaim, inputs []dto.AltaiCheckItemUpdateRequest) (string, rest_err.APIError) {
	if len(inputs) == 0 {
		return "", rest_err.NewBadRequestError("tidak ada perubahan")
	}

	parentOid, errT := primitive.ObjectIDFromHex(inputs[0].ParentID)
	if errT != nil {
		return "", rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	inputDatas := make([]dto.AltaiCheckItemUpdate, len(inputs))
	for i, input := range inputs {
		inputDatas[i] = dto.AltaiCheckItemUpdate{
			FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
				FilterParentID: parentOid,
				FilterChildID:  input.ChildID,
				FilterBranch:   user.Branch,
			},
			CheckedAt: timeNow,
			CheckedBy: user.Name,
			IsChecked: input.IsChecked,
			IsOffline: input.IsOffline,
		}
	}
	result, err := c.daoC.BulkUpdateItem(inputDatas)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d data telah diubah", result), nil
}

func (c *altaiCheckService) FinishCheck(user mjwt.CustomClaim, detailID string) (*dto.AltaiCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(detailID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// 1. cek altai existing, untuk mendapatkan keterangan apakah ada case
	genItems, err := c.daoG.FindUnit(dto.GenUnitFilter{
		Branch:   user.Branch,
		Category: category.Altai,
		Pings:    false,
	})
	if err != nil {
		return nil, err
	}

	// 2. altai yang memiliki case tersebut di exclude dari update berikutnya
	var altaiExcludeID []string
	for _, altai := range genItems {
		if altai.CasesSize != 0 {
			altaiExcludeID = append(altaiExcludeID, altai.ID)
		}
	}

	// 3. filter items check yang memiliki is_offline true
	var altaiOfflineID []string
	altaiChecklistDetail, err := c.daoC.GetCheckByID(oid, user.Branch)
	if err != nil {
		return nil, err
	}
	for _, altai := range altaiChecklistDetail.AltaiCheckItems {
		if sfunc.InSlice(altai.ID, altaiExcludeID) {
			continue
		}
		if altai.IsOffline {
			altaiOfflineID = append(altaiOfflineID, altai.ID)
		}
	}

	// send to background
	go func() {
		// Insert History isoffline
		if len(altaiOfflineID) != 0 {
			for _, altaiID := range altaiOfflineID {
				_, _ = c.servHistory.InsertHistory(user, dto.HistoryRequest{
					ParentID:       altaiID,
					Status:         "",
					Problem:        "ALTAI Offline",
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
	altaiChecklistDetail, err = c.daoC.EditCheck(dto.AltaiCheckEdit{
		FilterIDBranch: dto.FilterIDBranch{
			FilterID:     oid,
			FilterBranch: user.Branch,
		},
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		TimeStarted: altaiChecklistDetail.TimeStarted,
		TimeEnded:   timeNow,
		IsFinish:    true,
		Note:        "",
	})
	if err != nil {
		return nil, err
	}

	return altaiChecklistDetail, nil
}

func (c *altaiCheckService) GetAltaiCheckByID(altaiCheckID string, branchIfSpecific string) (*dto.AltaiCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(altaiCheckID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	altaiCheck, err := c.daoC.GetCheckByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return altaiCheck, nil
}

func (c *altaiCheckService) FindAltaiCheck(branch string, filter dto.FilterTimeRangeLimit) ([]dto.AltaiCheck, rest_err.APIError) {
	altaiCheckList, err := c.daoC.FindCheck(branch, filter, false)
	if err != nil {
		return []dto.AltaiCheck{}, err
	}
	return altaiCheckList, nil
}
