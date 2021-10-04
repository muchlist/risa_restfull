package improvedao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImproveDaoAssumer interface {
	ImproveSaver
	ImproveLoader
}

type ImproveSaver interface {
	InsertImprove(ctx context.Context, input dto.Improve) (*string, rest_err.APIError)
	EditImprove(ctx context.Context, input dto.ImproveEdit) (*dto.Improve, rest_err.APIError)
	ChangeImprove(ctx context.Context, filterA dto.FilterIDBranch, data dto.ImproveChange) (*dto.Improve, rest_err.APIError)
	DeleteImprove(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.Improve, rest_err.APIError)
	ActivateImprove(ctx context.Context, improveID primitive.ObjectID, user mjwt.CustomClaim, isEnable bool) (*dto.Improve, rest_err.APIError)
}

type ImproveLoader interface {
	GetImproveByID(ctx context.Context, improveID primitive.ObjectID, branchIfSpecific string) (*dto.Improve, rest_err.APIError)
	FindImprove(ctx context.Context, filterA dto.FilterBranchCompleteTimeRangeLimit) (dto.ImproveResponseMinList, rest_err.APIError)
}
