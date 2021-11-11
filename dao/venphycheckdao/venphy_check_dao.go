package venphycheckdao

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
	keyCollection  = "venPhyCheck"
	keyID          = "_id"
	keyName        = "name"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
	keyUpdatedBy   = "updated_by"
	keyUpdatedByID = "updated_by_id"
	keyBranch      = "branch"

	keyTimeEnded        = "time_ended"
	keyIsFinish         = "is_finish"
	keyIsQuarterly      = "quarterly_mode"
	keyVenPhyCheckItems = "ven_phy_check_items"
	keyNote             = "note"

	keyChXId           = "ven_phy_check_items.id"
	keyChXCheckedAt    = "ven_phy_check_items.$.checked_at"
	keyChXCheckedBy    = "ven_phy_check_items.$.checked_by"
	keyChXIsChecked    = "ven_phy_check_items.$.is_checked"
	keyChXIsMaintained = "ven_phy_check_items.$.is_maintained"
	keyChXIsBlur       = "ven_phy_check_items.$.is_blur"
	keyChXIsOffline    = "ven_phy_check_items.$.is_offline"
	keyChXImagePath    = "ven_phy_check_items.$.image_path"
)

func NewVenPhyCheckDao() CheckVenPhyDaoAssumer {
	return &checkVenPhyDao{}
}

type checkVenPhyDao struct {
}

func (c *checkVenPhyDao) InsertCheck(ctx context.Context, input dto.VenPhyCheck) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// Default value for slice
	input.Branch = strings.ToUpper(input.Branch)
	if input.VenPhyCheckItems == nil {
		input.VenPhyCheckItems = []dto.VenPhyCheckItemEmbed{}
	}

	result, err := coll.InsertOne(ctxt, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan vendor check fisik ke database", err)
		logger.Error("Gagal menyimpan vendor check fisik ke database, (InsertCheck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *checkVenPhyDao) EditCheck(ctx context.Context, input dto.VenPhyCheckEdit) (*dto.VenPhyCheck, rest_err.APIError) {
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

	set := bson.M{
		keyUpdatedAt:   input.UpdatedAt,
		keyUpdatedBy:   input.UpdatedBy,
		keyUpdatedByID: input.UpdatedByID,
		keyTimeEnded:   input.TimeEnded,
		keyIsFinish:    input.IsFinish,
		keyNote:        input.Note,
	}
	if input.Name != "" {
		set[keyName] = input.Name
	}

	update := bson.M{
		"$set": set,
	}

	var check dto.VenPhyCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("vendor check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan vendor check fisik dari database (EditCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan vendor check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVenPhyDao) UndoFinishCheck(ctx context.Context, filterID primitive.ObjectID, filterBranch string) (*dto.VenPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// Default value for slice

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:     filterID,
		keyBranch: filterBranch,
	}

	set := bson.M{
		keyIsFinish: false,
	}

	update := bson.M{
		"$set": set,
	}

	var check dto.VenPhyCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("vendor check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan vendor check fisik dari database (UndoFinishCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan vendor check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

// OverwriteChecklist akan menindih hasil checklist sebelumnya, digunakan untuk system merename ulang nama cctv yang sudah di edit pada parrent
func (c *checkVenPhyDao) OverwriteChecklist(ctx context.Context, id primitive.ObjectID, checkItems []dto.VenPhyCheckItemEmbed) (*dto.VenPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	// Default value for slice

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID: id,
	}

	update := bson.M{
		"$set": bson.M{
			keyVenPhyCheckItems: checkItems,
		},
	}

	var check dto.VenPhyCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("vendor check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan vendor check fisik dari database (OverwriteChecklist)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan vendor check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVenPhyDao) DeleteCheck(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.VenPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyID:        input.FilterID,
		keyBranch:    input.FilterBranch,
		keyCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var check dto.VenPhyCheck
	err := coll.FindOneAndDelete(ctxt, filter).Decode(&check)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Vendor Check fisik tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus Vendor Check fisik dari database (VendorDeleteCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus Vendor Check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVenPhyDao) UploadChildImage(ctx context.Context, filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.VenPhyCheck, rest_err.APIError) {
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

	var check dto.VenPhyCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, vendor check fisik dengan id %s -> %s tidak ditemukan", filterA.FilterParentID.Hex(), filterA.FilterChildID))
		}

		logger.Error("Memasukkan path image child vendor check fisik ke db gagal, (VendorUploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image vendor check fisik ke db gagal", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVenPhyDao) UpdateCheckItem(ctx context.Context, input dto.VenPhyCheckItemUpdate) (*dto.VenPhyCheck, rest_err.APIError) {
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
			keyUpdatedAt:       input.CheckedAt,
			keyChXCheckedAt:    input.CheckedAt,
			keyChXCheckedBy:    input.CheckedBy,
			keyChXIsChecked:    input.IsChecked,
			keyChXIsMaintained: input.IsMaintained,
			keyChXIsBlur:       input.IsBlur,
			keyChXIsOffline:    input.IsOffline,
		},
	}

	var check dto.VenPhyCheck
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("vendor check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (VendorUpdateCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan vendor check fisik dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVenPhyDao) BulkUpdateItem(ctx context.Context, inputs []dto.VenPhyCheckItemUpdate) (int64, rest_err.APIError) {
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

		setMap := bson.M{
			keyUpdatedAt:       input.CheckedAt,
			keyChXCheckedAt:    input.CheckedAt,
			keyChXCheckedBy:    input.CheckedBy,
			keyChXIsBlur:       input.IsBlur,
			keyChXIsOffline:    input.IsOffline,
			keyChXIsChecked:    input.IsChecked,
			keyChXIsMaintained: input.IsMaintained,
		}

		update := bson.M{
			"$set": setMap,
		}

		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := coll.BulkWrite(ctxt, operations, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, rest_err.NewBadRequestError("vendor check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("gagal bulk write checkItem dari database (BulkUpdateItem)", err)
		apiErr := rest_err.NewInternalServerError("gagal bulk write vendor check fisik dari database", err)
		return 0, apiErr
	}

	return result.ModifiedCount, nil
}

// BulkUpdateItemForUpdateCheckItem
// jika ismaintenance false, atau ischecked false maka tidak akan merubah nilai
func (c *checkVenPhyDao) BulkUpdateItemForUpdateCheckItem(ctx context.Context, inputs []dto.VenPhyCheckItemUpdate) (int64, rest_err.APIError) {
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

		setMap := bson.M{
			keyUpdatedAt:    input.CheckedAt,
			keyChXCheckedAt: input.CheckedAt,
			keyChXCheckedBy: input.CheckedBy,
			keyChXIsBlur:    input.IsBlur,
			keyChXIsOffline: input.IsOffline,
		}

		if input.IsChecked {
			setMap[keyChXIsChecked] = input.IsChecked
		}
		if input.IsMaintained {
			setMap[keyChXIsMaintained] = input.IsMaintained
		}

		update := bson.M{
			"$set": setMap,
		}

		operations[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false)
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := coll.BulkWrite(ctxt, operations, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, rest_err.NewBadRequestError("vendor check fisik tidak diupdate : validasi id branch isFinish")
		}

		logger.Error("gagal bulk write checkItem dari database (BulkUpdateItem)", err)
		apiErr := rest_err.NewInternalServerError("gagal bulk write vendor check fisik dari database", err)
		return 0, apiErr
	}

	return result.ModifiedCount, nil
}

func (c *checkVenPhyDao) GetCheckByID(ctx context.Context, checkID primitive.ObjectID, branchIfSpecific string) (*dto.VenPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: checkID}
	// filter condition
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var check dto.VenPhyCheck
	if err := coll.FindOne(ctxt, filter).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("vendor check dengan ID %s tidak ditemukan. validation : id branch", checkID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan vendor check dari database (GetCheckByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan vendor check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkVenPhyDao) FindCheck(ctx context.Context, branch string, filterA dto.FilterTimeRangeLimit, detail bool) ([]dto.VenPhyCheck, rest_err.APIError) {
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
			keyVenPhyCheckItems: 0,
		})
	}
	opts.SetSort(bson.D{{keyUpdatedAt, -1}}) //nolint:govet
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar vendor check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.VenPhyCheck{}, apiErr
	}

	var checkList []dto.VenPhyCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (VendorFindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.VenPhyCheck{}, apiErr
	}

	return checkList, nil
}

func (c *checkVenPhyDao) FindCheckStillOpen(ctx context.Context, branch string, detail bool) ([]dto.VenPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
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
			keyVenPhyCheckItems: 0,
		})
	}
	opts.SetSort(bson.D{{Key: keyUpdatedAt, Value: -1}})
	opts.SetLimit(100)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar vendor check dari database (FindCheckStillOpen)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.VenPhyCheck{}, apiErr
	}

	var checkList []dto.VenPhyCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (FindCheckStillOpen)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.VenPhyCheck{}, apiErr
	}

	return checkList, nil
}

func (c *checkVenPhyDao) GetLastCheckCreateRange(ctx context.Context, start, end int64, branch string, isQuarter bool) (*dto.VenPhyCheck, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
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

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("gagal mendapatkan daftar cctv dari database (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var checkList []dto.VenPhyCheck
	if err = cursor.All(ctxt, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (GetLastCheckCreateRange)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return nil, apiErr
	}

	var check dto.VenPhyCheck
	if len(checkList) != 0 {
		check = checkList[0]
	}

	return &check, nil
}
