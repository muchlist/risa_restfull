package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/check_item_dao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewCheckItemService(checkItemDao check_item_dao.CheckItemDaoAssumer) CheckItemServiceAssumer {
	return &checkItemService{
		daoC: checkItemDao,
	}
}

type checkItemService struct {
	daoC check_item_dao.CheckItemDaoAssumer
}
type CheckItemServiceAssumer interface {
	InsertCheckItem(user mjwt.CustomClaim, input dto.CheckItemRequest) (*string, rest_err.APIError)
	EditCheckItem(user mjwt.CustomClaim, checkItemID string, input dto.CheckItemEditRequest) (*dto.CheckItem, rest_err.APIError)
	DeleteCheckItem(user mjwt.CustomClaim, id string) rest_err.APIError
	DisableCheckItem(checkItemID string, user mjwt.CustomClaim, value bool) (*dto.CheckItem, rest_err.APIError)

	GetCheckItemByID(checkItemID string, branchIfSpecific string) (*dto.CheckItem, rest_err.APIError)
	FindCheckItem(filter dto.FilterBranchNameDisable, filterHaveProblem bool) (dto.CheckItemResponseMinList, rest_err.APIError)
}

func (c *checkItemService) InsertCheckItem(user mjwt.CustomClaim, input dto.CheckItemRequest) (*string, rest_err.APIError) {
	// Filling data
	timeNow := time.Now().Unix()
	data := dto.CheckItem{
		CreatedAt:   timeNow,
		CreatedBy:   user.Name,
		CreatedByID: user.Identity,
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		Branch:      user.Branch,
		Name:        input.Name,
		Location:    input.Location,
		LocationLat: input.LocationLat,
		LocationLon: input.LocationLon,
		Tag:         input.Tag,
		TagExtra:    input.TagExtra,
		Type:        input.Type,
		Note:        input.Note,
		Shifts:      input.Shifts,
	}

	// DB
	insertedID, err := c.daoC.InsertCheckItem(data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (c *checkItemService) EditCheckItem(user mjwt.CustomClaim, checkItemID string, input dto.CheckItemEditRequest) (*dto.CheckItem, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(checkItemID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.CheckItemEdit{
		FilterIDBranchTimestamp: dto.FilterIDBranchTimestamp{
			FilterID:        oid,
			FilterBranch:    user.Branch,
			FilterTimestamp: input.FilterTimestamp,
		},
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		Name:        input.Name,
		Location:    input.Location,
		LocationLat: input.LocationLat,
		LocationLon: input.LocationLon,
		Tag:         input.Tag,
		TagExtra:    input.TagExtra,
		Type:        input.Type,
		Note:        input.Note,
		Shifts:      input.Shifts,
	}

	// DB
	checkItemEdited, err := c.daoC.EditCheckItem(data)
	if err != nil {
		return nil, err
	}

	return checkItemEdited, nil
}

func (c *checkItemService) DeleteCheckItem(user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// DB
	_, err := c.daoC.DeleteCheckItem(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	})
	if err != nil {
		return err
	}

	return nil
}

// DisableCheckItem if value true , checkItem will disabled
func (c *checkItemService) DisableCheckItem(checkItemID string, user mjwt.CustomClaim, value bool) (*dto.CheckItem, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(checkItemID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable checkItem
	checkItem, err := c.daoC.DisableCheckItem(oid, user, value)
	if err != nil {
		return nil, err
	}
	return checkItem, nil
}

func (c *checkItemService) GetCheckItemByID(checkItemID string, branchIfSpecific string) (*dto.CheckItem, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(checkItemID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	checkItem, err := c.daoC.GetCheckItemByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return checkItem, nil
}

func (c *checkItemService) FindCheckItem(filter dto.FilterBranchNameDisable, filterHaveProblem bool) (dto.CheckItemResponseMinList, rest_err.APIError) {
	checkItemList, err := c.daoC.FindCheckItem(filter, filterHaveProblem)
	if err != nil {
		return nil, err
	}
	return checkItemList, nil
}
