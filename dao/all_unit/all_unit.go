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
	//GetUserByIDWithPassword(userID string) (*dto.User, rest_err.APIError)
	//FindUser() (dto.UserResponseList, rest_err.APIError)
	//CheckIDAvailable(email string) (bool, rest_err.APIError)
	//EditUser(userID string, userRequest dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError)
	//DeleteUser(userID string) rest_err.APIError
	//PutAvatar(userID string, avatar string) (*dto.UserResponse, rest_err.APIError)
	//ChangePassword(data dto.UserChangePasswordRequest) rest_err.APIError
}

func (u *allUnitDao) InsertUnit(unit dto.AllUnitRequest) (*string, rest_err.APIError) {
	coll := db.Db.Collection(keyAllUnitColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

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
