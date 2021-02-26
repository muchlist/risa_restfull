package gen_unit_dao

import (
	"context"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	connectTimeout = 3
	keyGenUnitColl = "genUnit"

	keyGenID        = "_id"
	keyGenCategory  = "category"
	keyGenName      = "name"
	keyGenIP        = "ip"
	keyGenBranch    = "branch"
	keyGenCases     = "cases"
	keyGenCasesSize = "cases_size"
	keyGenPingState = "ping_state"
	keyGenLastPing  = "last_ping"
	keyGenDisable   = "disable"

	keyCaseID   = "case_id"
	keyCaseNote = "case_note"
)

func NewGenUnitDao() GenUnitDaoAssumer {
	return &genUnitDao{}
}

type genUnitDao struct {
}

type GenUnitDaoAssumer interface {
	InsertUnit(unit dto.GenUnitRequest) (*string, rest_err.APIError)
	EditUnit(unitID string, unitRequest dto.GenUnitEditRequest) (*dto.GenUnitResponse, rest_err.APIError)
	DeleteUnit(unitID string) rest_err.APIError
	InsertCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	DeleteCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	DisableUnit(unitID string, value bool) (*dto.GenUnitResponse, rest_err.APIError)
	//insertPing

	GetUnitByID(unitID string, branchSpecific string) (*dto.GenUnitResponse, rest_err.APIError)
	FindUnit(filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError)
}

func (u *genUnitDao) InsertUnit(unit dto.GenUnitRequest) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	unit.Name = strings.ToUpper(unit.Name)

	insertDoc := bson.M{
		keyGenID:        unit.ID,
		keyGenCategory:  unit.Category,
		keyGenBranch:    unit.Branch,
		keyGenName:      unit.Name,
		keyGenIP:        unit.IP,
		keyGenCases:     []string{},
		keyGenCasesSize: 0,
		keyGenPingState: []dto.PingState{},
		keyGenLastPing:  "",
		keyGenDisable:   false,
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan unit ke database", err)
		logger.Error("Gagal menyimpan unit ke database, InsertUnit", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(string)

	return &insertID, nil
}

func (u *genUnitDao) EditUnit(unitID string, unitRequest dto.GenUnitEditRequest) (*dto.GenUnitResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	unitRequest.Name = strings.ToUpper(unitRequest.Name)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyGenID: unitID,
	}

	update := bson.M{
		"$set": bson.M{
			keyGenName:     unitRequest.Name,
			keyGenCategory: unitRequest.Category,
			keyGenBranch:   unitRequest.Branch,
			keyGenIP:       unitRequest.IP,
		},
	}

	var unit dto.GenUnitResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&unit); err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error("Gagal mengedit unit dari database (EditUnit)", err)
			return nil, rest_err.NewBadRequestError("Unit tidak diupdate karena ID tidak valid")
		}

		logger.Error("Gagal mengedit unit dari database (EditUnit)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mengedit unit dari database", err)
		return nil, apiErr
	}

	return &unit, nil
}

func (u *genUnitDao) DeleteUnit(unitID string) rest_err.APIError {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyGenID: unitID,
	}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest_err.NewBadRequestError("Unit gagal dihapus, dokumen tidak ditemukan")
		}

		logger.Error("Gagal menghapus unit dari database (DeleteUnit)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menghapus unit dari database", err)
		return apiErr
	}

	if result.DeletedCount == 0 {
		return rest_err.NewBadRequestError("Unit gagal dihapus, dokumen tidak ditemukan")
	}

	return nil
}

func (u *genUnitDao) DisableUnit(unitID string, value bool) (*dto.GenUnitResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyGenID: unitID,
	}

	update := bson.M{
		"$set": bson.M{
			keyGenDisable: value,
		},
	}

	var unit dto.GenUnitResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&unit); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Unit tidak diupdate karena ID tidak valid")
		}

		logger.Error("Gagal mendapatkan unit dari database (DisableUnit)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan unit dari database", err)
		return nil, apiErr
	}

	return &unit, nil
}

func (u *genUnitDao) InsertCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	payload.FilterBranch = strings.ToUpper(payload.FilterBranch)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyGenID:     payload.UnitID,
		keyGenBranch: payload.FilterBranch,
	}

	update := bson.M{
		"$push": bson.M{
			keyGenCases: bson.M{keyCaseID: payload.CaseID, keyCaseNote: payload.CaseNote},
		},
		"$inc": bson.M{
			keyGenCasesSize: 1,
		},
	}

	var unit dto.GenUnitResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&unit); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Unit tidak diupdate karena ID atau timestamp tidak valid")
		}

		logger.Error("Gagal mendapatkan unit dari database (InsertCase)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan unit dari database", err)
		return nil, apiErr
	}

	return &unit, nil
}

func (u *genUnitDao) DeleteCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	payload.FilterBranch = strings.ToUpper(payload.FilterBranch)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyGenID:     payload.UnitID,
		keyGenBranch: payload.FilterBranch,
	}

	update := bson.M{
		"$pull": bson.M{
			keyGenCases: bson.M{keyCaseID: payload.CaseID},
		},
		"$inc": bson.M{
			keyGenCasesSize: -1,
		},
	}

	var unit dto.GenUnitResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&unit); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Unit tidak diupdate karena ID atau timestamp tidak valid")
		}

		logger.Error("Gagal mendapatkan unit dari database (InsertCase)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan unit dari database", err)
		return nil, apiErr
	}

	return &unit, nil
}

func (u *genUnitDao) GetUnitByID(unitID string, branchIfSpecific string) (*dto.GenUnitResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	branchIfSpecific = strings.ToUpper(branchIfSpecific)

	var unit dto.GenUnitResponse
	opts := options.FindOne()

	filter := bson.M{keyGenID: unitID}

	if branchIfSpecific != "" {
		filter[keyGenBranch] = branchIfSpecific
	}

	if err := coll.FindOne(ctx, filter, opts).Decode(&unit); err != nil {

		if err == mongo.ErrNoDocuments {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Unit dengan ID %s tidak ditemukan", unitID))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan unit dari database (GetUnitByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan unit dari database", err)
		return nil, apiErr
	}

	return &unit, nil
}

func (u *genUnitDao) FindUnit(filterInput dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError) {
	coll := db.Db.Collection(keyGenUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterInput.Name = strings.ToUpper(filterInput.Name)
	filterInput.Category = strings.ToUpper(filterInput.Category)

	filter := bson.M{
		keyGenBranch:  filterInput.Branch,
		keyGenDisable: false,
	}
	if filterInput.Category != "" {
		filter[keyGenCategory] = filterInput.Category
	}
	if filterInput.Name != "" {
		filter[keyGenName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterInput.Name),
		} // {'$regex': f'.*{cctv_name.upper()}.*'}
	}
	if filterInput.IP != "" {
		filter[keyGenIP] = filterInput.IP
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyGenName, 1}})
	opts.SetLimit(200)
	sortCursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan unit dari database (FindUnit)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.GenUnitResponseList{}, apiErr
	}

	units := dto.GenUnitResponseList{}
	if err = sortCursor.All(ctx, &units); err != nil {
		logger.Error("Gagal decode unitsCursor ke objek slice (FindUnit)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.GenUnitResponseList{}, apiErr
	}

	return units, nil
}
