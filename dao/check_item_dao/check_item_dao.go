package check_item_dao

import (
	"context"
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
	keyChCollection = "checkItem"

	keyChID             = "_id"
	keyChCreatedAt      = "created_at"
	keyChCreatedBy      = "created_by"
	keyChCreatedByID    = "created_by_id"
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
