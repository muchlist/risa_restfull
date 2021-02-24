package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dao/gen_unit_dao"
	"github.com/muchlist/risa_restfull/dao/history_dao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewHistoryService(histDao history_dao.HistoryDaoAssumer, genDao gen_unit_dao.GenUnitDaoAssumer) HistoryServiceAssumer {
	return &historyService{
		daoH: histDao,
		daoG: genDao,
	}
}

type historyService struct {
	daoH history_dao.HistoryDaoAssumer
	daoG gen_unit_dao.GenUnitDaoAssumer
}
type HistoryServiceAssumer interface {
	InsertHistory(user mjwt.CustomClaim, input dto.HistoryRequest) (*string, rest_err.APIError)
	EditHistory(user mjwt.CustomClaim, historyID primitive.ObjectID, input dto.HistoryEditRequest) (*dto.HistoryResponse, rest_err.APIError)

	FindHistory(filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForParent(parentID string) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForUser(userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	GetHistoryCount(branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError)
}

func (h *historyService) InsertHistory(user mjwt.CustomClaim, input dto.HistoryRequest) (*string, rest_err.APIError) {

	// Cek apakah image tersedia untuk ID tersebut TODO
	imagePath := ""

	// Tambahkan Case jika history status bukan Complete, akan gagal jika ID dan Cabang tidak sesuai
	// jika complete gunakan GetUnitByID untuk memastikan ID dan Cabang sesuai
	var parent *dto.GenUnitResponse
	var err rest_err.APIError
	if input.CompleteStatus == enum.HComplete {
		//DB
		parent, err = h.daoG.InsertCase(dto.GenUnitCaseRequest{
			UnitID:       input.ParentID,
			FilterBranch: user.Branch,
			CaseID:       input.ID.Hex(), // gunakan history id sebagai caseID
			CaseNote:     fmt.Sprintf("#%s# %s : %s", enum.GetProgressString(input.CompleteStatus), input.Status, input.Problem),
		})
	} else {
		//DB
		parent, err = h.daoG.GetUnitByID(input.ParentID, user.Branch)
	}
	if err != nil {
		// menggabungkan err dari case insert dengan diawali pesan error tambahan
		combineErr := rest_err.NewBadRequestError(
			fmt.Sprintf("-> %s -> %s", "id unit atau cabang tidak sesuai", err.Message()),
		)
		return nil, combineErr
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.History{
		CreatedAt:      timeNow,
		CreatedBy:      user.Name,
		CreatedByID:    user.Identity,
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Category:       parent.Category,
		Branch:         user.Branch,
		ID:             primitive.ObjectID{},
		ParentID:       input.ParentID,
		ParentName:     parent.Name,
		Status:         input.Status,
		Problem:        input.Problem,
		ProblemResolve: input.ProblemResolve,
		CompleteStatus: input.CompleteStatus,
		DateStart:      input.DateStart,
		DateEnd:        input.DateEnd,
		Tag:            input.Tag,
		Image:          imagePath,
	}

	//DB
	insertedID, err := h.daoH.InsertHistory(data)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
}

func (h *historyService) EditHistory(user mjwt.CustomClaim, historyID primitive.ObjectID, input dto.HistoryEditRequest) (*dto.HistoryResponse, rest_err.APIError) {

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.HistoryEdit{
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		Status:          input.Status,
		Problem:         input.Problem,
		ProblemResolve:  input.ProblemResolve,
		CompleteStatus:  input.CompleteStatus,
		DateEnd:         input.DateEnd,
		Tag:             input.Tag,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
	}

	//DB
	historyEdited, err := h.daoH.EditHistory(historyID, data)
	if err != nil {
		return nil, err
	}

	//Hapus Case pada parrent (jika complete status tidak complete maka perlu ditambahkan lagi case baru)
	// DB
	_, err = h.daoG.DeleteCase(dto.GenUnitCaseRequest{
		UnitID:       historyEdited.ParentID,
		FilterBranch: user.Branch,
		CaseID:       historyID.Hex(),
		CaseNote:     "",
	})
	if err != nil {
		return nil, err
	}

	// jika complete status tidak complete maka perlu ditambahkan lagi case baru
	if input.CompleteStatus != enum.HComplete {
		//DB
		_, err = h.daoG.InsertCase(dto.GenUnitCaseRequest{
			UnitID:       historyEdited.ParentID,
			FilterBranch: user.Branch,
			CaseID:       historyID.Hex(), // gunakan history id sebagai caseID
			CaseNote:     fmt.Sprintf("#%s# %s : %s", enum.GetProgressString(input.CompleteStatus), input.Status, input.Problem),
		})
	}

	return historyEdited, nil
}

func (h *historyService) FindHistory(filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	historyList, err := h.daoH.FindHistory(filterA, filterB)
	if err != nil {
		return nil, err
	}
	return historyList, nil
}

func (h *historyService) FindHistoryForUser(userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	historyList, err := h.daoH.FindHistoryForUser(userID, filter)
	if err != nil {
		return nil, err
	}
	return historyList, nil
}

func (h *historyService) FindHistoryForParent(parentID string) (dto.HistoryResponseMinList, rest_err.APIError) {
	historyList, err := h.daoH.FindHistoryForParent(parentID)
	if err != nil {
		return nil, err
	}
	return historyList, nil
}

func (h *historyService) GetHistoryCount(branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError) {
	historyCountList, err := h.daoH.GetHistoryCount(branchIfSpecific, statusComplete)
	if err != nil {
		return nil, err
	}
	return historyCountList, nil
}
