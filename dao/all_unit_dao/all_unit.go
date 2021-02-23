package all_unit

import (
	"context"
	"fmt"
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
	connectTimeout = 3
	keyAllUnitColl = "allUnit"

	keyAllID        = "_id"
	keyAllCategory  = "category"
	keyAllName      = "name"
	keyAllIP        = "ip"
	keyAllBranch    = "branch"
	keyAllCases     = "cases"
	keyAllCasesSize = "cases_size"
	keyAllPingState = "ping_state"
	keyAllLastPing  = "last_ping"
)

func AllUnitDao() AllUnitDaoAssumer {
	return &allUnitDao{}
}

type allUnitDao struct {
}

type AllUnitDaoAssumer interface {
	InsertUnit(unit dto.AllUnitRequest) (*string, rest_err.APIError)
	GetUnitByID(unitID string) (*dto.AllUnitResponse, rest_err.APIError)
	EditUnit(unitID string, unitRequest dto.AllUnitEditRequest) (*dto.AllUnitResponse, rest_err.APIError)
	DeleteUnit(unitID string) rest_err.APIError
	FindUnit(filter dto.AllUnitFilter) (dto.AllUnitResponseList, rest_err.APIError)

	//insertCase
	//deleteCase
	//insertPing
}

func (u *allUnitDao) InsertUnit(unit dto.AllUnitRequest) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyAllUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	unit.Name = strings.ToUpper(unit.Name)

	insertDoc := bson.D{
		{keyAllID, unit.ID},
		{keyAllCategory, unit.Category},
		{keyAllBranch, unit.Branch},
		{keyAllName, unit.Name},
		{keyAllIP, unit.IP},
		{keyAllCases, []string{}},
		{keyAllCasesSize, 0},
		{keyAllPingState, []dto.PingState{}},
		{keyAllLastPing, []string{}},
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan unit ke database", err)
		logger.Error("Gagal menyimpan unit ke database, InsertUnit", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (u *allUnitDao) GetUnitByID(unitID string) (*dto.AllUnitResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyAllUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	var unit dto.AllUnitResponse
	opts := options.FindOne()

	if err := coll.FindOne(ctx, bson.M{keyAllID: unitID}, opts).Decode(&unit); err != nil {

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

func (u *allUnitDao) EditUnit(unitID string, unitRequest dto.AllUnitEditRequest) (*dto.AllUnitResponse, rest_err.APIError) {
	coll := db.Db.Collection(keyAllUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	unitRequest.Name = strings.ToUpper(unitRequest.Name)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyAllID: unitID,
	}

	update := bson.M{
		"$set": bson.M{
			keyAllName:     unitRequest.Name,
			keyAllCategory: unitRequest.Category,
			keyAllBranch:   unitRequest.Branch,
			keyAllIP:       unitRequest.IP,
		},
	}

	var unit dto.AllUnitResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&unit); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, rest_err.NewBadRequestError("Unit tidak diupdate karena ID atau timestamp tidak valid")
		}

		logger.Error("Gagal mendapatkan unit dari database (EditUnit)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan unit dari database", err)
		return nil, apiErr
	}

	return &unit, nil
}

func (u *allUnitDao) DeleteUnit(unitID string) rest_err.APIError {
	coll := db.Db.Collection(keyAllUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyAllID: unitID,
	}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return rest_err.NewBadRequestError("Unit gagal dihapus, dokumen tidak ditemukan")
		}

		logger.Error("Gagal menghapus unit dari database (DeleteUnit)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan unit dari database", err)
		return apiErr
	}

	if result.DeletedCount == 0 {
		return rest_err.NewBadRequestError("Unit gagal dihapus, dokumen tidak ditemukan")
	}

	return nil
}

func (u *allUnitDao) FindUnit(filterInput dto.AllUnitFilter) (dto.AllUnitResponseList, rest_err.APIError) {
	coll := db.Db.Collection(keyAllUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterInput.Name = strings.ToUpper(filterInput.Name)
	filterInput.Category = strings.ToUpper(filterInput.Category)

	filter := bson.M{
		keyAllBranch: filterInput.Branch,
	}
	if filterInput.Category != "" {
		filter[keyAllCategory] = filterInput.Category
	}
	if filterInput.Name != "" {
		filter[keyAllName] = filterInput.Name
	}
	if filterInput.IP != "" {
		filter[keyAllIP] = filterInput.IP
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyAllName, 1}})
	sortCursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan unit dari database (FindUnit)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.AllUnitResponseList{}, apiErr
	}

	units := dto.AllUnitResponseList{}
	if err = sortCursor.All(ctx, &units); err != nil {
		logger.Error("Gagal decode unitsCursor ke objek slice (FindUnit)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.AllUnitResponseList{}, apiErr
	}

	return units, nil
}
