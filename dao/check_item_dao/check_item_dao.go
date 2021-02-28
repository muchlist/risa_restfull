package check_item_dao

import (
	"context"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
