package computerdao

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
	connectTimeout  = 3
	keyPCCollection = "computer"

	keyPCID          = "_id"
	keyPCName        = "name"
	keyPCCreatedAt   = "created_at"
	keyPCUpdatedAt   = "updated_at"
	keyPCUpdatedBy   = "updated_by"
	keyPCUpdatedByID = "updated_by_id"
	keyPCBranch      = "branch"
	keyPCDisable     = "disable"

	keyPCHostname       = "hostname"
	keyPCDivision       = "division"
	keyPCSeatManagement = "seat_management"
	keyPCOS             = "os"
	keyPCProcessor      = "processor"
	keyPCRam            = "ram"
	keyPCHardisk        = "hardisk"

	keyPCIP              = "ip"
	keyPCInventoryNumber = "inventory_number"
	keyPCLocation        = "location"
	keyPCLocationLat     = "location_lat"
	keyPCLocationLon     = "location_lon"
	keyPCDate            = "date"
	keyPCTag             = "tag"
	keyPCImage           = "image"
	keyPCBrand           = "brand"
	keyPCType            = "type"
	keyPCNote            = "note"
)

func NewComputerDao() ComputerDaoAssumer {
	return &computerDao{}
}

type computerDao struct {
}

type ComputerDaoAssumer interface {
	InsertPc(input dto.Computer) (*string, rest_err.APIError)
	EditPc(input dto.ComputerEdit) (*dto.Computer, rest_err.APIError)
	DeletePc(input dto.FilterIDBranchCreateGte) (*dto.Computer, rest_err.APIError)
	DisablePc(pcID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.Computer, rest_err.APIError)
	UploadImage(pcID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Computer, rest_err.APIError)

	GetPcByID(pcID primitive.ObjectID, branchIfSpecific string) (*dto.Computer, rest_err.APIError)
	FindPc(filter dto.FilterComputer) (dto.ComputerResponseMinList, rest_err.APIError)
}

func (c *computerDao) InsertPc(input dto.Computer) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyPCCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Branch = strings.ToUpper(input.Branch)
	input.Location = strings.ToUpper(input.Location)
	input.Division = strings.ToUpper(input.Division)
	if input.Tag == nil {
		input.Tag = []string{}
	}
	input.Disable = false

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan pc ke database", err)
		logger.Error("Gagal menyimpan pc ke database, (InsertPc)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *computerDao) EditPc(input dto.ComputerEdit) (*dto.Computer, rest_err.APIError) {
	coll := db.DB.Collection(keyPCCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.Location = strings.ToUpper(input.Location)
	input.Division = strings.ToUpper(input.Division)
	if input.Tag == nil {
		input.Tag = []string{}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyPCID:        input.ID,
		keyPCBranch:    input.FilterBranch,
		keyPCUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyPCName:        input.Name,
			keyPCUpdatedAt:   input.UpdatedAt,
			keyPCUpdatedBy:   input.UpdatedBy,
			keyPCUpdatedByID: input.UpdatedByID,

			keyPCIP:              input.IP,
			keyPCInventoryNumber: input.InventoryNumber,
			keyPCLocation:        input.Location,
			keyPCLocationLat:     input.LocationLat,
			keyPCLocationLon:     input.LocationLon,

			keyPCHostname:       input.Hostname,
			keyPCDivision:       input.Division,
			keyPCSeatManagement: input.SeatManagement,
			keyPCOS:             input.OS,
			keyPCProcessor:      input.Processor,
			keyPCRam:            input.Ram,
			keyPCHardisk:        input.Hardisk,

			keyPCDate:  input.Date,
			keyPCTag:   input.Tag,
			keyPCBrand: input.Brand,
			keyPCType:  input.Type,
			keyPCNote:  input.Note,
		},
	}

	var pc dto.Computer
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&pc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Pc tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan pc dari database (EditPc)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan pc dari database", err)
		return nil, apiErr
	}

	return &pc, nil
}

func (c *computerDao) DeletePc(input dto.FilterIDBranchCreateGte) (*dto.Computer, rest_err.APIError) {
	coll := db.DB.Collection(keyPCCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyPCID:        input.FilterID,
		keyPCBranch:    input.FilterBranch,
		keyPCCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var pc dto.Computer
	err := coll.FindOneAndDelete(ctx, filter).Decode(&pc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Pc tidak dihapus : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus pc dari database (DeletePc)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus pc dari database", err)
		return nil, apiErr
	}

	return &pc, nil
}

// DisablePc if value true , pc will disabled
func (c *computerDao) DisablePc(pcID primitive.ObjectID, user mjwt.CustomClaim, value bool) (*dto.Computer, rest_err.APIError) {
	coll := db.DB.Collection(keyPCCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyPCID:     pcID,
		keyPCBranch: user.Branch,
	}

	update := bson.M{
		"$set": bson.M{
			keyPCDisable:     value,
			keyPCUpdatedAt:   time.Now().Unix(),
			keyPCUpdatedByID: user.Identity,
			keyPCUpdatedBy:   user.Name,
		},
	}

	var pc dto.Computer
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&pc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Pc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendisable pc dari database (DisablePc)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendisable pc dari database", err)
		return nil, apiErr
	}

	return &pc, nil
}

func (c *computerDao) UploadImage(pcID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Computer, rest_err.APIError) {
	coll := db.DB.Collection(keyPCCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyPCID:     pcID,
		keyPCBranch: strings.ToUpper(filterBranch),
	}
	update := bson.M{
		"$set": bson.M{
			keyPCImage: imagePath,
		},
	}

	var pc dto.Computer
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&pc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, pc dengan id %s tidak ditemukan", pcID.Hex()))
		}

		logger.Error("Memasukkan path image pc ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image pc ke db gagal", err)
		return nil, apiErr
	}

	return &pc, nil
}

func (c *computerDao) GetPcByID(pcID primitive.ObjectID, branchIfSpecific string) (*dto.Computer, rest_err.APIError) {
	coll := db.DB.Collection(keyPCCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyPCID: pcID}
	if branchIfSpecific != "" {
		filter[keyPCBranch] = strings.ToUpper(branchIfSpecific)
	}

	var pc dto.Computer
	if err := coll.FindOne(ctx, filter).Decode(&pc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Pc dengan ID %s tidak ditemukan", pcID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan pc dari database (GetPcByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan pc dari database", err)
		return nil, apiErr
	}

	return &pc, nil
}

func (c *computerDao) FindPc(filterA dto.FilterComputer) (dto.ComputerResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyPCCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterName = strings.ToUpper(filterA.FilterName)
	filterA.FilterDivision = strings.ToUpper(filterA.FilterDivision)
	filterA.FilterLocation = strings.ToUpper(filterA.FilterLocation)

	// filter
	filter := bson.M{
		keyPCDisable: filterA.FilterDisable,
	}
	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyPCBranch] = filterA.FilterBranch
	}
	if filterA.FilterName != "" {
		filter[keyPCName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterName),
		}
	}
	if filterA.FilterLocation != "" {
		filter[keyPCLocation] = filterA.FilterLocation
	}
	if filterA.FilterDivision != "" {
		filter[keyPCDivision] = filterA.FilterDivision
	}
	if filterA.FilterIP != "" {
		filter[keyPCIP] = filterA.FilterIP
	}

	switch filterA.FilterSeatManagement {
	case 0:
		filter[keyPCSeatManagement] = false
	case 1:
		filter[keyPCSeatManagement] = true
	default:
		// do nothing
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyPCLocation, -1}, {keyPCDivision, -1}}) //nolint:govet
	opts.SetLimit(500)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar pc dari database (FindPc)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.ComputerResponseMinList{}, apiErr
	}

	pcList := dto.ComputerResponseMinList{}
	if err = cursor.All(ctx, &pcList); err != nil {
		logger.Error("Gagal decode pcList cursor ke objek slice (FindPc)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.ComputerResponseMinList{}, apiErr
	}

	return pcList, nil
}
