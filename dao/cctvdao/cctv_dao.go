package cctvdao

import (
	"context"
	"errors"
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
	keyCtvUpdatedAt   = "updated_at"
	keyCtvUpdatedBy   = "updated_by"
	keyCtvUpdatedByID = "updated_by_id"
	keyCtvBranch      = "branch"
	keyCtvDisable     = "disable"
	keyCtvDisVendor   = "dis_vendor"

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

func (c *cctvDao) InsertCctv(ctx context.Context, input dto.Cctv) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCtvCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)
	if input.Tag == nil {
		input.Tag = []string{}
	}
	input.Disable = false

	result, err := coll.InsertOne(ctxt, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan cctv ke database", err)
		logger.Error("Gagal menyimpan cctv ke database, (InsertCctv)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *cctvDao) EditCctv(ctx context.Context, input dto.CctvEdit) (*dto.Cctv, rest_err.APIError) {
	coll := db.DB.Collection(keyCtvCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
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

			keyCtvDate:      input.Date,
			keyCtvTag:       input.Tag,
			keyCtvBrand:     input.Brand,
			keyCtvType:      input.Type,
			keyCtvNote:      input.Note,
			keyCtvDisVendor: input.DisVendor,
		},
	}

	var cctv dto.Cctv
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&cctv); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Cctv tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan cctv dari database (EditCctv)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) DeleteCctv(ctx context.Context, input dto.FilterIDBranchCreateGte) (*dto.Cctv, rest_err.APIError) {
	coll := db.DB.Collection(keyCtvCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyCtvID:        input.FilterID,
		keyCtvBranch:    input.FilterBranch,
		keyCtvCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var cctv dto.Cctv
	err := coll.FindOneAndDelete(ctxt, filter).Decode(&cctv)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Cctv tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus cctv dari database (DeleteCctv)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

// DisableCctv if value true , cctv will disabled
func (c *cctvDao) DisableCctv(ctx context.Context, cctvID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError) {
	coll := db.DB.Collection(keyCtvCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
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
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&cctv); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Cctv tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendisable cctv dari database (DisableCctv)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendisable cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) UploadImage(ctx context.Context, cctvID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Cctv, rest_err.APIError) {
	coll := db.DB.Collection(keyCtvCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
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
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&cctv); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, cctv dengan id %s tidak ditemukan", cctvID.Hex()))
		}

		logger.Error("Memasukkan path image cctv ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image cctv ke db gagal", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) GetCctvByID(ctx context.Context, cctvID primitive.ObjectID, branchIfSpecific string) (*dto.Cctv, rest_err.APIError) {
	coll := db.DB.Collection(keyCtvCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyCtvID: cctvID}
	if branchIfSpecific != "" {
		filter[keyCtvBranch] = strings.ToUpper(branchIfSpecific)
	}

	var cctv dto.Cctv
	if err := coll.FindOne(ctxt, filter).Decode(&cctv); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Cctv dengan ID %s tidak ditemukan", cctvID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan cctv dari database (GetCctvByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan cctv dari database", err)
		return nil, apiErr
	}

	return &cctv, nil
}

func (c *cctvDao) FindCctv(ctx context.Context, filterA dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyCtvCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterName = strings.ToUpper(filterA.FilterName)

	// filter
	filter := bson.M{
		keyCtvDisable: filterA.FilterDisable,
	}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyCtvBranch] = filterA.FilterBranch
	}
	if filterA.FilterName != "" {
		filter[keyCtvName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterName),
		}
	}
	if filterA.FilterLocation != "" {
		filter[keyCtvLocation] = filterA.FilterLocation
	}
	if filterA.FilterIP != "" {
		filter[keyCtvIP] = filterA.FilterIP
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyCtvLocation, -1}, {keyCtvName, 1}}) //nolint:govet
	opts.SetLimit(500)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar cctv dari database (FindCctv)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CctvResponseMinList{}, apiErr
	}

	cctvList := dto.CctvResponseMinList{}
	if err = cursor.All(ctxt, &cctvList); err != nil {
		logger.Error("Gagal decode cctvList cursor ke objek slice (FindCctv)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.CctvResponseMinList{}, apiErr
	}

	return cctvList, nil
}
