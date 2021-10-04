package stockdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockDaoAssumer interface {
	StockSaver
	StockLoader
}

type StockSaver interface {
	InsertStock(ctx context.Context, input dto.Stock) (*string, rest_err.APIError)
	EditStock(ctx context.Context, input dto.StockEdit) (*dto.Stock, rest_err.APIError)
	DeleteStock(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.Stock, rest_err.APIError)
	DisableStock(ctx context.Context, stockID primitive.ObjectID, user mjwt.CustomClaim, isDisable bool) (*dto.Stock, rest_err.APIError)
	UploadImage(ctx context.Context, stockID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Stock, rest_err.APIError)
	ChangeQtyStock(ctx context.Context, filterA dto.FilterIDBranch, data dto.StockChange) (*dto.Stock, rest_err.APIError)
}

type StockLoader interface {
	GetStockByID(ctx context.Context, stockID primitive.ObjectID, branchIfSpecific string) (*dto.Stock, rest_err.APIError)
	FindStock(ctx context.Context, filterA dto.FilterBranchNameCatDisable) (dto.StockResponseMinList, rest_err.APIError)
	FindStockNeedRestock(ctx context.Context, filterA dto.FilterBranchCatDisable) ([]dto.Stock, rest_err.APIError)
}
