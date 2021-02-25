package cctv_dao

import (
	"context"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

const (
	connectTimeout   = 3
	keyCtvCollection = "cctv"

	keyCtvID          = "_id"
	keyCtvName        = "name"
	keyCtvCreatedAt   = "created_at"
	keyCtvCreatedBy   = "created_by"
	keyCtvCreatedByID = "created_by_id"
	keyCtvUpdatedAt   = "updated_at"
	keyCtvUpdatedBy   = "updated_by"
	keyCtvUpdatedByID = "updated_by_id"
	keyCtvBranch      = "branch"

	keyCtvIP              = "ip"
	keyCtvInventoryNumber = "inventory_number"
	keyCtvLocation        = "location"
	keyCtvLocationLat     = "location_lat"
	keyCtvLocationLon     = "location_lon"
	keyCtvDate            = "date"
	keyCtvTag             = "tag"
	keyCtvImage           = "image"
	keyCtvBrand           = "brand"
	keyCtvType            = "type"
	keyCtvNote            = "note"
)

func NewCctvDao() CctvDaoAssumer {
	return &cctvDao{}
}

type cctvDao struct {
}

type CctvDaoAssumer interface {
	InsertCctv(input dto.Cctv) (*string, rest_err.APIError)
}

func (c *cctvDao) InsertCctv(input dto.Cctv) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)

	insertDoc := bson.M{
		keyCtvID:          input.ID,
		keyCtvName:        input.Name,
		keyCtvCreatedAt:   input.CreatedAt,
		keyCtvCreatedBy:   input.CreatedBy,
		keyCtvCreatedByID: input.CreatedByID,
		keyCtvUpdatedAt:   input.UpdatedAt,
		keyCtvUpdatedBy:   input.UpdatedBy,
		keyCtvUpdatedByID: input.UpdatedByID,
		keyCtvBranch:      input.Branch,

		keyCtvIP:              input.IP,
		keyCtvInventoryNumber: input.InventoryNumber,
		keyCtvLocation:        input.Location,
		keyCtvLocationLat:     input.LocationLat,
		keyCtvLocationLon:     input.LocationLon,

		keyCtvDate:  input.Date,
		keyCtvTag:   input.Tag,
		keyCtvImage: input.Image,
		keyCtvBrand: input.Brand,
		keyCtvType:  input.Type,
		keyCtvNote:  input.Note,
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan cctv ke database", err)
		logger.Error("Gagal menyimpan cctv ke database, (InsertCctv)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}
