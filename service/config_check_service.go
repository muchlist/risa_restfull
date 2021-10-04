package service

import (
	"context"
	"fmt"
	"github.com/muchlist/risa_restfull/constants/enum"
	"sort"
	"time"

	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/dao/configcheckdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/otherdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewConfigCheckService(
	configCheckDao configcheckdao.CheckConfigDaoAssumer,
	genUnitDao genunitdao.GenUnitDaoAssumer,
	configDao otherdao.OtherDaoAssumer,
	histService HistoryServiceAssumer,
) ConfigCheckServiceAssumer {
	return &configCheckService{
		daoC:        configCheckDao,
		daoG:        genUnitDao,
		daoNetwork:  configDao,
		servHistory: histService,
	}
}

type configCheckService struct {
	daoC        configcheckdao.CheckConfigDaoAssumer
	daoG        genunitdao.GenUnitDaoAssumer
	daoNetwork  otherdao.OtherDaoAssumer
	servHistory HistoryServiceAssumer
}
type ConfigCheckServiceAssumer interface {
	InsertConfigCheck(ctx context.Context, user mjwt.CustomClaim) (*string, rest_err.APIError)
	DeleteConfigCheck(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	GetConfigCheckByID(ctx context.Context, configCheckID string, branchIfSpecific string) (*dto.ConfigCheck, rest_err.APIError)
	UpdateManyConfigCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.ConfigCheckUpdateManyRequest) (*dto.ConfigCheck, rest_err.APIError)
	FindConfigCheck(ctx context.Context, branch string, filter dto.FilterTimeRangeLimit) ([]dto.ConfigCheck, rest_err.APIError)
	UpdateConfigCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.ConfigCheckItemUpdateRequest) (*dto.ConfigCheck, rest_err.APIError)
	FinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.ConfigCheck, rest_err.APIError)
}

func (c *configCheckService) InsertConfigCheck(ctx context.Context, user mjwt.CustomClaim) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	networkItems, err := c.daoNetwork.FindOther(ctx, dto.FilterOther{
		FilterBranch:      user.Branch,
		FilterSubCategory: fmt.Sprintf("%s,%s", category.Network, category.Altai),
	})
	if err != nil {
		return nil, err
	}

	// translit networkItem menjadi checkItem
	configCheckItem := make([]dto.ConfigCheckItemEmbed, len(networkItems))
	for i, networkItem := range networkItems {
		configCheckItem[i] = dto.ConfigCheckItemEmbed{
			ID:       networkItem.ID.Hex(),
			Name:     networkItem.Name,
			Location: networkItem.Location,
		}
	}

	// sorting berdasarkan nama
	sort.Slice(configCheckItem, func(i, j int) bool {
		return configCheckItem[i].Name < configCheckItem[j].Name
	})

	data := dto.ConfigCheck{
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
		ConfigCheckItems: configCheckItem,
		Note:             "",
	}

	// DB
	insertedID, err := c.daoC.InsertCheck(ctx, data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (c *configCheckService) DeleteConfigCheck(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := c.daoC.DeleteCheck(ctx, dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *configCheckService) UpdateConfigCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.ConfigCheckItemUpdateRequest) (*dto.ConfigCheck, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// DB
	data := dto.ConfigCheckItemUpdate{
		FilterParentIDChildIDBranch: dto.FilterParentIDChildIDBranch{
			FilterParentID: parentOid,
			FilterChildID:  input.ChildID,
			FilterBranch:   user.Branch,
		},
		CheckedAt: timeNow,
		CheckedBy: user.Name,
		IsUpdated: input.IsUpdated,
	}
	configCheck, err := c.daoC.UpdateCheckItem(ctx, data)
	if err != nil {
		return nil, err
	}
	return configCheck, nil
}

func (c *configCheckService) UpdateManyConfigCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.ConfigCheckUpdateManyRequest) (*dto.ConfigCheck, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	// data fill updatedCheck
	data := dto.ConfigCheckUpdateMany{
		ParentID:       parentOid,
		ChildIDsUpdate: input.ChildUpdate,
		UpdatedValue:   true,
		Branch:         user.Branch,
		Updater:        user.Name,
	}

	err := c.daoC.UpdateManyItem(ctx, data)
	if err != nil {
		return nil, err
	}

	// data fill updatedCheck
	data = dto.ConfigCheckUpdateMany{
		ParentID:       parentOid,
		ChildIDsUpdate: input.ChildNotUpdate,
		UpdatedValue:   false,
		Branch:         user.Branch,
		Updater:        user.Name,
	}

	err = c.daoC.UpdateManyItem(ctx, data)
	if err != nil {
		return nil, err
	}

	configCheck, err := c.daoC.GetCheckByID(ctx, parentOid, "")
	if err != nil {
		return nil, err
	}
	return configCheck, nil
}

func (c *configCheckService) FinishCheck(ctx context.Context, user mjwt.CustomClaim, detailID string) (*dto.ConfigCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(detailID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// find config item yang sudah diupdate
	var configUpdatedIDs []string
	configChecklistDetail, err := c.daoC.GetCheckByID(ctx, oid, user.Branch)
	if err != nil {
		return nil, err
	}

	for _, configItem := range configChecklistDetail.ConfigCheckItems {
		if configItem.IsUpdated {
			configUpdatedIDs = append(configUpdatedIDs, configItem.ID)
		}
	}

	// send to background
	go func() {
		if len(configUpdatedIDs) != 0 {
			for _, configID := range configUpdatedIDs {
				_, _ = c.servHistory.InsertHistory(user, dto.HistoryRequest{
					ParentID:       configID,
					Status:         "",
					Problem:        "Pengecekan auto backup",
					ProblemResolve: "update terkonfirmasi",
					CompleteStatus: enum.HDataInfo,
					DateStart:      timeNow,
					DateEnd:        0,
					Tag:            []string{},
				})
			}
		}
	}()

	// 7. tandai isFinish true dan end_date ke waktu sekarang
	configChecklistDetail, err = c.daoC.EditCheck(ctx, dto.ConfigCheckEdit{
		FilterIDBranch: dto.FilterIDBranch{
			FilterID:     oid,
			FilterBranch: user.Branch,
		},
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		TimeStarted: configChecklistDetail.TimeStarted,
		TimeEnded:   timeNow,
		IsFinish:    true,
		Note:        "",
	})
	if err != nil {
		return nil, err
	}

	return configChecklistDetail, nil
}

func (c *configCheckService) GetConfigCheckByID(ctx context.Context, configCheckID string, branchIfSpecific string) (*dto.ConfigCheck, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(configCheckID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	configCheck, err := c.daoC.GetCheckByID(ctx, oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return configCheck, nil
}

func (c *configCheckService) FindConfigCheck(ctx context.Context, branch string, filter dto.FilterTimeRangeLimit) ([]dto.ConfigCheck, rest_err.APIError) {
	configCheckList, err := c.daoC.FindCheck(ctx, branch, filter, false)
	if err != nil {
		return []dto.ConfigCheck{}, err
	}
	return configCheckList, nil
}
