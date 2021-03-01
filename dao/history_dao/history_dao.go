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
	keyHistCreatedByID    = "created_by_id"
	keyHistUpdatedAt      = "updated_at"
	keyHistUpdatedBy      = "updated_by"
	keyHistUpdatedByID    = "updated_by_id"
	keyHistBranch         = "branch"
	keyHistCategory       = "category"
	keyHistParentID       = "parent_id"
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
	DeleteHistory(input dto.FilterIDBranchCreateGte) (*dto.HistoryResponse, rest_err.APIError)
	UploadImage(historyID primitive.ObjectID, imagePath string, filterBranch string) (*dto.HistoryResponse, rest_err.APIError)

	GetHistoryByID(historyID primitive.ObjectID, branchIfSpecific string) (*dto.HistoryResponse, rest_err.APIError)
	FindHistory(filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForParent(parentID string) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForUser(userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	GetHistoryCount(branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError)
}

func (h *historyDao) InsertHistory(input dto.History) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Branch = strings.ToUpper(input.Branch)
	input.Category = strings.ToUpper(input.Category)
	if input.Tag == nil {
		input.Tag = []string{}
	}

	result, err := coll.InsertOne(ctx, input)
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

	if input.Tag == nil {
		input.Tag = []string{}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyHistID:             historyID,
		keyHistBranch:         input.FilterBranch,
		keyHistUpdatedAt:      input.FilterTimestamp,
		keyHistCompleteStatus: bson.M{"$nin": bson.A{enum.HComplete, enum.HInfo}},
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
			return nil, rest_err.NewBadRequestError("History tidak diupdate : validasi id branch timestamp status_complete")
		}

		logger.Error("Gagal mendapatkan history dari database (EditHistory)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan history dari database", err)
		return nil, apiErr
	}

	return &history, nil
}

func (h *historyDao) DeleteHistory(input dto.FilterIDBranchCreateGte) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyHistID:        input.FilterID,
		keyHistBranch:    input.FilterBranch,
		keyHistCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var history dto.HistoryResponse
	err := coll.FindOneAndDelete(ctx, filter).Decode(&history)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("History tidak diupdate : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus history dari database (DeleteHistory)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan history dari database", err)
		return nil, apiErr
	}

	return &history, nil
}

func (h *historyDao) GetHistoryByID(historyID primitive.ObjectID, branchIfSpecific string) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{keyHistID: historyID}

	// filter conditional
	if branchIfSpecific != "" {
		filter[keyHistBranch] = strings.ToUpper(branchIfSpecific)
	}
	var history dto.HistoryResponse
	if err := coll.FindOne(ctx, filter).Decode(&history); err != nil {

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

func (h *historyDao) FindHistory(filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterCategory = strings.ToUpper(filterA.FilterCategory)

	// set default limit
	if filterB.Limit == 0 {
		filterB.Limit = 100
	}

	// empty filter
	filter := bson.M{}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyHistBranch] = filterA.FilterBranch
	}
	if filterA.FilterCategory != "" {
		filter[keyHistCategory] = filterA.FilterCategory
	}
	if filterA.FilterCompleteStatus != 0 {
		filter[keyHistCompleteStatus] = filterA.FilterCompleteStatus
	}

	// option range
	if filterB.FilterStart != 0 {
		filter[keyHistDateStart] = bson.M{"$gte": filterB.FilterStart}
	}
	if filterB.FilterEnd != 0 {
		filter[keyHistDateEnd] = bson.M{"$lte": filterB.FilterStart}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyHistUpdatedAt, -1}})
	opts.SetLimit(filterB.Limit)

	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (FindHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories := dto.HistoryResponseMinList{}
	if err = cursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (FindHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	return histories, nil
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
	if filterOpt.FilterStart != 0 {
		filter[keyHistDateStart] = bson.M{"$gte": filterOpt.FilterStart}
	}
	if filterOpt.FilterEnd != 0 {
		filter[keyHistDateEnd] = bson.M{"$lte": filterOpt.FilterStart}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyHistID, -1}})
	opts.SetLimit(filterOpt.Limit)

	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (FindHistoryForUser)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories := dto.HistoryResponseMinList{}
	if err = cursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (FindHistoryForuser)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	return histories, nil
}

//get_histories_in_progress_count
func (h *historyDao) GetHistoryCount(branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// Jika branch ada isinya maka hanya menampilkan branch tersebut,
	// umumnya digunakan dengan branch kosong untuk melihat semua cabang
	branchIfSpecific = strings.ToUpper(branchIfSpecific)

	filter := bson.M{
		keyHistCompleteStatus: statusComplete,
	}
	if branchIfSpecific != "" {
		filter[keyHistBranch] = branchIfSpecific
	}

	matchStage := bson.D{
		{"$match", filter},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$branch"},
			{"count", bson.M{"$sum": 1}},
		}},
	}
	sortStage := bson.D{
		{"$sort", bson.D{
			{"count", -1},
			{"_id", -1},
		}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage})

	if err != nil {
		logger.Error("Gagal mendapatkan history count dari database (GetHistoryCount)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryCountList{}, apiErr
	}

	histories := dto.HistoryCountList{}
	if err = cursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode history count ke objek slice (GetHistoryCount)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryCountList{}, apiErr
	}

	return histories, nil
}

// UploadImage tidak digunakan saat pembuatan history dengan langsung
// menyertakan image, hanya untuk keperluan update pada dokumen yang sudah ada
func (h *historyDao) UploadImage(historyID primitive.ObjectID, imagePath string, filterBranch string) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyHistID:     historyID,
		keyHistBranch: strings.ToUpper(filterBranch),
	}
	update := bson.M{
		"$set": bson.M{
			keyHistImage: imagePath,
		},
	}

	var history dto.HistoryResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&history); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, history dengan id %s tidak ditemukan", historyID.Hex()))
		}

		logger.Error("Memasukkan path image history ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image history ke db gagal", err)
		return nil, apiErr
	}

	return &history, nil
}