package service

import (
	"errors"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/checkitemdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewCheckService(checkDao checkdao.CheckDaoAssumer,
	checkItemDao checkitemdao.CheckItemDaoAssumer,
	genUnitDao genunitdao.GenUnitDaoAssumer,
	histService HistoryServiceAssumer,
) CheckServiceAssumer {
	return &checkService{
		daoC:        checkDao,
		daoCI:       checkItemDao,
		daoG:        genUnitDao,
		servHistory: histService,
	}
}

type checkService struct {
	daoC        checkdao.CheckDaoAssumer
	daoCI       checkitemdao.CheckItemDaoAssumer
	daoG        genunitdao.GenUnitDaoAssumer
	servHistory HistoryServiceAssumer
}
type CheckServiceAssumer interface {
	InsertCheck(user mjwt.CustomClaim, input dto.CheckRequest) (*string, rest_err.APIError)
	EditCheck(user mjwt.CustomClaim, checkID string, input dto.CheckEditRequest) (*dto.Check, rest_err.APIError)
	DeleteCheck(user mjwt.CustomClaim, id string) rest_err.APIError
	UpdateCheckItem(user mjwt.CustomClaim, input dto.CheckChildUpdateRequest) (*dto.Check, rest_err.APIError)
	PutChildImage(user mjwt.CustomClaim, parentID string, childID string, imagePath string) (*dto.Check, rest_err.APIError)

	GetCheckByID(checkID string, branchIfSpecific string) (*dto.Check, rest_err.APIError)
	FindCheck(branch string, filter dto.FilterTimeRangeLimit) (dto.CheckResponseMinList, rest_err.APIError)
}

func (c *checkService) InsertCheck(user mjwt.CustomClaim, input dto.CheckRequest) (*string, rest_err.APIError) {
	itemResultCheckItemChan := make(chan []dto.CheckItemEmbed)
	itemResultCctvChan := make(chan []dto.CheckItemEmbed)

	go func() {
		// ambil check item berdasarkan cabang yang di input
		checkItems, err := c.daoCI.FindCheckItem(dto.FilterBranchNameDisable{
			FilterBranch: user.Branch,
		}, false)
		if err != nil {
			// if error, send result 0 slice
			itemResultCheckItemChan <- []dto.CheckItemEmbed{}
			return
		}
		// filter check item yang memiliki shift sama dengan input dan memeliki problem
		var checkItemsSelected []dto.CheckItemEmbed
		for _, v := range checkItems {
			shiftMatch := sfunc.IntInSlice(input.Shift, v.Shifts)
			if shiftMatch || v.HaveProblem {
				checkItemsSelected = append(checkItemsSelected, dto.CheckItemEmbed{
					ID:             v.ID.Hex(),
					Name:           v.Name,
					Location:       v.Location,
					Type:           v.Type,
					Tag:            v.Tag, // tag
					TagExtra:       v.TagExtra,
					HaveProblem:    v.HaveProblem,
					CompleteStatus: v.CompleteStatus,
				})
			}
		}

		itemResultCheckItemChan <- checkItemsSelected
	}()

	go func() {
		// ambil data cctv yang harus di cek
		cctvList, err := c.daoG.FindUnit(dto.GenUnitFilter{
			Branch:   user.Branch,
			Category: category.Cctv,
			Disable:  false,
			Pings:    true,
			LastPing: enum.GetPingString(enum.PingDown),
		})
		if err != nil {
			// if error, send result 0 slice
			itemResultCctvChan <- []dto.CheckItemEmbed{}
			return
		}

		var checkCctvSelected []dto.CheckItemEmbed
		for _, cctv := range cctvList {
			cctvHaveZeroCase := cctv.CasesSize == 0
			cctvPing2IsDown := true
			if len(cctv.PingsState) > 1 {
				// memeriksa ping index 1 (ping kedua) karena ping pertama sudah pasti 0
				// berdasarkan filter FindUnit
				cctvPing2IsDown = cctv.PingsState[1].Code == 0
			}
			if cctvHaveZeroCase && cctvPing2IsDown {
				checkCctvSelected = append(checkCctvSelected, dto.CheckItemEmbed{
					ID:          cctv.ID,
					Name:        cctv.Name,
					Type:        category.Cctv,
					Tag:         []string{},
					TagExtra:    []string{},
					HaveProblem: true,
				})
			}
		}
		itemResultCctvChan <- checkCctvSelected
	}()

	checkItemsSelected := <-itemResultCheckItemChan
	checkCctvSelected := <-itemResultCctvChan

	checkItemsSelected = append(checkItemsSelected, checkCctvSelected...)

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.Check{
		CreatedAt:   timeNow,
		CreatedBy:   user.Name,
		CreatedByID: user.Identity,
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		Branch:      user.Branch,
		Shift:       input.Shift,
		CheckItems:  checkItemsSelected,
	}

	// DB
	insertedID, err := c.daoC.InsertCheck(data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (c *checkService) EditCheck(user mjwt.CustomClaim, checkID string, input dto.CheckEditRequest) (*dto.Check, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(checkID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.CheckEdit{
		FilterIDBranchAuthor: dto.FilterIDBranchAuthor{
			FilterID:       oid,
			FilterBranch:   user.Branch,
			FilterAuthorID: user.Identity,
		},
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		IsFinish:    input.IsFinish,
		Note:        input.Note,
	}

	// DB
	checkEdited, err := c.daoC.EditCheck(data)
	if err != nil {
		return nil, err
	}

	// IMPROVEMENT : make looping insert history to use goroutine
	// Jika isFinish true , maka masukkan semua checkItem yang bertipe cctv
	// looping ke insert history
	var errorList []rest_err.APIError
	if checkEdited.IsFinish {
		for _, checkItem := range checkEdited.CheckItems {
			if checkItem.Type == category.Cctv {
				// jika checkItemnya tidak di check lewati
				if !checkItem.IsChecked {
					continue
				}
				// cek complete status tidak boleh 0 atau 3, set default ke 1
				if !(checkItem.CompleteStatus == enum.HComplete) {
					checkItem.CompleteStatus = enum.HProgress
				}
				dataHistory := dto.HistoryRequest{
					ParentID:       checkItem.ID,
					Status:         "Checklist",
					Problem:        checkItem.CheckedNote,
					ProblemResolve: "",
					CompleteStatus: checkItem.CompleteStatus,
					DateStart:      timeNow,
					DateEnd:        timeNow,
					Tag:            []string{},
				}
				_, err := c.servHistory.InsertHistory(user, dataHistory)
				if err != nil {
					errorList = append(errorList, err)
				}
			}
		}
	}

	// mengkoleksi semua error hasil looping insert history
	errMessage := ""
	if len(errorList) != 0 {
		for _, err := range errorList {
			errMessage = errMessage + ". " + err.Message()
		}
		logger.Error(fmt.Sprintf("Check berhasil diubah namun menambahkan history cctv gagal (EditCheck isFinish) : %s", errMessage), errors.New("internal error"))
		return nil, rest_err.NewInternalServerError("Check berhasil diubah namun menambahkan history cctv gagal", errors.New("internal error"))
	}

	return checkEdited, nil
}

func (c *checkService) DeleteCheck(user mjwt.CustomClaim, id string) rest_err.APIError {
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

// PutImage memasukkan lokasi file (path) ke dalam database check dengan mengecek kesesuaian branch
func (c *checkService) PutChildImage(user mjwt.CustomClaim, parentID string, childID string, imagePath string) (*dto.Check, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(parentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	filter := dto.FilterParentIDChildIDAuthor{
		FilterParentID: parentOid,
		FilterChildID:  childID,
		FilterAuthorID: user.Identity,
	}

	check, err := c.daoC.UploadChildImage(filter, imagePath)
	if err != nil {
		return nil, err
	}
	return check, nil
}

func (c *checkService) UpdateCheckItem(user mjwt.CustomClaim, input dto.CheckChildUpdateRequest) (*dto.Check, rest_err.APIError) {
	parentOid, errT := primitive.ObjectIDFromHex(input.ParentID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	timeNow := time.Now().Unix()

	// DB
	data := dto.CheckChildUpdate{
		FilterParentIDChildIDAuthor: dto.FilterParentIDChildIDAuthor{
			FilterParentID: parentOid,
			FilterChildID:  input.ChildID,
			FilterAuthorID: user.Identity,
		},
		UpdatedAt:        timeNow,
		CheckedAt:        timeNow,
		IsChecked:        input.IsChecked,
		TagSelected:      input.TagSelected,
		TagExtraSelected: input.TagExtraSelected,
		CheckedNote:      input.CheckedNote,
		HaveProblem:      input.HaveProblem,
		CompleteStatus:   input.CompleteStatus,
	}
	check, err := c.daoC.UpdateCheckItem(data)
	if err != nil {
		return nil, err
	}

	// Cek index dan Type childID yang diupdate
	var updatedType string
	// var indexItems int //index digunakan untuk memepermudah mendapatkan nama Checkitem yang diupdate
	for _, v := range check.CheckItems {
		if v.ID == input.ChildID {
			updatedType = v.Type
			// indexItems = i
		}
	}

	// Mengupdate value di check item, agar pada pembuatan check berikutnya pesan tetap berlanjut
	// Kecuali cctv yang mana tidak memiliki check item
	childOid, errT := primitive.ObjectIDFromHex(input.ChildID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Child ObjectID yang dimasukkan salah")
	}

	if updatedType != category.Cctv {
		_, err = c.daoCI.EditCheckItemValue(
			dto.CheckItemEditBySys{
				FilterID:       childOid,
				UpdatedAt:      0,
				CheckedNote:    "",
				HaveProblem:    false,
				CompleteStatus: 0,
			})
		if err != nil {
			return nil, err
		}
	}

	return check, nil
}

func (c *checkService) GetCheckByID(checkID string, branchIfSpecific string) (*dto.Check, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(checkID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	check, err := c.daoC.GetCheckByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return check, nil
}

func (c *checkService) FindCheck(branch string, filter dto.FilterTimeRangeLimit) (dto.CheckResponseMinList, rest_err.APIError) {
	checkList, err := c.daoC.FindCheck(branch, filter)
	if err != nil {
		return nil, err
	}
	return checkList, nil
}
