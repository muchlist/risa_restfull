package vendorcheckdao

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
	keyCollection  = "vendorCheck"
	keyID          = "_id"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
	keyUpdatedBy   = "updated_by"
	keyUpdatedByID = "updated_by_id"
	keyBranch      = "branch"

	keyTimeStarted      = "time_started"
	keyTimeEnded        = "time_ended"
	keyIsFinish         = "is_finish"
	keyIsCheckByVendor  = "is_check_by_vendor"
	keyVendorCheckItems = "vendor_check_items"
	keyNote             = "note"

	keyChXId        = "vendor_check_items.id"
	keyChXCheckedAt = "vendor_check_items.$.checked_at"
	keyChXCheckedBy = "vendor_check_items.$.checked_by"
	keyChXIsChecked = "vendor_check_items.$.is_checked"
	keyChXIsBlur    = "vendor_check_items.$.is_blur"
	keyChXIsOffline = "vendor_check_items.$.is_offline"
	keyChXImagePath = "vendor_check_items.$.image_path"
)

func NewVendorCheckDao() CheckVendorDaoAssumer {
	return &checkVendorDao{}
}

type checkVendorDao struct {
}

func (c *checkVendorDao) InsertCheck(ctx context.Context, input dto.VendorCheck) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// Default value for slice
	input.Branch = strings.ToUpper(input.Branch)
	if input.VendorCheckItems == nil {
		input.VendorCheckItems = []dto.VendorCheckItemEmbed{}
	}

	result, err := coll.InsertOne(ctxt, input)

	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan cctv check ke database", err)
		logger.Error("Gagal menyimpan cctv check ke database, (VendorInsertCheck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *checkVendorDao) EditCheck(ctx context.Context, input dto.VendorCheckEdit) (*dto.VendorCheck, rest_err.APIError) {
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

	var check dto.VendorCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("cctv check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan cctv check dari database (EditCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVendorDao) DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.VendorCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyID:        input.FilterID,
		keyBranch:    input.FilterBranch,
		keyCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var check dto.VendorCheck
	err := coll.FindOneAndDelete(ctxt, filter).Decode(&check)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Cctv Check tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus cctvCheck dari database (VendorDeleteCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus cctv check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVendorDao) UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.VendorCheck, rest_err.APIError) {
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

	var check dto.VendorCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, cctv check dengan id %s -> %s tidak ditemukan", filterA.FilterParentID.Hex(), filterA.FilterChildID))
		}

		logger.Error("Memasukkan path image child cctv Check ke db gagal, (VendorUploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image cctv check ke db gagal", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVendorDao) UpdateCheckItem(ctx context.Context, input dto.VendorCheckItemUpdate) (*dto.VendorCheck, rest_err.APIError) {
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
			keyChXIsBlur:    input.IsBlur,
			keyChXIsOffline: input.IsOffline,
		},
	}

	var check dto.VendorCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("cctv check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (VendorUpdateCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVendorDao) BulkUpdateItem(ctx context.Context, inputs []dto.VendorCheckItemUpdate) (int64, rest_err.APIError) {
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
				keyChXIsBlur:    input.IsBlur,
				keyChXIsOffline: input.IsOffline,
			},
		}
		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := coll.BulkWrite(ctxt, operations, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, rest_err.NewBadRequestError("cctv check tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("gagal bulk write checkItem dari database (BulkUpdateItem)", err)
		apiErr := rest_err.NewInternalServerError("gagal bulk write cctv check dari database", err)
		return 0, apiErr
	}

	return result.ModifiedCount, nil
}

func (c *checkVendorDao) GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.VendorCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: checkID}
	// filter condition
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var check dto.VendorCheck
	if err := coll.FindOne(ctxt, filter).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("cctv check dengan ID %s tidak ditemukan. validation : id branch", checkID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan cctv check dari database (GetCheckByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVendorDao) FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.VendorCheck, rest_err.APIError) {
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
			keyVendorCheckItems: 0,
		})
	}
	opts.SetSort(bson.D{{Key: keyUpdatedAt, Value: -1}}) //nolint:govet
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar cctv check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.VendorCheck{}, apiErr
	}

	var checkList []dto.VendorCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (VendorFindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.VendorCheck{}, apiErr
	}

	return checkList, nil
}

func (c *checkVendorDao) GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string) (*dto.VendorCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{
		keyIsCheckByVendor: true,
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

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar cctv check dari database (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var checkList []dto.VendorCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode cctv virtual ke objek slice (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var check dto.VendorCheck
	if len(checkList) != 0 {
		check = checkList[0]
	}

	return &check, nil
}
