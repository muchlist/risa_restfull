package historydao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistoryDaoAssumer interface {
	HistorySaver
	HistoryLoader
}

type HistorySaver interface {
	InsertHistory(ctx context.Context, input dto.History, isVendor bool) (*string, rest_err.APIError)
	InsertManyHistory(ctx context.Context, dataList []dto.History, isVendor bool) (int, rest_err.APIError)
	EditHistory(ctx context.Context, historyID primitive.ObjectID, input dto.HistoryEdit, isVendor bool) (*dto.HistoryResponse, rest_err.APIError)
	DeleteHistory(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.HistoryResponse, rest_err.APIError)
	UploadImage(ctx context.Context, historyID primitive.ObjectID, imagePath string, filterBranch string) (*dto.HistoryResponse, rest_err.APIError)
}

type HistoryLoader interface {
	GetHistoryByID(ctx context.Context, historyID primitive.ObjectID, branchIfSpecific string) (*dto.HistoryResponse, rest_err.APIError)
	FindHistory(ctx context.Context, filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	SearchHistory(ctx context.Context, search string, filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForParent(ctx context.Context, parentID string) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForUser(ctx context.Context, userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	GetHistoryCount(ctx context.Context, branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError)
	FindHistoryForReport(ctx context.Context, branchIfSpecific string, start int64, end int64) (dto.HistoryResponseMinList, rest_err.APIError)
	UnwindHistory(ctx context.Context, filterA dto.FilterBranchCatInCompleteIn, filterB dto.FilterTimeRangeLimit) (dto.HistoryUnwindResponseList, rest_err.APIError)
}
