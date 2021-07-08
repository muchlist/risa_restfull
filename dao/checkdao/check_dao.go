package checkdao

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
	connectTimeout  = 3
	keyChCollection = "check"

	keyChID          = "_id"
	keyChCreatedAt   = "created_at"
	keyChUpdatedAt   = "updated_at"
	keyChUpdatedBy   = "updated_by"
	keyChUpdatedByID = "updated_by_id"
	keyChCreatedByID = "created_by_id"
	keyChBranch      = "branch"

	keyChIsFinish   = "is_finish"
	keyChCheckItems = "check_items"
	keyChNote       = "note"

	keyCiXId = "check_items.id"

	keyCiXCheckedAt        = "check_items.$.checked_at"
	keyCiXIsChecked        = "check_items.$.is_checked"
	keyCiXTagSelected      = "check_items.$.tag_selected"
	keyCiXTagExtraSelected = "check_items.$.tag_extra_selected"
	keyCiXImagePath        = "check_items.$.image_path"
	keyCiXCheckedNote      = "check_items.$.checked_note"
	keyCiXHaveProblem      = "check_items.$.have_problem"
	keyCiXCompleteStatus   = "check_items.$.complete_status"
)

func NewCheckDao() CheckDaoAssumer {
	return &checkDao{}
}

type checkDao struct {
}

type CheckDaoAssumer interface {
	InsertCheck(input dto.Check) (*string, rest_err.APIError)
	EditCheck(input dto.CheckEdit) (*dto.Check, rest_err.APIError)
	DeleteCheck(input dto.FilterIDBranchCreateGte) (*dto.Check, rest_err.APIError)
	UploadChildImage(filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.Check, rest_err.APIError)
	UpdateCheckItem(input dto.CheckChildUpdate) (*dto.Check, rest_err.APIError)

	GetCheckByID(checkID primitive.ObjectID, branchIfSpecific string) (*dto.Check, rest_err.APIError)
	FindCheck(branch string, filterA dto.FilterTimeRangeLimit) (dto.CheckResponseMinList, rest_err.APIError)
	FindCheckForReports(branch string, filterA dto.FilterTimeRangeLimit) ([]dto.Check, rest_err.APIError)
}

func (c *checkDao) InsertCheck(input dto.Check) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// Default value for slice
	input.Branch = strings.ToUpper(input.Branch)
	if input.CheckItems == nil {
		input.CheckItems = []dto.CheckItemEmbed{}
	}
	for _, v := range input.CheckItems {
		if v.Tag == nil {
			v.Tag = []string{}
		}
		if v.TagExtra == nil {
			v.TagExtra = []string{}
		}
	}

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan check ke database", err)
		logger.Error("Gagal menyimpan check ke database, (InsertCheck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *checkDao) EditCheck(input dto.CheckEdit) (*dto.Check, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// Default value for slice

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyChID:          input.FilterID,
		keyChBranch:      input.FilterBranch,
		keyChCreatedByID: input.FilterAuthorID,
		keyChIsFinish:    false,
	}

	update := bson.M{
		"$set": bson.M{
			keyChUpdatedAt:   input.UpdatedAt,
			keyChUpdatedBy:   input.UpdatedBy,
			keyChUpdatedByID: input.UpdatedByID,

			keyChIsFinish: input.IsFinish,
			keyChNote:     input.Note,
		},
	}

	var check dto.Check
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Check tidak diupdate : validasi id branch author isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (EditCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) DeleteCheck(input dto.FilterIDBranchCreateGte) (*dto.Check, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyChID:        input.FilterID,
		keyChBranch:    input.FilterBranch,
		keyChCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var check dto.Check
	err := coll.FindOneAndDelete(ctx, filter).Decode(&check)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Check tidak diupdate : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus check dari database (DeleteCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) UploadChildImage(filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.Check, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyChID:          filterA.FilterParentID,
		keyCiXId:         filterA.FilterChildID,
		keyChCreatedByID: filterA.FilterAuthorID,
	}

	update := bson.M{
		"$set": bson.M{
			keyCiXImagePath: imagePath,
		},
	}

	var check dto.Check
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, check dengan id %s -> %s tidak ditemukan", filterA.FilterParentID.Hex(), filterA.FilterChildID))
		}

		logger.Error("Memasukkan path image child Check ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image check ke db gagal", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) UpdateCheckItem(input dto.CheckChildUpdate) (*dto.Check, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyChID:          input.FilterParentID,
		keyCiXId:         input.FilterChildID,
		keyChCreatedByID: input.FilterAuthorID,
		keyChIsFinish:    false,
	}

	update := bson.M{
		"$set": bson.M{
			keyCiXIsChecked:        input.IsChecked,
			keyCiXCheckedAt:        input.CheckedAt,
			keyCiXCheckedNote:      input.CheckedNote,
			keyCiXHaveProblem:      input.HaveProblem,
			keyCiXCompleteStatus:   input.CompleteStatus,
			keyCiXTagSelected:      input.TagSelected,
			keyCiXTagExtraSelected: input.TagExtraSelected,
		},
	}

	var check dto.Check
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Check tidak diupdate : validasi id branch author isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (UpdateCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) GetCheckByID(checkID primitive.ObjectID, branchIfSpecific string) (*dto.Check, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyChID: checkID}
	// filter condition
	if branchIfSpecific != "" {
		filter[keyChBranch] = strings.ToUpper(branchIfSpecific)
	}

	var check dto.Check
	if err := coll.FindOne(ctx, filter).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Check dengan ID %s tidak ditemukan. validation : id branch", checkID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan check dari database (GetCheckByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) FindCheck(branch string, filterA dto.FilterTimeRangeLimit) (dto.CheckResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	if filterA.Limit == 0 {
		filterA.Limit = 100
	}

	// filter
	filter := bson.M{}
	// filter condition
	if branch != "" {
		filter[keyChBranch] = strings.ToUpper(branch)
	}
	if filterA.FilterStart != 0 {
		filter[keyChCreatedAt] = bson.M{"$gte": filterA.FilterStart}
	}
	if filterA.FilterEnd != 0 {
		filter[keyChCreatedAt] = bson.M{"$lte": filterA.FilterEnd}
	}

	opts := options.Find()
	opts.SetProjection(bson.M{
		keyChCheckItems: 0,
	})
	opts.SetSort(bson.D{{keyChID, -1}}) //nolint:govet
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CheckResponseMinList{}, apiErr
	}

	checkList := dto.CheckResponseMinList{}
	if err = cursor.All(ctx, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CheckResponseMinList{}, apiErr
	}

	return checkList, nil
}

// FindCheckForReports mengembalikan check detail dengan limit 2
func (c *checkDao) FindCheckForReports(branch string, filterA dto.FilterTimeRangeLimit) ([]dto.Check, rest_err.APIError) {
	coll := db.DB.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	if filterA.Limit == 0 {
		filterA.Limit = 2
	}

	// filter
	filter := bson.M{}
	// filter condition
	if branch != "" {
		filter[keyChBranch] = strings.ToUpper(branch)
	}
	if filterA.FilterStart != 0 {
		filter[keyChCreatedAt] = bson.M{"$gte": filterA.FilterStart}
	}
	if filterA.FilterEnd != 0 {
		filter[keyChCreatedAt] = bson.M{"$lte": filterA.FilterEnd}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyChID, -1}}) //nolint:govet
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.Check{}, apiErr
	}

	var checkList []dto.Check
	if err = cursor.All(ctx, &checkList); err != nil {
		logger.Error("Gagal decode checkList cursor ke objek slice (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.Check{}, apiErr
	}

	return checkList, nil
}
