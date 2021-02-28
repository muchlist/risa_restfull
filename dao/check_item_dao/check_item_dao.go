package check_item_dao

import (
	"context"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	connectTimeout  = 3
	keyChCollection = "checkItem"

	keyChID             = "_id"
	keyChUpdatedAt      = "updated_at"
	keyChUpdatedBy      = "updated_by"
	keyChUpdatedByID    = "updated_by_id"
	keyChBranch         = "branch"
	keyChDisable        = "disable"
	keyChName           = "name"
	keyChLocation       = "location"
	keyChLocationLat    = "location_lat"
	keyChLocationLon    = "location_lon"
	keyChType           = "type"
	keyChTag            = "tag"
	keyChTagExtra       = "tag_extra"
	keyChNote           = "note"
	keyChShifts         = "shifts"
	keyChCheckedNote    = "checked_note"
	keyChHaveProblem    = "have_problem"
	keyChCompleteStatus = "complete_status"
)

func NewCheckItemDao() CheckItemDaoAssumer {
	return &checkItemDao{}
}

type checkItemDao struct {
}

type CheckItemDaoAssumer interface {
	InsertCheckItem(input dto.CheckItem) (*string, rest_err.APIError)
}

func (c *checkItemDao) InsertCheckItem(input dto.CheckItem) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// Default value
	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)
	if input.Shifts == nil {
		input.Shifts = []string{}
	}
	if input.Tag == nil {
		input.Tag = []string{}
	}
	if input.TagExtra == nil {
		input.TagExtra = []string{}
	}
	input.Disable = false

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan checkItem ke database", err)
		logger.Error("Gagal menyimpan checkItem ke database, (InsertCheckItem)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *checkItemDao) EditCheckItem(input dto.CheckItemEdit) (*dto.CheckItem, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	if input.Shifts == nil {
		input.Shifts = []string{}
	}
	if input.Tag == nil {
		input.Tag = []string{}
	}
	if input.TagExtra == nil {
		input.TagExtra = []string{}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyChID:        input.FilterID,
		keyChBranch:    input.FilterBranch,
		keyChUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyChName:        input.Name,
			keyChUpdatedAt:   input.UpdatedAt,
			keyChUpdatedBy:   input.UpdatedBy,
			keyChUpdatedByID: input.UpdatedByID,

			keyChLocation:    input.Location,
			keyChLocationLat: input.LocationLat,
			keyChLocationLon: input.LocationLon,

			keyChTag:      input.Tag,
			keyChTagExtra: input.TagExtra,
			keyChType:     input.Type,
			keyChNote:     input.Note,
			keyChShifts:   input.Shifts,
		},
	}

	var checkItem dto.CheckItem
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&checkItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("CheckItem tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (EditCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan checkItem dari database", err)
		return nil, apiErr
	}

	return &checkItem, nil
}

func (c *checkItemDao) EditCheckItemBySystem(input dto.CheckItemEditBySys) (*dto.CheckItem, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyChID: input.FilterID,
	}

	update := bson.M{
		"$set": bson.M{
			keyChUpdatedAt:      input.UpdatedAt,
			keyChCheckedNote:    input.CheckedNote,
			keyChHaveProblem:    input.HaveProblem,
			keyChCompleteStatus: input.CompleteStatus,
		},
	}

	var checkItem dto.CheckItem
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&checkItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("CheckItem tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (EditCheckItemBySystem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan checkItem dari database", err)
		return nil, apiErr
	}

	return &checkItem, nil
}

func (c *checkItemDao) DeleteCheckItem(input dto.FilterIDBranch) (*dto.CheckItem, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyChID:     input.FilterID,
		keyChBranch: input.FilterBranch,
	}

	var checkItem dto.CheckItem
	err := coll.FindOneAndDelete(ctx, filter).Decode(&checkItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("CheckItem tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal menghapus checkItem dari database (DeleteCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus checkItem dari database", err)
		return nil, apiErr
	}

	return &checkItem, nil
}

// DisableCheckItem if value true , checkItem will disabled
func (c *checkItemDao) DisableCheckItem(checkItemID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.CheckItem, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyChID:     checkItemID,
		keyChBranch: user.Branch,
	}

	update := bson.M{
		"$set": bson.M{
			keyChDisable:     value,
			keyChUpdatedAt:   time.Now().Unix(),
			keyChUpdatedByID: user.Identity,
			keyChUpdatedBy:   user.Name,
		},
	}

	var checkItem dto.CheckItem
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&checkItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("CheckItem tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendisable checkItem dari database (DisableCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendisable checkItem dari database", err)
		return nil, apiErr
	}

	return &checkItem, nil
}

func (c *checkItemDao) GetCheckItemByID(checkItemID primitive.ObjectID, branchIfSpecific string) (*dto.CheckItem, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyChID: checkItemID}
	if branchIfSpecific != "" {
		filter[keyChBranch] = strings.ToUpper(branchIfSpecific)
	}

	var checkItem dto.CheckItem
	if err := coll.FindOne(ctx, filter).Decode(&checkItem); err != nil {

		if err == mongo.ErrNoDocuments {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("CheckItem dengan ID %s tidak ditemukan", checkItemID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan checkItem dari database (GetCheckItemByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan checkItem dari database", err)
		return nil, apiErr
	}

	return &checkItem, nil
}

func (c *checkItemDao) FindCheckItem(filterA dto.FilterBranchNameDisable, filterHaveProblem bool) (dto.CheckItemResponseMinList, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterName = strings.ToUpper(filterA.FilterName)

	// filter
	filter := bson.M{
		keyChDisable: filterA.FilterDisable,
	}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyChBranch] = filterA.FilterBranch
	}
	if filterA.FilterName != "" {
		filter[keyChName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterName),
		}
	}
	if filterHaveProblem {
		filter[keyChHaveProblem] = true
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyChType, -1}, {keyChName, 1}})

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar checkItem dari database (FindCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CheckItemResponseMinList{}, apiErr
	}

	checkItemList := dto.CheckItemResponseMinList{}
	if err = cursor.All(ctx, &checkItemList); err != nil {
		logger.Error("Gagal decode checkItemList cursor ke objek slice (FindCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CheckItemResponseMinList{}, apiErr
	}

	return checkItemList, nil
}
