package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dao/check_dao"
	"github.com/muchlist/risa_restfull/dao/check_item_dao"
	"github.com/muchlist/risa_restfull/dao/gen_unit_dao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewCheckService(checkDao check_dao.CheckDaoAssumer,
	checkItemDao check_item_dao.CheckItemDaoAssumer,
	genUnitDao gen_unit_dao.GenUnitDaoAssumer) CheckServiceAssumer {
	return &checkService{
		daoC:  checkDao,
		daoCI: checkItemDao,
		daoG:  genUnitDao,
	}
}

type checkService struct {
	daoC  check_dao.CheckDaoAssumer
	daoCI check_item_dao.CheckItemDaoAssumer
	daoG  gen_unit_dao.GenUnitDaoAssumer
}
type CheckServiceAssumer interface {
	InsertCheck(user mjwt.CustomClaim, input dto.CheckRequest) (*string, rest_err.APIError)
	EditCheck(user mjwt.CustomClaim, checkID string, input dto.CheckEditRequest) (*dto.Check, rest_err.APIError)
	DeleteCheck(user mjwt.CustomClaim, id string) rest_err.APIError
	PutChildImage(user mjwt.CustomClaim, parentId string, childId string, imagePath string) (*dto.Check, rest_err.APIError)

	//GetCheckByID(checkID string, branchIfSpecific string) (*dto.Check, rest_err.APIError)
	//FindCheck(filter dto.FilterBranchLocIPNameDisable) (dto.CheckResponseMinList, rest_err.APIError)
}

func (c *checkService) InsertCheck(user mjwt.CustomClaim, input dto.CheckRequest) (*string, rest_err.APIError) {

	// IMPREVEMENT gunakan goroutine untuk mengambil data dari dua dao
	// ambil check item berdasarkan cabang yang di input
	checkItems, err := c.daoCI.FindCheckItem(dto.FilterBranchNameDisable{
		FilterBranch: user.Branch,
	}, false)
	if err != nil {
		return nil, err
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

	// ambil data cctv yang harus di cek
	cctvList, err := c.daoG.FindUnit(dto.GenUnitFilter{
		Branch:   user.Branch,
		Category: category.Cctv,
		Disable:  false,
		Pings:    true,
		LastPing: enum.GetPingString(enum.PingDown),
	})
	if err != nil {
		return nil, err
	}
	for _, cctv := range cctvList {
		cctvHaveZeroCase := cctv.CasesSize == 0
		cctvPing2IsDown := true
		if len(cctv.PingsState) > 1 {
			// memeriksa ping index 1 (ping kedua) karena ping pertama sudah pasti 0
			// berdasarkan filter FindUnit
			cctvPing2IsDown = cctv.PingsState[1].Code == 0
		}
		if cctvHaveZeroCase && cctvPing2IsDown {
			checkItemsSelected = append(checkItemsSelected, dto.CheckItemEmbed{
				ID:          cctv.ID,
				Name:        cctv.Name,
				Type:        category.Cctv,
				Tag:         []string{},
				TagExtra:    []string{},
				HaveProblem: true,
			})
		}
	}

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

	//DB
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

	//DB
	checkEdited, err := c.daoC.EditCheck(data)
	if err != nil {
		return nil, err
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

//PutImage memasukkan lokasi file (path) ke dalam database check dengan mengecek kesesuaian branch
func (c *checkService) PutChildImage(user mjwt.CustomClaim, parentId string, childId string, imagePath string) (*dto.Check, rest_err.APIError) {

	parentOid, errT := primitive.ObjectIDFromHex(parentId)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Parent ObjectID yang dimasukkan salah")
	}

	filter := dto.FilterParentIDChildIDAuthor{
		FilterParentID: parentOid,
		FilterChildID:  childId,
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

	// Mengupdate value di check item, agar pada
	childOid, errT := primitive.ObjectIDFromHex(input.ChildID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("Child ObjectID yang dimasukkan salah")
	}

	_, err = c.daoCI.EditCheckItemValue(
		dto.CheckItemEditBySys{
			FilterID:       childOid,
			UpdatedAt:      0,
			CheckedNote:    "",
			HaveProblem:    false,
			CompleteStatus: 0,
		})

	return check, nil
}

//
//func (c *checkService) GetCheckByID(checkID string, branchIfSpecific string) (*dto.Check, rest_err.APIError) {
//	oid, errT := primitive.ObjectIDFromHex(checkID)
//	if errT != nil {
//		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
//	}
//
//	check, err := c.daoC.GetCheckByID(oid, branchIfSpecific)
//	if err != nil {
//		return nil, err
//	}
//	return check, nil
//}
//
//func (c *checkService) FindCheck(filter dto.FilterBranchLocIPNameDisable) (dto.CheckResponseMinList, rest_err.APIError) {
//
//	// cek apakah ip address valid, jika valid maka set filter.FilterName ke kosong supaya pencarian berdasarkan IP
//	if filter.FilterIP != "" {
//		if net.ParseIP(filter.FilterIP) == nil {
//			return nil, rest_err.NewBadRequestError("IP Address tidak valid")
//		}
//		filter.FilterName = ""
//	}
//
//	checkList, err := c.daoC.FindCheck(filter)
//	if err != nil {
//		return nil, err
//	}
//	return checkList, nil
//}
