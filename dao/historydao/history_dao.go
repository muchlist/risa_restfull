package historydao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout = 3
	keyHistColl    = "history"

	keyHistID                  = "_id"
	keyHistCreatedAt           = "created_at"
	keyHistCreatedByID         = "created_by_id"
	keyHistUpdatedAt           = "updated_at"
	keyHistUpdatedBy           = "updated_by"
	keyHistUpdatedByID         = "updated_by_id"
	keyHistBranch              = "branch"
	keyHistCategory            = "category"
	keyHistParentID            = "parent_id"
	keyHistStatus              = "status"
	keyHistProblem             = "problem"
	keyHistProblemResolve      = "problem_resolve"
	keyHistProblemLower        = "problem_lower"
	keyHistProblemResolveLower = "problem_resolve_lower"
	keyHistCompleteStatus      = "complete_status"
	keyHistDateStart           = "date_start"
	keyHistDateEnd             = "date_end"
	keyHistTag                 = "tag"
	keyHistImage               = "image"
	keyHistUpdates             = "updates"
)

func NewHistoryDao() HistoryDaoAssumer {
	return &historyDao{}
}

type historyDao struct {
}

type HistoryDaoAssumer interface {
	InsertHistory(input dto.History, isVendor bool) (*string, rest_err.APIError)
	InsertManyHistory(dataList []dto.History, isVendor bool) (int, rest_err.APIError)
	EditHistory(historyID primitive.ObjectID, input dto.HistoryEdit, isVendor bool) (*dto.HistoryResponse, rest_err.APIError)
	DeleteHistory(input dto.FilterIDBranchCreateGte) (*dto.HistoryResponse, rest_err.APIError)
	UploadImage(historyID primitive.ObjectID, imagePath string, filterBranch string) (*dto.HistoryResponse, rest_err.APIError)

	GetHistoryByID(historyID primitive.ObjectID, branchIfSpecific string) (*dto.HistoryResponse, rest_err.APIError)
	FindHistory(filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	SearchHistory(search string, filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForParent(parentID string) (dto.HistoryResponseMinList, rest_err.APIError)
	FindHistoryForUser(userID string, filter dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError)
	GetHistoryCount(branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError)
	FindHistoryForReport(branchIfSpecific string, start int64, end int64) (dto.HistoryResponseMinList, rest_err.APIError)
	UnwindHistory(filterA dto.FilterBranchCatInCompleteIn, filterB dto.FilterTimeRangeLimit) (dto.HistoryUnwindResponseList, rest_err.APIError)
}

func (h *historyDao) InsertHistory(input dto.History, isVendor bool) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Branch = strings.ToUpper(input.Branch)
	input.Category = strings.ToUpper(input.Category)
	if input.Tag == nil {
		input.Tag = []string{}
	}
	if input.Updates == nil {
		input.Updates = []dto.HistoryUpdate{}
	}

	// History versi 2 akan menambahkan detail riwayat perubahan dalam bentuk array
	input.Version = 2
	input.Updates = []dto.HistoryUpdate{{
		Time:           input.CreatedAt,
		UpdatedBy:      input.UpdatedBy,
		UpdatedByID:    input.UpdatedByID,
		Problem:        input.Problem,
		ProblemResolve: input.ProblemResolve,
		CompleteStatus: input.CompleteStatus,
		Vendor:         isVendor,
	},
	}

	// History problem dan problem resolve lower (digunakan untuk pencarian text)
	input.ProblemLower = strings.ToLower(input.Problem)
	input.ProblemResolveLower = strings.ToLower(input.ProblemResolve)

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan history ke database", err)
		logger.Error("Gagal menyimpan history ke database, (InsertHistory)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (h *historyDao) InsertManyHistory(dataList []dto.History, isVendor bool) (int, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	var dataForInserts []interface{}
	for _, data := range dataList {
		data.Branch = strings.ToUpper(data.Branch)
		data.Category = strings.ToUpper(data.Category)
		if data.Tag == nil {
			data.Tag = []string{}
		}
		if data.Updates == nil {
			data.Updates = []dto.HistoryUpdate{}
		}
		// History problem dan problem resolve lower (digunakan untuk pencarian text)
		data.ProblemLower = strings.ToLower(data.Problem)
		data.ProblemResolveLower = strings.ToLower(data.ProblemResolve)

		data.Updates = []dto.HistoryUpdate{{
			Time:           data.CreatedAt,
			UpdatedBy:      data.UpdatedBy,
			UpdatedByID:    data.UpdatedByID,
			Problem:        data.Problem,
			ProblemResolve: data.ProblemResolve,
			CompleteStatus: data.CompleteStatus,
			Vendor:         isVendor,
		},
		}
		dataForInserts = append(dataForInserts, data)
	}

	if len(dataForInserts) == 0 {
		return 0, nil
	}

	result, err := coll.InsertMany(ctx, dataForInserts)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan banyak history ke database", err)
		logger.Error("Gagal menyimpan banyak history ke database, (InsertManyHistory)", err)
		return 0, apiErr
	}

	totalInserted := len(result.InsertedIDs)

	return totalInserted, nil
}

func (h *historyDao) EditHistory(historyID primitive.ObjectID, input dto.HistoryEdit, isVendor bool) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	if input.Tag == nil {
		input.Tag = []string{}
	}

	// History problem dan problem resolve lower (digunakan untuk pencarian text)
	input.ProblemLower = strings.ToLower(input.Problem)
	input.ProblemResolveLower = strings.ToLower(input.ProblemResolve)

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
			keyHistUpdatedAt:           input.UpdatedAt,
			keyHistUpdatedBy:           input.UpdatedBy,
			keyHistUpdatedByID:         input.UpdatedByID,
			keyHistStatus:              input.Status,
			keyHistProblem:             input.Problem,
			keyHistProblemResolve:      input.ProblemResolve,
			keyHistProblemLower:        input.ProblemLower,
			keyHistProblemResolveLower: input.ProblemResolveLower,
			keyHistCompleteStatus:      input.CompleteStatus,
			keyHistDateEnd:             input.DateEnd,
			keyHistTag:                 input.Tag,
		},
		"$push": bson.M{
			keyHistUpdates: dto.HistoryUpdate{
				Time:           input.UpdatedAt,
				UpdatedBy:      input.UpdatedBy,
				UpdatedByID:    input.UpdatedByID,
				Problem:        input.Problem,
				ProblemResolve: input.ProblemResolve,
				CompleteStatus: input.CompleteStatus,
				Vendor:         isVendor,
			},
		},
	}

	var history dto.HistoryResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&history); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("History tidak diupdate : validasi id branch timestamp status_complete")
		}

		logger.Error("Gagal mendapatkan history dari database (EditHistory)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan history dari database", err)
		return nil, apiErr
	}

	return &history, nil
}

func (h *historyDao) DeleteHistory(input dto.FilterIDBranchCreateGte) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("History tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus history dari database (DeleteHistory)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan history dari database", err)
		return nil, apiErr
	}

	return &history, nil
}

func (h *historyDao) GetHistoryByID(historyID primitive.ObjectID, branchIfSpecific string) (*dto.HistoryResponse, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
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
		if errors.Is(err, mongo.ErrNoDocuments) {
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
	coll := db.DB.Collection(keyHistColl)
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
	// branch
	if filterA.FilterBranch != "" {
		filter[keyHistBranch] = filterA.FilterBranch
	}

	// category
	if filterA.FilterCategory != "" {
		// cek kategori jika multi category (pisah dengan koma)
		if strings.Contains(filterA.FilterCategory, ",") {
			categories := strings.Split(filterA.FilterCategory, ",")
			filter[keyHistCategory] = bson.M{"$in": categories}
		} else {
			filter[keyHistCategory] = filterA.FilterCategory
		}
	}

	// complete status
	if filterA.FilterCompleteStatus != nil && len(filterA.FilterCompleteStatus) != 0 {
		if len(filterA.FilterCompleteStatus) == 1 {
			filter[keyHistCompleteStatus] = filterA.FilterCompleteStatus[0]
		} else {
			filter[keyHistCompleteStatus] = bson.M{"$in": filterA.FilterCompleteStatus}
		}
	}

	// option range
	if filterB.FilterStart != 0 {
		filter[keyHistUpdatedAt] = bson.M{"$gte": filterB.FilterStart}
	}
	if filterB.FilterEnd != 0 {
		filter[keyHistCreatedAt] = bson.M{"$lte": filterB.FilterEnd}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyHistUpdatedAt, Value: -1}}) //nolint:govet
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

//db.history.createIndex(
//  {
//    problem_lower: "text",
//    problem_resolve_lower: "text"
//  },
//  {
//    weights: {
//      problem_lower: 5,
//      problem_resolve_lower: 3
//    }
//  }
//)

func (h *historyDao) SearchHistory(search string, filterA dto.FilterBranchCatComplete, filterB dto.FilterTimeRangeLimit) (dto.HistoryResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	search = strings.ToLower(search)

	// validate search
	if len(search) == 0 {
		return dto.HistoryResponseMinList{}, nil
	}

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterCategory = strings.ToUpper(filterA.FilterCategory)

	// set default limit
	if filterB.Limit == 0 {
		filterB.Limit = 100
	}

	// empty filter
	filter := bson.M{}

	filter["$text"] = bson.M{"$search": search}

	// filter condition
	// branch
	if filterA.FilterBranch != "" {
		filter[keyHistBranch] = filterA.FilterBranch
	}

	// category
	if filterA.FilterCategory != "" {
		// cek kategori jika multi category (pisah dengan koma)
		if strings.Contains(filterA.FilterCategory, ",") {
			categories := strings.Split(filterA.FilterCategory, ",")
			filter[keyHistCategory] = bson.M{"$in": categories}
		} else {
			filter[keyHistCategory] = filterA.FilterCategory
		}
	}

	// complete status
	if filterA.FilterCompleteStatus != nil && len(filterA.FilterCompleteStatus) != 0 {
		if len(filterA.FilterCompleteStatus) == 1 {
			filter[keyHistCompleteStatus] = filterA.FilterCompleteStatus[0]
		} else {
			filter[keyHistCompleteStatus] = bson.M{"$in": filterA.FilterCompleteStatus}
		}
	}

	// option range
	if filterB.FilterStart != 0 {
		filter[keyHistUpdatedAt] = bson.M{"$gte": filterB.FilterStart}
	}
	if filterB.FilterEnd != 0 {
		filter[keyHistCreatedAt] = bson.M{"$lte": filterB.FilterEnd}
	}

	opts := options.Find()
	opts.SetSort(bson.M{
		"score": bson.M{
			"$meta": "textScore",
		},
	})
	opts.SetLimit(filterB.Limit)

	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (SearchHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories := dto.HistoryResponseMinList{}
	if err = cursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (SearchHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	return histories, nil
}

// UnwindHistory mengembalikan unwind history dengan urutan
//{Key: "$sort", Value: bson.D{{Key: keyHistUpdatedAt, Value: -1}, {Key: "updates.time", Value: 1}}},
func (h *historyDao) UnwindHistory(filterA dto.FilterBranchCatInCompleteIn, filterB dto.FilterTimeRangeLimit) (dto.HistoryUnwindResponseList, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
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
	// branch
	if filterA.FilterBranch != "" {
		filter[keyHistBranch] = filterA.FilterBranch
	}

	// category
	if filterA.FilterCategory != "" {
		// cek kategori jika multi category (pisah dengan koma)
		if strings.Contains(filterA.FilterCategory, ",") {
			categories := strings.Split(filterA.FilterCategory, ",")
			filter[keyHistCategory] = bson.M{"$in": categories}
		} else {
			filter[keyHistCategory] = filterA.FilterCategory
		}
	}

	// complete status
	if filterA.FilterCompleteStatus != "" {
		// cek complete status jika multi status (pisah dengan koma)
		if strings.Contains(filterA.FilterCompleteStatus, ",") {
			statusStr := strings.Split(filterA.FilterCompleteStatus, ",")
			var statusInt []int
			for _, status := range statusStr {
				statusConverted := sfunc.StrToInt(strings.Trim(status, " "), -1)
				if statusConverted == -1 {
					continue
				}
				statusInt = append(statusInt, sfunc.StrToInt(status, statusConverted))
			}
			filter[keyHistCompleteStatus] = bson.M{"$in": statusInt}
		} else {
			statusConverted := sfunc.StrToInt(strings.Trim(filterA.FilterCompleteStatus, " "), -1)
			if statusConverted != -1 {
				filter[keyHistCompleteStatus] = statusConverted
			}
		}
	}

	// option range
	if filterB.FilterStart != 0 {
		filter[keyHistUpdatedAt] = bson.M{"$gte": filterB.FilterStart}
	}
	if filterB.FilterEnd != 0 {
		filter[keyHistCreatedAt] = bson.M{"$lte": filterB.FilterEnd}
	}

	groupStage := bson.D{
		{Key: "$match", Value: filter},
	}
	unwindStage := bson.D{
		{Key: "$unwind", Value: "$updates"},
	}
	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{{Key: keyHistUpdatedAt, Value: -1}, {Key: "updates.time", Value: 1}}},
	}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{groupStage, unwindStage, sortStage})

	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (UnwindHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryUnwindResponseList{}, apiErr
	}

	histories := dto.HistoryUnwindResponseList{}
	if err = cursor.All(ctx, &histories); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (UnwindHistory)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryUnwindResponseList{}, apiErr
	}

	return histories, nil
}

func (h *historyDao) FindHistoryForParent(parentID string) (dto.HistoryResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyHistParentID: parentID,
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyHistID, Value: -1}}) //nolint:govet
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
	coll := db.DB.Collection(keyHistColl)
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
	opts.SetSort(bson.D{{Key: keyHistID, Value: -1}}) //nolint:govet
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

// get_histories_in_progress_count
func (h *historyDao) GetHistoryCount(branchIfSpecific string, statusComplete int) (dto.HistoryCountList, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
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
		{Key: "$match", Value: filter},
	}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$branch"},
			{Key: "count", Value: bson.M{"$sum": 1}},
		}},
	}
	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "count", Value: -1},
			{Key: "_id", Value: -1},
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
	coll := db.DB.Collection(keyHistColl)
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, history dengan id %s tidak ditemukan", historyID.Hex()))
		}

		logger.Error("Memasukkan path image history ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image history ke db gagal", err)
		return nil, apiErr
	}

	return &history, nil
}

func (h *historyDao) FindHistoryForReport(branchIfSpecific string, start int64, end int64) (dto.HistoryResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyHistColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()
	branch := strings.ToUpper(branchIfSpecific)

	// kenapa dipisah ?
	// karena yang dimunculkan adalah yang pending dan progress sebelum waktu end, sedangkan complete dan info sesuai range waktu
	// Find complete (4) and info (0)
	filter := bson.M{
		keyHistBranch:         branch,
		keyHistCompleteStatus: bson.M{"$in": bson.A{0, 4}},
		keyHistUpdatedAt:      bson.M{"$gte": start, "$lte": end},
	}
	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyHistUpdatedAt, Value: -1}}) //nolint:govet

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (FindHistoryForReport cursor)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories04 := dto.HistoryResponseMinList{}
	if err = cursor.All(ctx, &histories04); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (FindHistoryForReport 04)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	// Find progress (1) and pending (2, 3)
	filter = bson.M{
		keyHistBranch:         branch,
		keyHistCompleteStatus: bson.M{"$in": bson.A{1, 2, 3}},
		keyHistUpdatedAt:      bson.M{"$lte": end},
	}

	cursor, err = coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar history dari database (FindHistoryForReport cursor2)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories123 := dto.HistoryResponseMinList{}
	if err = cursor.All(ctx, &histories123); err != nil {
		logger.Error("Gagal decode histories cursor ke objek slice (FindHistoryForReport 123)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.HistoryResponseMinList{}, apiErr
	}

	histories := append(histories04, histories123...)
	return histories, nil
}
