package altaiphycheckdao

import (
	"context"
	"errors"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
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
	keyCollection  = "altaiPhyCheck"
	keyID          = "_id"
	keyName        = "name"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
	keyUpdatedBy   = "updated_by"
	keyUpdatedByID = "updated_by_id"
	keyBranch      = "branch"

	keyTimeEnded          = "time_ended"
	keyIsFinish           = "is_finish"
	keyIsQuarterly        = "quarterly_mode"
	keyAltaiPhyCheckItems = "altai_phy_check_items"
	keyNote               = "note"

	keyChXId           = "altai_phy_check_items.id"
	keyChXCheckedAt    = "altai_phy_check_items.$.checked_at"
	keyChXCheckedBy    = "altai_phy_check_items.$.checked_by"
	keyChXIsChecked    = "altai_phy_check_items.$.is_checked"
	keyChXIsMaintained = "altai_phy_check_items.$.is_maintained"
	keyChXIsOffline    = "altai_phy_check_items.$.is_offline"
	keyChXImagePath    = "altai_phy_check_items.$.image_path"
)

func NewAltaiPhyCheckDao() CheckAltaiPhyDaoAssumer {
	return &checkAltaiPhyDao{}
}

type checkAltaiPhyDao struct {
}

type CheckAltaiPhyDaoAssumer interface {
	InsertCheck(input dto.AltaiPhyCheck) (*string, rest_err.APIError)
	EditCheck(input dto.AltaiPhyCheckEdit) (*dto.AltaiPhyCheck, rest_err.APIError)
	DeleteCheck(input dto.FilterIDBranchCreateGte) (*dto.AltaiPhyCheck, rest_err.APIError)
	UploadChildImage(filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.AltaiPhyCheck, rest_err.APIError)
	UpdateCheckItem(input dto.AltaiPhyCheckItemUpdate) (*dto.AltaiPhyCheck, rest_err.APIError)
	BulkUpdateItem(inputs []dto.AltaiPhyCheckItemUpdate) (int64, rest_err.APIError)

	GetCheckByID(checkID primitive.ObjectID, branchIfSpecific string) (*dto.AltaiPhyCheck, rest_err.APIError)
	FindCheck(branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.AltaiPhyCheck, rest_err.APIError)
	FindCheckStillOpen(branch string, detail bool) ([]dto.AltaiPhyCheck, rest_err.APIError)
	GetLastCheckCreateRange(start, end int64, branch string, isQuarter bool) (*dto.AltaiPhyCheck, rest_err.APIError)
}

func (c *checkAltaiPhyDao) InsertCheck(input dto.AltaiPhyCheck) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// Default value for slice
	input.Branch = strings.ToUpper(input.Branch)
	if input.AltaiPhyCheckItems == nil {
		input.AltaiPhyCheckItems = []dto.AltaiPhyCheckItemEmbed{}
	}

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan altai check fisik ke database", err)
		logger.Error("Gagal menyimpan altai check fisik ke database, (InsertCheck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *checkAltaiPhyDao) EditCheck(input dto.AltaiPhyCheckEdit) (*dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
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

			keyName:        input.Name,
			keyUpdatedAt:   input.UpdatedAt,
			keyUpdatedBy:   input.UpdatedBy,
			keyUpdatedByID: input.UpdatedByID,
			keyTimeEnded:   input.TimeEnded,
			keyIsFinish:    input.IsFinish,
			keyNote:        input.Note,
		},
	}

	var check dto.AltaiPhyCheck
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("altai check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan altai check fisik dari database (EditCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan altai check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiPhyDao) DeleteCheck(input dto.FilterIDBranchCreateGte) (*dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyID:        input.FilterID,
		keyBranch:    input.FilterBranch,
		keyCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var check dto.AltaiPhyCheck
	err := coll.FindOneAndDelete(ctx, filter).Decode(&check)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Altai Check fisik tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus Altai Check fisik dari database (AltaiDeleteCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus Altai Check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiPhyDao) UploadChildImage(filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
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

	var check dto.AltaiPhyCheck
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, altai check fisik dengan id %s -> %s tidak ditemukan", filterA.FilterParentID.Hex(), filterA.FilterChildID))
		}

		logger.Error("Memasukkan path image child altai check fisik ke db gagal, (AltaiUploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image altai check fisik ke db gagal", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiPhyDao) UpdateCheckItem(input dto.AltaiPhyCheckItemUpdate) (*dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
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
			keyUpdatedAt:       input.CheckedAt,
			keyChXCheckedAt:    input.CheckedAt,
			keyChXCheckedBy:    input.CheckedBy,
			keyChXIsChecked:    input.IsChecked,
			keyChXIsMaintained: input.IsMaintained,
			keyChXIsOffline:    input.IsOffline,
		},
	}

	var check dto.AltaiPhyCheck
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("altai check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (AltaiUpdateCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan altai check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkAltaiPhyDao) BulkUpdateItem(inputs []dto.AltaiPhyCheckItemUpdate) (int64, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
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
				keyUpdatedAt:       input.CheckedAt,
				keyChXCheckedAt:    input.CheckedAt,
				keyChXCheckedBy:    input.CheckedBy,
				keyChXIsChecked:    input.IsChecked,
				keyChXIsMaintained: input.IsMaintained,
				keyChXIsOffline:    input.IsOffline,
			},
		}
		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := coll.BulkWrite(ctx, operations, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, rest_err.NewBadRequestError("altai check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("gagal bulk write checkItem dari database (BulkUpdateItem)", err)
		apiErr := rest_err.NewInternalServerError("gagal bulk write altai check fisik dari database", err)
		return 0, apiErr
	}

	return result.ModifiedCount, nil
}

func (c *checkAltaiPhyDao) GetCheckByID(checkID primitive.ObjectID, branchIfSpecific string) (*dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: checkID}
	// filter condition
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var check dto.AltaiPhyCheck
	if err := coll.FindOne(ctx, filter).Decode(&check); err != nil {
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

func (c *checkAltaiPhyDao) FindCheck(branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
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
			keyAltaiPhyCheckItems: 0,
		})
	}
	opts.SetSort(bson.D{{keyUpdatedAt, -1}}) //nolint:govet
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar altai check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.AltaiPhyCheck{}, apiErr
	}

	var checkList []dto.AltaiPhyCheck
	if err = cursor.All(ctx, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (AltaiFindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.AltaiPhyCheck{}, apiErr
	}

	return checkList, nil
}

func (c *checkAltaiPhyDao) FindCheckStillOpen(branch string, detail bool) ([]dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{}
	// filter condition
	if branch != "" {
		filter[keyBranch] = strings.ToUpper(branch)
	}
	filter[keyIsFinish] = false

	opts := options.Find()

	if !detail {
		opts.SetProjection(bson.M{
			keyAltaiPhyCheckItems: 0,
		})
	}
	opts.SetSort(bson.D{{Key: keyUpdatedAt, Value: -1}})
	opts.SetLimit(100)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar altai check dari database (FindCheckStillOpen)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.AltaiPhyCheck{}, apiErr
	}

	var checkList []dto.AltaiPhyCheck
	if err = cursor.All(ctx, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (FindCheckStillOpen)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.AltaiPhyCheck{}, apiErr
	}

	return checkList, nil
}

func (c *checkAltaiPhyDao) GetLastCheckCreateRange(start, end int64, branch string, isQuarter bool) (*dto.AltaiPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{
		keyIsQuarterly: isQuarter,
	}
	// filter condition
	if branch != "" {
		filter[keyBranch] = strings.ToUpper(branch)
	}
	filter[keyCreatedAt] = bson.M{"$gte": start}
	filter[keyCreatedAt] = bson.M{"$lte": end}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyCreatedAt, Value: -1}})
	opts.SetLimit(1)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar altai check dari database (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var checkList []dto.AltaiPhyCheck
	if err = cursor.All(ctx, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var check dto.AltaiPhyCheck

	if len(checkList) != 0 {
		check = checkList[0]
	}

	return &check, nil
}
