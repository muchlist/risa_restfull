package check_dao

import (
	"context"
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

	keyChShift      = "shift"
	keyChIsFinish   = "is_finish"
	keyChCheckItems = "check_items"
	keyChNote       = "note"

	keyCiXId        = "check_items._id"
	keyCiXName      = "check_items.$.name"
	keyCiXLocation  = "check_items.$.location"
	keyCiXType      = "check_items.$.type"
	keyCiXTag       = "check_items.$.tag"
	keyCiXTag_extra = "check_items.$.tag_extra"

	keyCiXChecked_at         = "check_items.$.checked_at"
	keyCiXIs_checked         = "check_items.$.is_checked"
	keyCiXTag_selected       = "check_items.$.tag_selected"
	keyCiXTag_extra_selected = "check_items.$.tag_extra_selected"
	keyCiXImage_path         = "check_items.$.image_path"
	keyCiXChecked_note       = "check_items.$.checked_note"
	keyCiXHave_problem       = "check_items.$.have_problem"
	keyCiXComplete_status    = "check_items.$.complete_status"
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
	//UploadImage(filterA FilterParentIDChildIDAuthor, imagePath string) (*dto.Check, rest_err.APIError)
	//
	//GetCheckByID(checkID primitive.ObjectID, branchIfSpecific string) (*dto.Check, rest_err.APIError)
	//FindCheck(filter dto.FilterBranchLocIPNameDisable) (dto.CheckResponseMinList, rest_err.APIError)
}

func (c *checkDao) InsertCheck(input dto.Check) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
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
	coll := db.Db.Collection(keyChCollection)
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
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Check tidak diupdate : validasi id branch author isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (EditCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) DeleteCheck(input dto.FilterIDBranchCreateGte) (*dto.Check, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
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
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Check tidak diupdate : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus check dari database (DeleteCheck)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) UploadChildImage(filterA dto.FilterParentIDChildIDAuthor, imagePath string) (*dto.Check, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
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
			keyCiXImage_path: imagePath,
		},
	}

	var check dto.Check
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, check dengan id %s -> %s tidak ditemukan", filterA.FilterParentID.Hex(), filterA.FilterChildID))
		}

		logger.Error("Memasukkan path image child Check ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image check ke db gagal", err)
		return nil, apiErr
	}

	return &check, nil
}

func (c *checkDao) UpdateCheckItem(input dto.CheckChildUpdate) (*dto.Check, rest_err.APIError) {
	coll := db.Db.Collection(keyChCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// Default value for slice

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyChID:          input.FilterParentID,
		keyCiXId:         input.FilterChildID,
		keyChCreatedByID: input.FilterAuthorID,
	}

	update := bson.M{
		"$set": bson.M{
			keyCiXIs_checked:         input.IsChecked,
			keyCiXChecked_note:       input.CheckedNote,
			keyCiXHave_problem:       input.HaveProblem,
			keyCiXComplete_status:    input.CompleteStatus,
			keyCiXTag_selected:       input.TagSelected,
			keyCiXTag_extra_selected: input.TagExtraSelected,
		},
	}

	var check dto.Check
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&check); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Check tidak diupdate : validasi id branch author isFinish")
		}

		logger.Error("Gagal mendapatkan checkItem dari database (UpdateCheckItem)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan check dari database", err)
		return nil, apiErr
	}

	return &check, nil
}
