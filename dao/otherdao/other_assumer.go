package otherdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OtherDaoAssumer interface {
	OtherSaver
	OtherLoader
}

type OtherSaver interface {
	InsertOther(ctx context.Context, input dto.Other) (*string, rest_err.APIError)
	EditOther(ctx context.Context, input dto.OtherEdit) (*dto.Other, rest_err.APIError)
	DeleteOther(ctx context.Context, input dto.FilterIDBranchCategoryCreateGte) (*dto.Other, rest_err.APIError)
	DisableOther(ctx context.Context, pcID primitive.ObjectID, user mjwt.CustomClaim, subCategory string, value bool) (*dto.Other, rest_err.APIError)
	UploadImage(ctx context.Context, pcID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Other, rest_err.APIError)
}

type OtherLoader interface {
	GetOtherByID(ctx context.Context, pcID primitive.ObjectID, branchIfSpecific string) (*dto.Other, rest_err.APIError)
	FindOther(ctx context.Context, filter dto.FilterOther) (dto.OtherResponseMinList, rest_err.APIError)
}
