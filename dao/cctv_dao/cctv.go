package cctv_dao

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
	keyCtvDisable     = "disable"

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
	EditCctv(input dto.CctvEdit) (*dto.Cctv, rest_err.APIError)
	DeleteCctv(input dto.FilterIDBranchTime) (*dto.Cctv, rest_err.APIError)
	DisableCctv(cctvID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError)
	UploadImage(cctvID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Cctv, rest_err.APIError)

	GetCctvByID(cctvID primitive.ObjectID, branchIfSpecific string) (*dto.Cctv, rest_err.APIError)
	FindCctv(filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError)
}

func (c *cctvDao) InsertCctv(input dto.Cctv) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)
	if input.Tag == nil {
		input.Tag = []string{}
	}

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
		keyCtvDisable:     false,

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

func (c *cctvDao) EditCctv(input dto.CctvEdit) (*dto.Cctv, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	if input.Tag == nil {
		input.Tag = []string{}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyCtvID:        input.ID,
		keyCtvBranch:    input.FilterBranch,
		keyCtvUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyCtvName:        input.Name,
			keyCtvUpdatedAt:   input.UpdatedAt,
			keyCtvUpdatedBy:   input.UpdatedBy,
			keyCtvUpdatedByID: input.UpdatedByID,

			keyCtvIP:              input.IP,
			keyCtvInventoryNumber: input.InventoryNumber,
			keyCtvLocation:        input.Location,
			keyCtvLocationLat:     input.LocationLat,
			keyCtvLocationLon:     input.LocationLon,

			keyCtvDate:  input.Date,
			keyCtvTag:   input.Tag,
			keyCtvBrand: input.Brand,
			keyCtvType:  input.Type,
			keyCtvNote:  input.Note,
		},
	}

	var cctv dto.Cctv
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&cctv); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Cctv tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan cctv dari database (EditCctv)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) DeleteCctv(input dto.FilterIDBranchTime) (*dto.Cctv, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyCtvID:        input.ID,
		keyCtvBranch:    input.Branch,
		keyCtvCreatedAt: bson.M{"$gte": input.Time},
	}

	var cctv dto.Cctv
	err := coll.FindOneAndDelete(ctx, filter).Decode(&cctv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Cctv tidak diupdate : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus cctv dari database (DeleteCctv)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

// DisableCctv if value true , cctv will disabled
func (c *cctvDao) DisableCctv(cctvID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyCtvID:     cctvID,
		keyCtvBranch: user.Branch,
	}

	update := bson.M{
		"$set": bson.M{
			keyCtvDisable:     value,
			keyCtvUpdatedAt:   time.Now().Unix(),
			keyCtvUpdatedByID: user.Identity,
			keyCtvUpdatedBy:   user.Name,
		},
	}

	var cctv dto.Cctv
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&cctv); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Cctv tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendapatkan cctv dari database (DisableCctv)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) UploadImage(cctvID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Cctv, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyCtvID:     cctvID,
		keyCtvBranch: strings.ToUpper(filterBranch),
	}
	update := bson.M{
		"$set": bson.M{
			keyCtvImage: imagePath,
		},
	}

	var cctv dto.Cctv
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&cctv); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, cctv dengan id %s tidak ditemukan", cctvID.Hex()))
		}

		logger.Error("Memasukkan path image cctv ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image cctv ke db gagal", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) GetCctvByID(cctvID primitive.ObjectID, branchIfSpecific string) (*dto.Cctv, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyCtvID: cctvID}
	if branchIfSpecific != "" {
		filter[keyCtvBranch] = strings.ToUpper(branchIfSpecific)
	}

	var cctv dto.Cctv
	if err := coll.FindOne(ctx, filter).Decode(&cctv); err != nil {

		if err == mongo.ErrNoDocuments {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Cctv dengan ID %s tidak ditemukan", cctvID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan cctv dari database (GetCctvByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) FindCctv(filterA dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError) {
	coll := db.Db.Collection(keyCtvCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.Branch = strings.ToUpper(filterA.Branch)
	filterA.Name = strings.ToUpper(filterA.Name)

	// filter
	filter := bson.M{
		keyCtvDisable: filterA.Disable,
	}

	// filter condition
	if filterA.Branch != "" {
		filter[keyCtvBranch] = filterA.Branch
	}
	if filterA.Name != "" {
		filter[keyCtvName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.Name),
		}
	}
	if filterA.Location != "" {
		filter[keyCtvLocation] = filterA.Location
	}
	if filterA.IP != "" {
		filter[keyCtvIP] = filterA.IP
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyCtvLocation, -1}})
	opts.SetLimit(500)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar cctv dari database (FindCctv)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CctvResponseMinList{}, apiErr
	}

	cctvList := dto.CctvResponseMinList{}
	if err = cursor.All(ctx, &cctvList); err != nil {
		logger.Error("Gagal decode cctvList cursor ke objek slice (FindCctv)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CctvResponseMinList{}, apiErr
	}

	return cctvList, nil
}
