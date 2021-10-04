package altaicheckdao

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
	keyCollection  = "altaiCheck"
	keyID          = "_id"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
	keyUpdatedBy   = "updated_by"
	keyUpdatedByID = "updated_by_id"
	keyBranch      = "branch"

	keyTimeStarted     = "time_started"
	keyTimeEnded       = "time_ended"
	keyIsFinish        = "is_finish"
	keyAltaiCheckItems = "altai_check_items"
	keyNote            = "note"

	keyChXId        = "altai_check_items.id"
	keyChXCheckedAt = "altai_check_items.$.checked_at"
	keyChXCheckedBy = "altai_check_items.$.checked_by"
	keyChXIsChecked = "altai_check_items.$.is_checked"
	keyChXIsOffline = "altai_check_items.$.is_offline"
	keyChXImagePath = "altai_check_items.$.image_path"
)

func NewAltaiCheckDao() CheckAltaiDaoAssumer {
	return &checkAltaiDao{}
}

type checkAltaiDao struct{}

func (c *checkAltaiDao) InsertCheck(ctx context.Context, input dto.AltaiCheck) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// Default value for slice
	input.Branch = strings.ToUpper(input.Branch)
	if input.AltaiCheckItems == nil {
		input.AltaiCheckItems = []dto.AltaiCheckItemEmbed{}
	}

	result, err := coll.InsertOne(ctxt, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan altai check ke database", err)
		logger.Error("Gagal menyimpan altai check ke database, (AltaiInsertCheck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *checkAltaiDao) EditCheck(ctx context.Context, input dto.AltaiCheckEdit) (*dto.AltaiCheck, rest_err.APIError) {
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

	var check dto.AltaiCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("altai check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan altai check dari database (EditCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan altai check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiDao) DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.AltaiCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyID:        input.FilterID,
		keyBranch:    input.FilterBranch,
		keyCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var check dto.AltaiCheck
	err := coll.FindOneAndDelete(ctxt, filter).Decode(&check)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Altai Check tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus cek altai dari database (AltaiDeleteCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus altai check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiDao) UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.AltaiCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:    filterA.FilterParentID,
		keyChXId: filterA.FilterChildID,
	}

	update := bson.M{
		"$set": bson.M{
			keyChXImagePath: imagePath,
		},
	}

	var check dto.AltaiCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, altai check dengan id %s -> %s tidak ditemukan", filterA.FilterParentID.Hex(), filterA.FilterChildID))
		}

		logger.Error("Memasukkan path image child altai Check ke db gagal, (AltaiUploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image altai check ke db gagal", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiDao) UpdateCheckItem(ctx context.Context, input dto.AltaiCheckItemUpdate) (*dto.AltaiCheck, rest_err.APIError) {
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
			keyChXIsChecked: input.IsChecked,
			keyChXIsOffline: input.IsOffline,
		},
	}

	var check dto.AltaiCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("altai check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (AltaiUpdateCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan altai check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiDao) BulkUpdateItem(ctx context.Context, inputs []dto.AltaiCheckItemUpdate) (int64, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	if len(inputs) == 0 {
		return 0, rest_err.NewBadRequestError("input tidak boleh kosong")
	}

	operations := make([]mongo.WriteModel, len(inputs))
	for i, input := range inputs {
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
				keyChXIsChecked: input.IsChecked,
				keyChXIsOffline: input.IsOffline,
			},
		}
		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := coll.BulkWrite(ctxt, operations, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, rest_err.NewBadRequestError("altai check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("gagal bulk write checkItem dari database (BulkUpdateItem)", err)
		apiErr := rest_err.NewInternalServerError("gagal bulk write altai check dari database", err)
		return 0, apiErr
	}

	return result.ModifiedCount, nil
}

func (c *checkAltaiDao) GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.AltaiCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: checkID}
	// filter condition
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var check dto.AltaiCheck
	if err := coll.FindOne(ctxt, filter).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("altai check dengan ID %s tidak ditemukan. validation : id branch", checkID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan altai check dari database (GetCheckByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan altai check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiDao) FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.AltaiCheck, rest_err.APIError) {
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
			keyAltaiCheckItems: 0,
		})
	}
	opts.SetSort(bson.D{{Key: keyUpdatedAt, Value: -1}}) //nolint:govet
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar altai check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.AltaiCheck{}, apiErr
	}

	var checkList []dto.AltaiCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (AltaiFindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.AltaiCheck{}, apiErr
	}

	return checkList, nil
}

func (c *checkAltaiDao) GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string) (*dto.AltaiCheck, rest_err.APIError) {
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
		logger.Error("gagal mendapatkan daftar altai check dari database (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var checkList []dto.AltaiCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var check dto.AltaiCheck
	if len(checkList) != 0 {
		check = checkList[0]
	}

	return &check, nil
}
