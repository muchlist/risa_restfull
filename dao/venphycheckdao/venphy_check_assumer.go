package venphycheckdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckVenPhyDaoAssumer interface {
	CheckVenPhySaver
	CheckVenPhyLoader
}

type CheckVenPhySaver interface {
	InsertCheck(ctx context.Context, input dto.VenPhyCheck) (*string, rest_err.APIError)
	EditCheck(ctx context.Context, input dto.VenPhyCheckEdit) (*dto.VenPhyCheck, rest_err.APIError)
	DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.VenPhyCheck, rest_err.APIError)
	UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.VenPhyCheck, rest_err.APIError)
	UpdateCheckItem(ctx context.Context, input dto.VenPhyCheckItemUpdate) (*dto.VenPhyCheck, rest_err.APIError)
	BulkUpdateItem(ctx context.Context, inputs []dto.VenPhyCheckItemUpdate) (int64, rest_err.APIError)
	BulkUpdateItemForUpdateCheckItem(ctx context.Context, inputs []dto.VenPhyCheckItemUpdate) (int64, rest_err.APIError)
	OverwriteChecklist(ctx context.Context, id primitive.ObjectID, checkItems []dto.VenPhyCheckItemEmbed) (*dto.VenPhyCheck, rest_err.APIError)
	UndoFinishCheck(ctx context.Context, filterID primitive.ObjectID, filterBranch string) (*dto.VenPhyCheck, rest_err.APIError)
}

type CheckVenPhyLoader interface {
	GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.VenPhyCheck, rest_err.APIError)
	FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.VenPhyCheck, rest_err.APIError)
	FindCheckStillOpen(ctx context.Context, branch string, detail bool) ([]dto.VenPhyCheck, rest_err.APIError)
	GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string, isQuarter bool) (*dto.VenPhyCheck, rest_err.APIError)
}
