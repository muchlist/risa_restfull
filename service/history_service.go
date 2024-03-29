package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/userdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewHistoryService(
	histDao historydao.HistoryDaoAssumer,
	genDao genunitdao.GenUnitDaoAssumer,
	userDao userdao.UserDaoAssumer,
	fcmClient fcm.ClientAssumer) HistoryServiceAssumer {
	return &historyService{
		daoH:      histDao,
		daoG:      genDao,
		daoU:      userDao,
		fcmClient: fcmClient,
	}
}

type historyService struct {
	daoH      historydao.HistoryDaoAssumer
	daoG      genunitdao.GenUnitDaoAssumer
	daoU      userdao.UserDaoAssumer
	fcmClient fcm.ClientAssumer
}
type HistoryServiceAssumer interface {
	InsertHistory(ctx context.Context, user mjwt.CustomClaim, input dto.HistoryRequest) (*string, rest_err.APIError)
	EditHistory(ctx context.Context, user mjwt.CustomClaim, historyID string, input dto.HistoryEditRequest) (*dto.HistoryResponse, rest_err.APIError)
	DeleteHistory(ctx context.Context, user mjwt.CustomClaim, id string, force bool) rest_err.APIError
	PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.HistoryResponse, rest_err.APIError)

	GetHistory(ctx context.Context, parentID string, branchIfSpecific string) (*dto.HistoryResponse, rest_err.APIError)
	FindHistory(ctx context.Context, search string, filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForHome(ctx context.Context, filterA dto.FilterBranchCatComplete) (dto.HistoryResponseMinList, rest_err.APIError)
	UnwindHistory(ctx context.Context, filterA dto.FilterBranchCatInCompleteIn, filterB dto.FilterTimeRangeLimit) (dto.HistoryUnwindResponseList, rest_err.APIError)
	FindHistoryForParent(ctx context.Context, parentID string) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForUser(ctx context.Context, userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	GetHistoryCount(ctx context.Context, branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError)
}

func (h *historyService) InsertHistory(ctx context.Context, user mjwt.CustomClaim, input dto.HistoryRequest) (*string, rest_err.APIError) {
	// Default value
	timeNow := time.Now().Unix()
	if input.DateStart == 0 {
		input.DateStart = timeNow
	}

	// jika ID tersedia, gunakan ID , jika tidak buatkan objectID
	// memastikan ID yang diinputkan bisa diubah ke ObjectID
	generatedID := primitive.NewObjectID()
	var errT error
	if input.ID != "" {
		generatedID, errT = primitive.ObjectIDFromHex(input.ID)
		if errT != nil {
			return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan dari frontend salah")
		}
	}

	// cek jika user biasa (bukan Approver) tidak boleh menambahkan history
	// History-with-note dan pending, alihkan ke req pending atau req complete-history
	if !sfunc.InSlice(roles.RoleApprove, user.Roles) {
		if input.CompleteStatus == enum.HPending {
			input.CompleteStatus = enum.HRequestPending
		}
		if input.CompleteStatus == enum.HCompleteWithBA {
			input.CompleteStatus = enum.HRequestComplete
		}
	}

	// Mengambil gen_unit
	// Tambahkan Case jika History status bukan Complete atau bukan info, akan gagal jika ID dan Cabang tidak sesuai
	// jika complete gunakan GetUnitByID untuk memastikan ID dan Cabang sesuai
	var parent *dto.GenUnitResponse
	var err rest_err.APIError
	historyIsComplete := input.CompleteStatus == enum.HComplete
	historyIsInfo := input.CompleteStatus == enum.HInfo
	historyIsDataInfo := input.CompleteStatus == enum.HDataInfo
	if !(historyIsComplete || historyIsInfo || historyIsDataInfo) {
		// DB
		parent, err = h.daoG.InsertCase(ctx, dto.GenUnitCaseRequest{
			UnitID:       input.ParentID,
			FilterBranch: user.Branch,
			CaseID:       generatedID.Hex(), // gunakan History id sebagai caseID
			CaseNote:     fmt.Sprintf("#%s# %s : %s", enum.GetProgressString(input.CompleteStatus), input.Status, input.Problem),
		})
	} else {
		// DB
		parent, err = h.daoG.GetUnitByID(ctx, input.ParentID, user.Branch)
	}
	if err != nil {
		// menggabungkan err dari case insert dengan diawali pesan error tambahan
		combineErr := rest_err.NewBadRequestError(
			fmt.Sprintf("-> %s -> %s", "id unit atau cabang tidak sesuai", err.Message()),
		)
		return nil, combineErr
	}

	// Filling data
	data := dto.History{
		ID:             generatedID,
		CreatedAt:      timeNow,
		CreatedBy:      user.Name,
		CreatedByID:    user.Identity,
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Category:       parent.Category,
		Branch:         user.Branch,
		ParentID:       input.ParentID,
		ParentName:     parent.Name,
		Status:         input.Status,
		Problem:        input.Problem,
		ProblemResolve: input.ProblemResolve,
		CompleteStatus: input.CompleteStatus,
		DateStart:      input.DateStart,
		DateEnd:        input.DateEnd,
		Tag:            input.Tag,
		Image:          input.Image,
	}

	isVendor := sfunc.InSlice(roles.RoleVendor, user.Roles)

	// DB
	insertedID, err := h.daoH.InsertHistory(ctx, data, isVendor)
	if err != nil {
		return nil, err
	}

	// goroutine notifikasi
	go func() {

		// jika status complete datanya adalah dataInfo skip notifikasi
		if historyIsDataInfo {
			return
		}

		users, err := h.daoU.FindUser(ctx, user.Branch)
		if err != nil {
			logger.Error("mendapatkan user gagal saat menambahkan fcm (INSERT HISTORY)", err)
		}

		var tokens []string
		for _, u := range users {
			// skip jika user == user yang mengirim
			if u.ID == user.Identity {
				continue
			}
			// skip jika category bukan cctv dan rolenya vendor
			if strings.ToUpper(data.Category) != category.Cctv && sfunc.InSlice(roles.RoleVendor, u.Roles) {
				continue
			}
			tokens = append(tokens, u.FcmToken)
		}
		// firebase
		h.fcmClient.SendMessage(fcm.Payload{
			Title:          fmt.Sprintf("Insiden %s ditambahkan", strings.ToLower(parent.Category)),
			Message:        fmt.Sprintf("%s :: %s - %s :: oleh %s", enum.GetProgressString(input.CompleteStatus), input.Problem, input.ProblemResolve, strings.ToLower(user.Name)),
			ReceiverTokens: tokens,
		})
	}()

	return insertedID, nil
}

func (h *historyService) EditHistory(ctx context.Context, user mjwt.CustomClaim, historyID string, input dto.HistoryEditRequest) (*dto.HistoryResponse, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(historyID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	if !sfunc.InSlice(roles.RoleApprove, user.Roles) {
		if input.CompleteStatus == enum.HPending {
			input.CompleteStatus = enum.HRequestPending
		}
		if input.CompleteStatus == enum.HCompleteWithBA {
			input.CompleteStatus = enum.HRequestComplete
		}
	}

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

	isVendor := sfunc.InSlice(roles.RoleVendor, user.Roles)

	// DB
	historyEdited, err := h.daoH.EditHistory(ctx, oid, data, isVendor)
	if err != nil {
		return nil, err
	}

	// Hapus Case pada parent
	// DB
	_, err = h.daoG.DeleteCase(ctx, dto.GenUnitCaseRequest{
		UnitID:       historyEdited.ParentID,
		FilterBranch: user.Branch,
		CaseID:       historyID,
		CaseNote:     "",
	})
	if err != nil {
		return nil, err
	}

	// jika complete_status tidak complete atau tidak info maka perlu ditambahkan lagi case baru
	historyIsComplete := input.CompleteStatus == enum.HComplete
	historyIsInfo := input.CompleteStatus == enum.HInfo
	if !(historyIsComplete || historyIsInfo) {
		// DB
		_, err = h.daoG.InsertCase(ctx, dto.GenUnitCaseRequest{
			UnitID:       historyEdited.ParentID,
			FilterBranch: user.Branch,
			CaseID:       historyID, // gunakan History id sebagai caseID
			CaseNote:     fmt.Sprintf("#%s# %s : %s", enum.GetProgressString(input.CompleteStatus), input.Status, input.Problem),
		})
		if err != nil {
			return nil, err
		}
	}

	go func() {
		users, err := h.daoU.FindUser(ctx, user.Branch)
		if err != nil {
			logger.Error("mendapatkan user gagal saat menambahkan fcm (EDIT HISTORY)", err)
		}

		var tokens []string

		for _, u := range users {
			// skip jika user == user yang mengirim
			if u.ID == user.Identity {
				continue
			}
			// skip jika category bukan cctv dan rolenya vendor
			if strings.ToUpper(historyEdited.Category) != category.Cctv && sfunc.InSlice(roles.RoleVendor, u.Roles) {
				continue
			}
			tokens = append(tokens, u.FcmToken)
		}

		// firebase
		h.fcmClient.SendMessage(fcm.Payload{
			Title:          fmt.Sprintf("History %s - %s diupdate", strings.ToLower(historyEdited.Category), strings.ToLower(historyEdited.ParentName)),
			Message:        fmt.Sprintf("%s :: %s - %s :: oleh %s", enum.GetProgressString(input.CompleteStatus), input.Problem, input.ProblemResolve, strings.ToLower(user.Name)),
			ReceiverTokens: tokens,
		})
	}()

	return historyEdited, nil
}

func (h *historyService) DeleteHistory(ctx context.Context, user mjwt.CustomClaim, id string, force bool) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1).Unix()
	if force {
		timeMinusOneDay = 0
	}
	// if admin can delete history without limit time
	if sfunc.InSlice(roles.RoleAdmin, user.Roles) {
		timeMinusOneDay = 0
	}

	// DB
	history, err := h.daoH.DeleteHistory(ctx, dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay,
	})
	if err != nil {
		return err
	}

	// Jika History yang dihapus tidak complete, berarti harus dihapus di parentnya karena masih ada sebagai case
	if history.CompleteStatus != enum.HComplete {
		// DB
		_, err = h.daoG.DeleteCase(ctx, dto.GenUnitCaseRequest{
			UnitID:       history.ParentID,
			FilterBranch: user.Branch,
			CaseID:       id,
			CaseNote:     "",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *historyService) GetHistory(ctx context.Context, parentID string, branchIfSpecific string) (*dto.HistoryResponse, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(parentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	history, err := h.daoH.GetHistoryByID(ctx, oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (h *historyService) FindHistory(ctx context.Context, search string, filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	historyList := dto.HistoryResponseMinList{}
	var err rest_err.APIError
	if len(search) != 0 {
		historyList, err = h.daoH.SearchHistory(ctx, search, filterA, filterB)
		if err != nil {
			return nil, err
		}
	} else {
		historyList, err = h.daoH.FindHistory(ctx, filterA, filterB)
		if err != nil {
			return nil, err
		}
	}

	return historyList, nil
}

// FindHistoryForHome get history complete + progress + pending + reqPending
func (h *historyService) FindHistoryForHome(ctx context.Context, filterA dto.FilterBranchCatComplete) (dto.HistoryResponseMinList, rest_err.APIError) {

	// kembalian dari golang channel
	type result struct {
		res dto.HistoryResponseMinList
		err rest_err.APIError
	}

	resultChan := make(chan result, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	getCompleteHistory := func() {
		defer wg.Done()
		// DB
		historyListCompleteInfo, err := h.daoH.FindHistory(ctx,
			dto.FilterBranchCatComplete{
				FilterBranch:         filterA.FilterBranch,
				FilterCategory:       filterA.FilterCategory,
				FilterCompleteStatus: []int{enum.HComplete, enum.HInfo},
			},
			dto.FilterTimeRangeLimit{
				Limit: 100,
			},
		)
		resultChan <- result{
			res: historyListCompleteInfo,
			err: err,
		}
	}

	getProgressHistory := func() {
		defer wg.Done()
		// DB
		historyListProgressPending, err := h.daoH.FindHistory(ctx,
			dto.FilterBranchCatComplete{
				FilterBranch:         filterA.FilterBranch,
				FilterCategory:       filterA.FilterCategory,
				FilterCompleteStatus: []int{enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete, enum.HCompleteWithBA},
			},
			dto.FilterTimeRangeLimit{
				Limit: 200,
			},
		)
		resultChan <- result{
			res: historyListProgressPending,
			err: err,
		}
	}

	go getCompleteHistory()
	go getProgressHistory()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	resultTemp := make([]dto.HistoryResponseMin, 0)
	var errString string
	for v := range resultChan {
		if v.err != nil {
			errString = errString + v.err.Message() + ". "
		}
		if len(v.res) != 0 {
			resultTemp = append(resultTemp, v.res...)
		}
	}

	if errString != "" {
		return nil, rest_err.NewInternalServerError(errString, errors.New("gagal mendapatkan data (Service: FindHistoryForHome)"))
	}

	// sort
	sort.Slice(resultTemp, func(i, j int) bool {
		return resultTemp[i].UpdatedAt > resultTemp[j].UpdatedAt
	})

	return resultTemp, nil
}

func (h *historyService) UnwindHistory(ctx context.Context, filterA dto.FilterBranchCatInCompleteIn, filterB dto.FilterTimeRangeLimit) (dto.HistoryUnwindResponseList, rest_err.APIError) {
	historyList, err := h.daoH.UnwindHistory(ctx, filterA, filterB)
	if err != nil {
		return nil, err
	}
	return historyList, nil
}

func (h *historyService) FindHistoryForUser(ctx context.Context, userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	historyList, err := h.daoH.FindHistoryForUser(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	return historyList, nil
}

func (h *historyService) FindHistoryForParent(ctx context.Context, parentID string) (dto.HistoryResponseMinList, rest_err.APIError) {
	historyList, err := h.daoH.FindHistoryForParent(ctx, parentID)
	if err != nil {
		return nil, err
	}
	return historyList, nil
}

func (h *historyService) GetHistoryCount(ctx context.Context, branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError) {
	historyCountList, err := h.daoH.GetHistoryCount(ctx, branchIfSpecific, statusComplete)
	if err != nil {
		return nil, err
	}
	return historyCountList, nil
}

// PutImage memasukkan lokasi file (path) ke dalam database History dengan mengecek kesesuaian branch
func (h *historyService) PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.HistoryResponse, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	history, err := h.daoH.UploadImage(ctx, oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return history, nil
}
