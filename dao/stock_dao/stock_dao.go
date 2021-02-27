package stock_dao

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
	EditStock(input dto.StockEdit) (*dto.Stock, rest_err.APIError)
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

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan stock ke database", err)
		logger.Error("Gagal menyimpan stock ke database, (InsertStock)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (s *stockDao) EditStock(input dto.StockEdit) (*dto.Stock, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	if input.Tag == nil {
		input.Tag = []string{}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyStoID:        input.ID,
		keyStoBranch:    input.FilterBranch,
		keyStoUpdatedAt: input.FilterTimestamp,
	}

	update := bson.D{
		{"$set", bson.M{
			keyStoName:        input.Name,
			keyStoUpdatedAt:   input.UpdatedAt,
			keyStoUpdatedBy:   input.UpdatedBy,
			keyStoUpdatedByID: input.UpdatedByID,
			keyStoUnit:        input.Unit,
			keyStoLocation:    input.Location,
			keyStoThreshold:   input.Threshold,
			keyStoTag:         input.Tag,
			keyStoNote:        input.Note,
		}},
	}

	var stock dto.Stock
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&stock); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Stock tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan stock dari database (EditStock)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan stock dari database", err)
		return nil, apiErr
	}

	return &stock, nil
}
