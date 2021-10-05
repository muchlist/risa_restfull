package configcheckdao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout = 3
	keyCollection  = "configCheck"
	keyID          = "_id"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
	keyUpdatedBy   = "updated_by"
	keyUpdatedByID = "updated_by_id"
	keyBranch      = "branch"

	keyTimeStarted      = "time_started"
	keyTimeEnded        = "time_ended"
	keyIsFinish         = "is_finish"
	keyConfigCheckItems = "config_check_items"
	keyNote             = "note"

	keyChXId            = "config_check_items.id"
	keyChXCheckedAt     = "config_check_items.$.checked_at"
	keyChXCheckedBy     = "config_check_items.$.checked_by"
	keyChXIsUpdated     = "config_check_items.$.is_updated"
	keyChXElemCheckedAt = "config_check_items.$[elem].checked_at"
	keyChXElemCheckedBy = "config_check_items.$[elem].checked_by"
	keyChXElemIsUpdated = "config_check_items.$[elem].is_updated"
)

func NewConfigCheckDao() CheckConfigDaoAssumer {
	return &checkConfigDao{}
}

type checkConfigDao struct {
}

func (c *checkConfigDao) InsertCheck(ctx context.Context, input dto.ConfigCheck) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// Default value for slice
	input.Branch = strings.ToUpper(input.Branch)
	if input.ConfigCheckItems == nil {
		input.ConfigCheckItems = []dto.ConfigCheckItemEmbed{}
	}

	result, err := coll.InsertOne(ctxt, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan config check ke database", err)
		logger.Error("Gagal menyimpan config check ke database, (ConfigInsertCheck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *checkConfigDao) EditCheck(ctx context.Context, input dto.ConfigCheckEdit) (*dto.ConfigCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// Default value for slice

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:       input.FilterID,
		keyBranch:   input.FilterBranch,
		keyIsFinish: false,
	}

	update := bson.M{
		"$set": bson.M{
			keyUpdatedAt:   input.UpdatedAt,
			keyUpdatedBy:   input.UpdatedBy,
			keyUpdatedByID: input.UpdatedByID,
			keyTimeStarted: input.TimeStarted,
			keyTimeEnded:   input.TimeEnded,
			keyIsFinish:    input.IsFinish,
			keyNote:        input.Note,
		},
	}

	var check dto.ConfigCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("config check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan config check dari database (EditCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan config check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkConfigDao) DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.ConfigCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyID:        input.FilterID,
		keyBranch:    input.FilterBranch,
		keyCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var check dto.ConfigCheck
	err := coll.FindOneAndDelete(ctxt, filter).Decode(&check)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Config Check tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus cek config dari database (ConfigDeleteCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus config check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkConfigDao) UpdateCheckItem(ctx context.Context, input dto.ConfigCheckItemUpdate) (*dto.ConfigCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:       input.FilterParentID,
		keyChXId:    input.FilterChildID,
		keyBranch:   strings.ToUpper(input.FilterBranch),
		keyIsFinish: false,
	}

	update := bson.M{
		"$set": bson.M{
			keyUpdatedAt:    input.CheckedAt,
			keyChXCheckedAt: input.CheckedAt,
			keyChXCheckedBy: input.CheckedBy,
			keyChXIsUpdated: input.IsUpdated,
		},
	}

	var check dto.ConfigCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("config check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (ConfigUpdateCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan config check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

// UpdateManyItem mengupdate list didalam checkConfigDetail
func (c *checkConfigDao) UpdateManyItem(ctx context.Context, input dto.ConfigCheckUpdateMany) rest_err.APIError {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	if len(input.ChildIDsUpdate) == 0 || input.ChildIDsUpdate == nil {
		return nil
	}

	filter := bson.M{
		keyID: input.ParentID,
		keyConfigCheckItems: bson.M{
			"$elemMatch": bson.M{
				"id": bson.M{"$in": input.ChildIDsUpdate},
			},
		},
		keyBranch:   strings.ToUpper(input.Branch),
		keyIsFinish: false,
	}

	update := bson.M{
		"$set": bson.M{
			keyChXElemIsUpdated: input.UpdatedValue,
			keyChXElemCheckedAt: time.Now().Unix(),
			keyChXElemCheckedBy: input.Updater,
		},
	}

	opts := options.Update()
	opts.SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.id": bson.M{"$in": input.ChildIDsUpdate}}},
	})

	_, err := coll.UpdateOne(ctxt, filter, update, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return rest_err.NewBadRequestError("config check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (UpdateManyItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan config check dari database", err)
		return apiErr
	}

	return nil
}

func (c *checkConfigDao) GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.ConfigCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: checkID}
	// filter condition
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var check dto.ConfigCheck
	if err := coll.FindOne(ctxt, filter).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("config check dengan ID %s tidak ditemukan. validation : id branch", checkID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan config check dari database (GetCheckByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan config check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkConfigDao) FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.ConfigCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	if filterA.Limit == 0 {
		filterA.Limit = 100
	}

	// filter
	filter := bson.M{}
	// filter condition
	if branch != "" {
		filter[keyBranch] = strings.ToUpper(branch)
	}
	if filterA.FilterStart != 0 {
		filter[keyUpdatedAt] = bson.M{"$gte": filterA.FilterStart}
	}
	if filterA.FilterEnd != 0 {
		filter[keyCreatedAt] = bson.M{"$lte": filterA.FilterEnd}
	}

	opts := options.Find()

	if !detail {
		opts.SetProjection(bson.M{
			keyConfigCheckItems: 0,
		})
	}
	opts.SetSort(bson.D{{Key: keyUpdatedAt, Value: -1}})
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar config check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.ConfigCheck{}, apiErr
	}

	var checkList []dto.ConfigCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (ConfigFindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.ConfigCheck{}, apiErr
	}

	return checkList, nil
}

func (c *checkConfigDao) GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string) (*dto.ConfigCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{}
	// filter condition
	if branch != "" {
		filter[keyBranch] = strings.ToUpper(branch)
	}
	filter[keyCreatedAt] = bson.M{"$gte": start}
	filter[keyCreatedAt] = bson.M{"$lte": end}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyCreatedAt, Value: -1}})
	opts.SetLimit(1)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar config check dari database (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var checkList []dto.ConfigCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var check dto.ConfigCheck
	if len(checkList) != 0 {
		check = checkList[0]
	}

	return &check, nil
}
