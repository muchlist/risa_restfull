package altaiphycheckdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckAltaiPhyDaoAssumer interface {
	CheckAltaiPhySaver
	CheckAltaiPhyLoader
}

type CheckAltaiPhySaver interface {
	InsertCheck(ctx context.Context, input dto.AltaiPhyCheck) (*string, rest_err.APIError)
	EditCheck(ctx context.Context, input dto.AltaiPhyCheckEdit) (*dto.AltaiPhyCheck, rest_err.APIError)
	DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.AltaiPhyCheck, rest_err.APIError)
	UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.AltaiPhyCheck, rest_err.APIError)
	UpdateCheckItem(ctx context.Context, input dto.AltaiPhyCheckItemUpdate) (*dto.AltaiPhyCheck, rest_err.APIError)
	BulkUpdateItem(ctx context.Context, inputs []dto.AltaiPhyCheckItemUpdate) (int64, rest_err.APIError)
}

type CheckAltaiPhyLoader interface {
	GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.AltaiPhyCheck, rest_err.APIError)
	FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.AltaiPhyCheck, rest_err.APIError)
	FindCheckStillOpen(ctx context.Context, branch string, detail bool) ([]dto.AltaiPhyCheck, rest_err.APIError)
	GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string, isQuarter bool) (*dto.AltaiPhyCheck, rest_err.APIError)
}
