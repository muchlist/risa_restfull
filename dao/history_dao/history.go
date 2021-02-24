package history_dao

import (
	"context"
	"fmt"
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
	keyHistImage          = "image"
)

func NewHistoryDao() HistoryDaoAssumer {
	return &historyDao{}
}

type historyDao struct {
}

type HistoryDaoAssumer interface {
	InsertHistory(input dto.History) (*string, rest_err.APIError)
	EditHistory(historyID primitive.ObjectID, input dto.HistoryEdit) (*dto.HistoryResponse, rest_err.APIError)
	DeleteHistory(input dto.FilterIDBranchTime) rest_err.APIError

	GetHistoryByID(historyID primitive.ObjectID) (*dto.HistoryResponse, rest_err.APIError)
	FindHistoryForParent(parentID string) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForUser(userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistory(filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	//InsertCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	//DeleteCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	////insertPing
	//
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
		keyHistImage:          input.Image,
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

func (h *historyDao) DeleteHistory(input dto.FilterIDBranchTime) rest_err.APIError {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyHistID:        input.ID,
		keyHistBranch:    input.Branch,
		keyHistCreatedAt: bson.M{"$gte": input.Time},
	}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest_err.NewBadRequestError("History gagal dihapus, limit waktu terlampaui, id atau cabang salah")
		}

		logger.Error("Gagal menghapus history dari database (DeleteHistory)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan history dari database", err)
		return apiErr
	}

	if result.DeletedCount == 0 {
		return rest_err.NewBadRequestError("History gagal dihapus, limit waktu terlampaui, id atau cabang salah")
	}

	return nil
}

func (h *historyDao) GetHistoryByID(historyID primitive.ObjectID) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	var history dto.HistoryResponse
	opts := options.FindOne()

	if err := coll.FindOne(ctx, bson.M{keyHistID: historyID}, opts).Decode(&history); err != nil {

		if err == mongo.ErrNoDocuments {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("History dengan ID %s tidak ditemukan", historyID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan history dari database (GetHistoryByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan history dari database", err)
		return nil, apiErr
	}

	return &history, nil
}

func (h *historyDao) FindHistoryForParent(parentID string) (dto.HistoryResponseMinList, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyHistParentID: parentID,
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyHistID, -1}})
	sortCursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan history dari database (FindHistoryForParent)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories := dto.HistoryResponseMinList{}
	if err = sortCursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (FindHistoryForParent)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	return histories, nil
}

func (h *historyDao) FindHistoryForUser(userID string, filterOpt dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// set default limit
	if filterOpt.Limit == 0 {
		filterOpt.Limit = 100
	}

	filter := bson.M{
		// menampilkan yang dibuat atau diupdate oleh UserID
		"$or": bson.A{
			bson.M{keyHistCreatedByID: userID},
			bson.M{keyHistUpdatedByID: userID},
		},
	}

	// option range
	if filterOpt.Start != 0 {
		filter[keyHistDateStart] = bson.M{"$gte": filterOpt.Start}
	}
	if filterOpt.End != 0 {
		filter[keyHistDateEnd] = bson.M{"$lte": filterOpt.Start}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyHistID, -1}})
	opts.SetLimit(filterOpt.Limit)

	Cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (FindHistoryForUser)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories := dto.HistoryResponseMinList{}
	if err = Cursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (FindHistoryForuser)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	return histories, nil
}

func (h *historyDao) FindHistory(filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.Branch = strings.ToUpper(filterA.Branch)
	filterA.Category = strings.ToUpper(filterA.Category)

	// set default limit
	if filterB.Limit == 0 {
		filterB.Limit = 100
	}

	// empty filter
	filter := bson.M{}

	// filter condition
	if filterA.Branch != "" {
		filter[keyHistBranch] = filterA.Branch
	}
	if filterA.Category != "" {
		filter[keyHistCategory] = filterA.Category
	}
	if filterA.CompleteStatus != 0 {
		filter[keyHistCompleteStatus] = filterA.CompleteStatus
	}

	// option range
	if filterB.Start != 0 {
		filter[keyHistDateStart] = bson.M{"$gte": filterB.Start}
	}
	if filterB.End != 0 {
		filter[keyHistDateEnd] = bson.M{"$lte": filterB.Start}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyHistID, -1}})
	opts.SetLimit(filterB.Limit)

	Cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (FindHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories := dto.HistoryResponseMinList{}
	if err = Cursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (FindHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	return histories, nil
}

//get_histories_in_progress_count
