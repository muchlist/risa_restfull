package genunitdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
)

type GenUnitDaoAssumer interface {
	GenUnitSaver
	GenUnitLoader
}

type GenUnitSaver interface {
	InsertUnit(ctx context.Context, unit dto.GenUnit) (*string, rest_err.APIError)
	EditUnit(ctx context.Context, unitID string, unitRequest dto.GenUnitEditRequest) (*dto.GenUnitResponse, rest_err.APIError)
	DeleteUnit(ctx context.Context, unitID string) rest_err.APIError
	InsertCase(ctx context.Context, payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	DeleteCase(ctx context.Context, payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	DisableUnit(ctx context.Context, unitID string, value bool) (*dto.GenUnitResponse, rest_err.APIError)
	AppendPingState(ctx context.Context, input dto.GenUnitPingStateRequest) (int64, rest_err.APIError)
}

type GenUnitLoader interface {
	GetUnitByID(ctx context.Context, unitID string, branchSpecific string) (*dto.GenUnitResponse, rest_err.APIError)
	FindUnit(ctx context.Context, filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError)
	GetIPList(ctx context.Context, branchIfSpecific string, category string) ([]string, rest_err.APIError)
}
