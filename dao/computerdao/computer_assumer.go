package computerdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComputerDaoAssumer interface {
	ComputerSaver
	ComputerLoader
}

type ComputerSaver interface {
	InsertPc(ctx context.Context, input dto.Computer) (*string, rest_err.APIError)
	EditPc(ctx context.Context, input dto.ComputerEdit) (*dto.Computer, rest_err.APIError)
	DeletePc(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.Computer, rest_err.APIError)
	DisablePc(ctx context.Context, pcID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.Computer, rest_err.APIError)
	UploadImage(ctx context.Context, pcID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Computer, rest_err.APIError)
}

type ComputerLoader interface {
	GetPcByID(ctx context.Context, pcID primitive.ObjectID, branchIfSpecific string) (*dto.Computer, rest_err.APIError)
	FindPc(ctx context.Context, filter dto.FilterComputer) (dto.ComputerResponseMinList, rest_err.APIError)
}
