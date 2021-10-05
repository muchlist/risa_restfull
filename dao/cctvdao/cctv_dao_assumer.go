package cctvdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CctvDaoAssumer interface {
	CctvSaver
	CctvLoader
}

type CctvSaver interface {
	InsertCctv(ctx context.Context, input dto.Cctv) (*string, rest_err.APIError)
	EditCctv(ctx context.Context, input dto.CctvEdit) (*dto.Cctv, rest_err.APIError)
	DeleteCctv(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.Cctv, rest_err.APIError)
	DisableCctv(ctx context.Context, cctvID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError)
	UploadImage(ctx context.Context, cctvID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Cctv, rest_err.APIError)
}
type CctvLoader interface {
	GetCctvByID(ctx context.Context, cctvID primitive.ObjectID, branchIfSpecific string) (*dto.Cctv, rest_err.APIError)
	FindCctv(ctx context.Context, filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError)
}
