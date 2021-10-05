package configcheckdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckConfigDaoAssumer interface {
	CheckConfigSaver
	CheckConfigLoader
}

type CheckConfigSaver interface {
	InsertCheck(ctx context.Context, input dto.ConfigCheck) (*string, rest_err.APIError)
	EditCheck(ctx context.Context, input dto.ConfigCheckEdit) (*dto.ConfigCheck, rest_err.APIError)
	DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.ConfigCheck, rest_err.APIError)
	UpdateCheckItem(ctx context.Context, input dto.ConfigCheckItemUpdate) (*dto.ConfigCheck, rest_err.APIError)
	UpdateManyItem(ctx context.Context, input dto.ConfigCheckUpdateMany) rest_err.APIError
}

type CheckConfigLoader interface {
	GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.ConfigCheck, rest_err.APIError)
	FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.ConfigCheck, rest_err.APIError)
	GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string) (*dto.ConfigCheck, rest_err.APIError)
}
