package altaicheckdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckAltaiDaoAssumer interface {
	CheckAltaiSaver
	CheckAltaiLoader
}

type CheckAltaiSaver interface {
	InsertCheck(ctx context.Context, input dto.AltaiCheck) (*string, rest_err.APIError)
	EditCheck(ctx context.Context, input dto.AltaiCheckEdit) (*dto.AltaiCheck, rest_err.APIError)
	DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.AltaiCheck, rest_err.APIError)
	UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.AltaiCheck, rest_err.APIError)
	UpdateCheckItem(ctx context.Context, input dto.AltaiCheckItemUpdate) (*dto.AltaiCheck, rest_err.APIError)
	BulkUpdateItem(ctx context.Context, inputs []dto.AltaiCheckItemUpdate) (int64, rest_err.APIError)
}

type CheckAltaiLoader interface {
	GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.AltaiCheck, rest_err.APIError)
	FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.AltaiCheck, rest_err.APIError)
	GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string) (*dto.AltaiCheck, rest_err.APIError)
}
