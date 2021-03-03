package stock_dao

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
	"math"
	"strings"
	"time"
)

const (
	connectTimeout   = 3
	keyStoCollection = "stock"

	keyStoID          = "_id"
	keyStoName        = "name"
	keyStoCreatedAt   = "created_at"
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
	DeleteStock(input dto.FilterIDBranchCreateGte) (*dto.Stock, rest_err.APIError)
	DisableStock(stockID primitive.ObjectID, user mjwt.CustomClaim, isDisable bool) (*dto.Stock, rest_err.APIError)
	UploadImage(stockID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Stock, rest_err.APIError)
	ChangeQtyStock(filterA dto.FilterIDBranch, data dto.StockChange) (*dto.Stock, rest_err.APIError)

	GetStockByID(stockID primitive.ObjectID, branchIfSpecific string) (*dto.Stock, rest_err.APIError)
	FindStock(filterA dto.FilterBranchNameCatDisable) (dto.StockResponseMinList, rest_err.APIError)
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

func (s *stockDao) DeleteStock(input dto.FilterIDBranchCreateGte) (*dto.Stock, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyStoID:        input.FilterID,
		keyStoBranch:    input.FilterBranch,
		keyStoCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var stock dto.Stock
	err := coll.FindOneAndDelete(ctx, filter).Decode(&stock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Stock tidak diupdate : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus stock dari database (DeleteStock)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan stock dari database", err)
		return nil, apiErr
	}

	return &stock, nil
}

// DisableStock if value true , stock will disabled
func (s *stockDao) DisableStock(stockID primitive.ObjectID, user mjwt.CustomClaim, isDisable bool) (*dto.Stock, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyStoID:     stockID,
		keyStoBranch: user.Branch,
	}

	update := bson.M{
		"$set": bson.M{
			keyStoDisable:     isDisable,
			keyStoUpdatedAt:   time.Now().Unix(),
			keyStoUpdatedByID: user.Identity,
			keyStoUpdatedBy:   user.Name,
		},
	}

	var stock dto.Stock
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&stock); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Stock tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendapatkan stock dari database (DisableStock)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan stock dari database", err)
		return nil, apiErr
	}

	return &stock, nil
}

func (s *stockDao) UploadImage(stockID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Stock, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyStoID:     stockID,
		keyStoBranch: strings.ToUpper(filterBranch),
	}
	update := bson.M{
		"$set": bson.M{
			keyStoImage: imagePath,
		},
	}

	var stock dto.Stock
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&stock); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, stock dengan id %s tidak ditemukan", stockID.Hex()))
		}

		logger.Error("Memasukkan path image stock ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image stock ke db gagal", err)
		return nil, apiErr
	}

	return &stock, nil
}

func (s *stockDao) GetStockByID(stockID primitive.ObjectID, branchIfSpecific string) (*dto.Stock, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyStoID: stockID}
	if branchIfSpecific != "" {
		filter[keyStoBranch] = strings.ToUpper(branchIfSpecific)
	}

	var stock dto.Stock
	if err := coll.FindOne(ctx, filter).Decode(&stock); err != nil {

		if err == mongo.ErrNoDocuments {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Stock dengan ID %s tidak ditemukan", stockID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan stock dari database (GetStockByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan stock dari database", err)
		return nil, apiErr
	}

	return &stock, nil
}

func (s *stockDao) FindStock(filterA dto.FilterBranchNameCatDisable) (dto.StockResponseMinList, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterName = strings.ToUpper(filterA.FilterName)

	// filter
	filter := bson.M{
		keyStoDisable: filterA.FilterDisable,
	}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyStoBranch] = filterA.FilterBranch
	}
	if filterA.FilterCategory != "" {
		filter[keyStoCategory] = filterA.FilterCategory
	}
	if filterA.FilterName != "" {
		filter[keyStoName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterName),
		}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyStoCategory, -1}})
	opts.SetLimit(500)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar stock dari database (FindStock)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.StockResponseMinList{}, apiErr
	}

	stockList := dto.StockResponseMinList{}
	if err = cursor.All(ctx, &stockList); err != nil {
		logger.Error("Gagal decode stockList cursor ke objek slice (FindStock)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.StockResponseMinList{}, apiErr
	}

	return stockList, nil
}

func (s *stockDao) ChangeQtyStock(filterA dto.FilterIDBranch, data dto.StockChange) (*dto.Stock, rest_err.APIError) {
	coll := db.Db.Collection(keyStoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyStoID:     filterA.FilterID,
		keyStoBranch: strings.ToUpper(filterA.FilterBranch),
	}

	// Jika qty minus (decrement) beri filter qty agar tidak mengurangi sampai dengan minus
	if data.Qty < 0 {
		// cari nilai positifnya
		positive := math.Abs(float64(data.Qty))
		filter[keyStoQty] = bson.M{"$gte": int(positive)}
	}

	var update bson.D
	if data.Qty < 0 {
		// Minus , lakukan decrement
		update = bson.D{
			{"$set", bson.M{keyStoUpdatedAt: time.Now().Unix()}},
			{"$inc", bson.M{keyStoQty: data.Qty}},
			{"$push", bson.M{keyStoDecrement: data}},
		}
	} else {
		// Increment
		update = bson.D{
			{"$set", bson.M{keyStoUpdatedAt: time.Now().Unix()}},
			{"$inc", bson.M{keyStoQty: data.Qty}},
			{"$push", bson.M{keyStoIncrement: data}},
		}
	}

	var stock dto.Stock
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&stock); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Stock tidak diupdate : validasi qty (tidak mencukupi) id branch"))
		}

		logger.Error("Merubah jumlah stock gagal, (ChangeQtyStock)", err)
		apiErr := rest_err.NewInternalServerError("Merubah jumlah stock gagal", err)
		return nil, apiErr
	}

	return &stock, nil
}
