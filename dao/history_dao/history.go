package history_dao

import (
	"context"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	connectTimeout = 3
	keyHistColl    = "history"

	keyHistID             = "_id"
	keyHistCreatedAt      = "created_at"
	keyHistCreatedBy      = "created_by"
	keyHistCreatedByID    = "created_by_id"
	keyHistUpdatedAt      = "updated_at"
	keyHistUpdatedBy      = "updated_by"
	keyHistUpdatedByID    = "updated_by_id"
	keyHistBranch         = "branch"
	keyHistCategory       = "category"
	keyHistParentID       = "parent_id"
	keyHistParentName     = "parent_name"
	keyHistStatus         = "status"
	keyHistProblem        = "problem"
	keyHistProblemResolve = "problem_resolve"
	keyHistCompleteStatus = "complete_status"
	keyHistDateStart      = "date_start"
	keyHistDateEnd        = "date_end"
	keyHistTag            = "tag"
)

func NewHistoryDao() HistoryDaoAssumer {
	return &historyDao{}
}

type historyDao struct {
}

type HistoryDaoAssumer interface {
	InsertHistory(input dto.History) (*string, rest_err.APIError)
	EditHistory(historyID primitive.ObjectID, input dto.HistoryEdit) (*dto.HistoryResponse, rest_err.APIError)
	//DeleteUnit(unitID string) rest_err.APIError
	//InsertCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	//DeleteCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	////insertPing
	//
	//GetUnitByID(unitID string) (*dto.GenUnitResponse, rest_err.APIError)
	//FindUnit(filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError)
}

func (h *historyDao) InsertHistory(input dto.History) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Branch = strings.ToUpper(input.Branch)
	input.Category = strings.ToUpper(input.Category)

	insertDoc := bson.M{
		keyHistID:             input.ID,
		keyHistCreatedAt:      input.CreatedAt,
		keyHistCreatedBy:      input.CreatedBy,
		keyHistCreatedByID:    input.CreatedByID,
		keyHistUpdatedAt:      input.UpdatedAt,
		keyHistUpdatedBy:      input.UpdatedBy,
		keyHistUpdatedByID:    input.UpdatedByID,
		keyHistBranch:         input.Branch,
		keyHistCategory:       input.Category,
		keyHistParentID:       input.ParentID,
		keyHistParentName:     input.ParentName,
		keyHistStatus:         input.Status,
		keyHistProblem:        input.Problem,
		keyHistProblemResolve: input.ProblemResolve,
		keyHistCompleteStatus: input.CompleteStatus,
		keyHistDateStart:      input.DateStart,
		keyHistDateEnd:        input.DateEnd,
		keyHistTag:            input.Tag,
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan history ke database", err)
		logger.Error("Gagal menyimpan history ke database, (InsertHistory)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (h *historyDao) EditHistory(historyID primitive.ObjectID, input dto.HistoryEdit) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyHistID:             historyID,
		keyHistBranch:         input.FilterBranch,
		keyHistUpdatedAt:      input.FilterTimestamp,
		keyHistCompleteStatus: bson.M{"$ne": enum.HComplete},
	}

	update := bson.M{
		"$set": bson.M{
			keyHistUpdatedAt:      input.UpdatedAt,
			keyHistUpdatedBy:      input.UpdatedBy,
			keyHistUpdatedByID:    input.UpdatedByID,
			keyHistStatus:         input.Status,
			keyHistProblem:        input.Problem,
			keyHistProblemResolve: input.ProblemResolve,
			keyHistCompleteStatus: input.CompleteStatus,
			keyHistDateEnd:        input.DateEnd,
			keyHistTag:            input.Tag,
		},
	}

	var history dto.HistoryResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&history); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("History tidak diupdate karena ID atau timestamp tidak valid")
		}

		logger.Error("Gagal mendapatkan history dari database (EditHistory)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan history dari database", err)
		return nil, apiErr
	}

	return &history, nil
}
