package service

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/checkitemdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewCheckItemService(checkItemDao checkitemdao.CheckItemDaoAssumer) CheckItemServiceAssumer {
	return &checkItemService{
		daoC: checkItemDao,
	}
}

type checkItemService struct {
	daoC checkitemdao.CheckItemDaoAssumer
}
type CheckItemServiceAssumer interface {
	InsertCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.CheckItemRequest) (*string, rest_err.APIError)
	EditCheckItem(ctx context.Context, user mjwt.CustomClaim, checkItemID string, input dto.CheckItemEditRequest) (*dto.CheckItem, rest_err.APIError)
	DeleteCheckItem(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	DisableCheckItem(ctx context.Context, checkItemID string, user mjwt.CustomClaim, value bool) (*dto.CheckItem, rest_err.APIError)

	GetCheckItemByID(ctx context.Context, checkItemID string, branchIfSpecific string) (*dto.CheckItem, rest_err.APIError)
	FindCheckItem(ctx context.Context, filter dto.FilterBranchNameDisable, filterHaveProblem bool) (dto.CheckItemResponseMinList, rest_err.APIError)
}

func (c *checkItemService) InsertCheckItem(ctx context.Context, user mjwt.CustomClaim, input dto.CheckItemRequest) (*string, rest_err.APIError) {
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
	insertedID, err := c.daoC.InsertCheckItem(ctx, data)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (c *checkItemService) EditCheckItem(ctx context.Context, user mjwt.CustomClaim, checkItemID string, input dto.CheckItemEditRequest) (*dto.CheckItem, rest_err.APIError) {
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
	checkItemEdited, err := c.daoC.EditCheckItem(ctx, data)
	if err != nil {
		return nil, err
	}

	return checkItemEdited, nil
}

func (c *checkItemService) DeleteCheckItem(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// DB
	_, err := c.daoC.DeleteCheckItem(ctx, dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	})
	if err != nil {
		return err
	}

	return nil
}

// DisableCheckItem if value true , checkItem will disabled
func (c *checkItemService) DisableCheckItem(ctx context.Context, checkItemID string, user mjwt.CustomClaim, value bool) (*dto.CheckItem, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(checkItemID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable checkItem
	checkItem, err := c.daoC.DisableCheckItem(ctx, oid, user, value)
	if err != nil {
		return nil, err
	}
	return checkItem, nil
}

func (c *checkItemService) GetCheckItemByID(ctx context.Context, checkItemID string, branchIfSpecific string) (*dto.CheckItem, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(checkItemID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	checkItem, err := c.daoC.GetCheckItemByID(ctx, oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return checkItem, nil
}

func (c *checkItemService) FindCheckItem(ctx context.Context, filter dto.FilterBranchNameDisable, filterHaveProblem bool) (dto.CheckItemResponseMinList, rest_err.APIError) {
	checkItemList, err := c.daoC.FindCheckItem(ctx, filter, filterHaveProblem)
	if err != nil {
		return nil, err
	}
	return checkItemList, nil
}
