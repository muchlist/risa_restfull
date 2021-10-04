package checkitemdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckItemDaoAssumer interface {
	CheckItemSaver
	CheckItemLoader
}

type CheckItemSaver interface {
	InsertCheckItem(ctx context.Context, input dto.CheckItem) (*string, rest_err.APIError)
	EditCheckItem(ctx context.Context, input dto.CheckItemEdit) (*dto.CheckItem, rest_err.APIError)
	EditCheckItemValue(ctx context.Context, input dto.CheckItemEditBySys) (*dto.CheckItem, rest_err.APIError)
	DeleteCheckItem(ctx context.Context, input dto.FilterIDBranch) (*dto.CheckItem, rest_err.APIError)
	DisableCheckItem(ctx context.Context, checkItemID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.CheckItem, rest_err.APIError)
}

type CheckItemLoader interface {
	GetCheckItemByID(ctx context.Context, checkItemID primitive.ObjectID, branchIfSpecific string) (*dto.CheckItem, rest_err.APIError)
	FindCheckItem(ctx context.Context, filterA dto.FilterBranchNameDisable, filterHaveProblem bool) (dto.CheckItemResponseMinList, rest_err.APIError)
}
