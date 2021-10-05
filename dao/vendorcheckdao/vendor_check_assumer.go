package vendorcheckdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckVendorDaoAssumer interface {
	CheckVendorSaver
	CheckVendorLoader
}

type CheckVendorSaver interface {
	InsertCheck(ctx context.Context, input dto.VendorCheck) (*string, rest_err.APIError)
	EditCheck(ctx context.Context, input dto.VendorCheckEdit) (*dto.VendorCheck, rest_err.APIError)
	DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.VendorCheck, rest_err.APIError)
	UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.VendorCheck, rest_err.APIError)
	UpdateCheckItem(ctx context.Context, input dto.VendorCheckItemUpdate) (*dto.VendorCheck, rest_err.APIError)
	BulkUpdateItem(ctx context.Context, inputs []dto.VendorCheckItemUpdate) (int64, rest_err.APIError)
}

type CheckVendorLoader interface {
	GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.VendorCheck, rest_err.APIError)
	FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.VendorCheck, rest_err.APIError)
	GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string) (*dto.VendorCheck, rest_err.APIError)
}
