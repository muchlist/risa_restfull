package checkdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckDaoAssumer interface {
	CheckSaver
	CheckLoader
}

type CheckSaver interface {
	InsertCheck(ctx context.Context, input dto.Check) (*string, rest_err.APIError)
	EditCheck(ctx context.Context, input dto.CheckEdit) (*dto.Check, rest_err.APIError)
	DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.Check, rest_err.APIError)
	UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.Check, rest_err.APIError)
	UpdateCheckItem(ctx context.Context, input dto.CheckChildUpdate) (*dto.Check, rest_err.APIError)
}

type CheckLoader interface {
	GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.Check, rest_err.APIError)
	FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit) (dto.CheckResponseMinList, rest_err.APIError)
	FindCheckForReports(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit) ([]dto.Check, rest_err.APIError)
}
