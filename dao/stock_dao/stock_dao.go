package stock_dao

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
	keyStoCollection = "stock"

	keyStoID          = "_id"
	keyStoName        = "name"
	keyStoCreatedAt   = "created_at"
	keyStoCreatedBy   = "created_by"
	keyStoCreatedByID = "created_by_id"
	keyStoUpdatedAt   = "updated_at"
	keyStoUpdatedBy   = "updated_by"
	keyStoUpdatedByID = "updated_by_id"
	keyStoBranch      = "branch"
	keyStoDisable     = "disable"

	keyStoCategory  = "stock_category"
	keyStoUnit      = "unit"
	keyStoQty       = "qty"
	keyStoThreshold = "threshold"
	keyStoIncrement = "increment"
	keyStoDecrement = "decrement"
	keyStoLocation  = "location"
	keyStoTag       = "tag"
	keyStoImage     = "image"
	keyStoNote      = "note"
)

func NewStockDao() StockDaoAssumer {
	return &stockDao{}
}

type stockDao struct {
}

type StockDaoAssumer interface {
	InsertStock(input dto.Stock) (*string, rest_err.APIError)
}

func (s *stockDao) InsertStock(input dto.Stock) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)
	if input.Tag == nil {
		input.Tag = []string{}
	}
	if input.Increment == nil {
		input.Increment = []dto.StockChange{}
	}
	if input.Decrement == nil {
		input.Decrement = []dto.StockChange{}
	}

	insertDoc := bson.M{
		keyStoID:          input.ID,
		keyStoName:        input.Name,
		keyStoCreatedAt:   input.CreatedAt,
		keyStoCreatedBy:   input.CreatedBy,
		keyStoCreatedByID: input.CreatedByID,
		keyStoUpdatedAt:   input.UpdatedAt,
		keyStoUpdatedBy:   input.UpdatedBy,
		keyStoUpdatedByID: input.UpdatedByID,
		keyStoBranch:      input.Branch,
		keyStoDisable:     false,

		keyStoCategory:  input.StockCategory,
		keyStoUnit:      input.Unit,
		keyStoQty:       input.Qty,
		keyStoThreshold: input.Threshold,
		keyStoIncrement: input.Increment,
		keyStoDecrement: input.Decrement,
		keyStoLocation:  input.Location,

		keyStoTag:   input.Tag,
		keyStoImage: input.Image,
		keyStoNote:  input.Note,
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan stock ke database", err)
		logger.Error("Gagal menyimpan stock ke database, (InsertStock)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}
