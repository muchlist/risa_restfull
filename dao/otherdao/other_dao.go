package otherdao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout     = 3
	keyOtherCollection = "other"

	keyOtherID          = "_id"
	keyOtherName        = "name"
	keyOtherCreatedAt   = "created_at"
	keyOtherUpdatedAt   = "updated_at"
	keyOtherUpdatedBy   = "updated_by"
	keyOtherUpdatedByID = "updated_by_id"
	keyOtherBranch      = "branch"
	keyOtherDisable     = "disable"

	keyOtherDetail      = "detail"
	keyOtherDivision    = "division"
	keyOtherSubCategory = "sub_category"

	keyOtherIP              = "ip"
	keyOtherInventoryNumber = "inventory_number"
	keyOtherLocation        = "location"
	keyOtherLocationLat     = "location_lat"
	keyOtherLocationLon     = "location_lon"
	keyOtherDate            = "date"
	keyOtherTag             = "tag"
	keyOtherImage           = "image"
	keyOtherBrand           = "brand"
	keyOtherType            = "type"
	keyOtherNote            = "note"
)

func NewOtherDao() OtherDaoAssumer {
	return &otherDao{}
}

type otherDao struct{}

func (c *otherDao) InsertOther(ctx context.Context, input dto.Other) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyOtherCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)
	input.SubCategory = strings.ToUpper(input.SubCategory)
	if input.Tag == nil {
		input.Tag = []string{}
	}
	input.Disable = false

	result, err := coll.InsertOne(ctxt, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError(fmt.Sprintf("Gagal menyimpan %s ke database", input.SubCategory), err)
		logger.Error(fmt.Sprintf("Gagal menyimpan %s ke database (InsertOther)", input.SubCategory), err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *otherDao) EditOther(ctx context.Context, input dto.OtherEdit) (*dto.Other, rest_err.APIError) {
	coll := db.DB.Collection(keyOtherCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.FilterBranch = strings.ToUpper(input.FilterBranch)
	input.FilterSubCategory = strings.ToUpper(input.FilterSubCategory)
	if input.Tag == nil {
		input.Tag = []string{}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyOtherID:          input.ID,
		keyOtherBranch:      input.FilterBranch,
		keyOtherSubCategory: input.FilterSubCategory,
		keyOtherUpdatedAt:   input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyOtherName:        input.Name,
			keyOtherUpdatedAt:   input.UpdatedAt,
			keyOtherUpdatedBy:   input.UpdatedBy,
			keyOtherUpdatedByID: input.UpdatedByID,

			keyOtherIP:              input.IP,
			keyOtherInventoryNumber: input.InventoryNumber,
			keyOtherLocation:        input.Location,
			keyOtherLocationLat:     input.LocationLat,
			keyOtherLocationLon:     input.LocationLon,

			keyOtherDivision: input.Division,
			keyOtherDetail:   input.Detail,

			keyOtherDate:  input.Date,
			keyOtherTag:   input.Tag,
			keyOtherBrand: input.Brand,
			keyOtherType:  input.Type,
			keyOtherNote:  input.Note,
		},
	}

	var other dto.Other
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&other); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("%s tidak diupdate : validasi id branch category timestamp", input.FilterSubCategory))
		}

		logger.Error(fmt.Sprintf("Gagal mendapatkan %s dari database (EditOther)", input.FilterSubCategory), err)
		apiErr := rest_err.NewInternalServerError(fmt.Sprintf("Gagal mendapatkan %s dari database", input.FilterSubCategory), err)
		return nil, apiErr
	}

	return &other, nil
}

func (c *otherDao) DeleteOther(ctx context.Context, input dto.FilterIDBranchCategoryCreateGte) (*dto.Other, rest_err.APIError) {
	coll := db.DB.Collection(keyOtherCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyOtherID:          input.FilterID,
		keyOtherBranch:      strings.ToUpper(input.FilterBranch),
		keyOtherSubCategory: strings.ToUpper(input.FilterSubCategory),
		keyOtherCreatedAt:   bson.M{"$gte": input.FilterCreateGTE},
	}

	var other dto.Other
	err := coll.FindOneAndDelete(ctxt, filter).Decode(&other)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Data tidak diupdate : validasi id branch category time_reach")
		}

		logger.Error("Gagal menghapus data dari database (DeleteOther)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus data dari database", err)
		return nil, apiErr
	}

	return &other, nil
}

// DisableOther if value true , other will disabled
func (c *otherDao) DisableOther(ctx context.Context, otherID primitive.ObjectID, user mjwt.CustomClaim, subCategory string, value bool) (*dto.Other, rest_err.APIError) {
	coll := db.DB.Collection(keyOtherCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyOtherID:          otherID,
		keyOtherBranch:      strings.ToUpper(user.Branch),
		keyOtherSubCategory: strings.ToUpper(subCategory),
	}

	update := bson.M{
		"$set": bson.M{
			keyOtherDisable:     value,
			keyOtherUpdatedAt:   time.Now().Unix(),
			keyOtherUpdatedByID: user.Identity,
			keyOtherUpdatedBy:   user.Name,
		},
	}

	var other dto.Other
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&other); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Data tidak diupdate : validasi id branch category")
		}

		logger.Error("Gagal mendisable data dari database (DisableOther)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendisable data dari database", err)
		return nil, apiErr
	}

	return &other, nil
}

func (c *otherDao) UploadImage(ctx context.Context, pcID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Other, rest_err.APIError) {
	coll := db.DB.Collection(keyOtherCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyOtherID:     pcID,
		keyOtherBranch: strings.ToUpper(filterBranch),
	}
	update := bson.M{
		"$set": bson.M{
			keyOtherImage: imagePath,
		},
	}

	var other dto.Other
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&other); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, data dengan id %s tidak ditemukan", pcID.Hex()))
		}

		logger.Error("Memasukkan path image data ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image data ke db gagal", err)
		return nil, apiErr
	}

	return &other, nil
}

func (c *otherDao) GetOtherByID(ctx context.Context, pcID primitive.ObjectID, branchIfSpecific string) (*dto.Other, rest_err.APIError) {
	coll := db.DB.Collection(keyOtherCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyOtherID: pcID}
	if branchIfSpecific != "" {
		filter[keyOtherBranch] = strings.ToUpper(branchIfSpecific)
	}

	var other dto.Other
	if err := coll.FindOne(ctxt, filter).Decode(&other); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Data dengan ID %s tidak ditemukan", pcID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan data dari database (GetOtherByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan data dari database", err)
		return nil, apiErr
	}

	return &other, nil
}

func (c *otherDao) FindOther(ctx context.Context, filterA dto.FilterOther) (dto.OtherResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyOtherCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterSubCategory = strings.ToUpper(filterA.FilterSubCategory)
	filterA.FilterName = strings.ToUpper(filterA.FilterName)
	filterA.FilterDivision = strings.ToUpper(filterA.FilterDivision)
	filterA.FilterLocation = strings.ToUpper(filterA.FilterLocation)

	// filter
	filter := bson.M{
		keyOtherDisable: filterA.FilterDisable,
	}
	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyOtherBranch] = filterA.FilterBranch
	}

	// sub category
	if filterA.FilterSubCategory != "" {
		// cek subkategori jika multi subcategory (pisah dengan koma)
		if strings.Contains(filterA.FilterSubCategory, ",") {
			subCategories := strings.Split(filterA.FilterSubCategory, ",")
			filter[keyOtherSubCategory] = bson.M{"$in": subCategories}
		} else {
			filter[keyOtherSubCategory] = filterA.FilterSubCategory
		}
	}

	if filterA.FilterName != "" {
		filter[keyOtherName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterName),
		}
	}
	if filterA.FilterLocation != "" {
		filter[keyOtherLocation] = filterA.FilterLocation
	}
	if filterA.FilterDivision != "" {
		filter[keyOtherDivision] = filterA.FilterDivision
	}
	if filterA.FilterIP != "" {
		filter[keyOtherIP] = filterA.FilterIP
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyOtherLocation, Value: -1}, {Key: keyOtherDivision, Value: -1}}) //nolint:govet
	opts.SetLimit(500)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error(fmt.Sprintf("Gagal mendapatkan daftar %s dari database (FindOther)", filterA.FilterSubCategory), err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.OtherResponseMinList{}, apiErr
	}

	otherList := dto.OtherResponseMinList{}
	if err = cursor.All(ctxt, &otherList); err != nil {
		logger.Error("Gagal decode otherList cursor ke objek slice (FindOther)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.OtherResponseMinList{}, apiErr
	}

	return otherList, nil
}
