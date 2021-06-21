package speedtestdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	connectTimeout  = 3
	keySpCollection = "speedTest"

	keySpID       = "id"
	keySpTime     = "time"
	keySpLatency  = "latency_ms"
	keySpUpload   = "upload"
	keySpDownload = "download"
)

func NewSpeedTestDao() SpedTestDaoAssumer {
	return &speedTestDao{}
}

type speedTestDao struct {
}

type SpedTestDaoAssumer interface {
	InsertSpeed(input dto.SpeedTest) (*string, rest_err.APIError)
	RetrieveSpeed() (dto.SpeedTestList, rest_err.APIError)
}

func (s *speedTestDao) InsertSpeed(input dto.SpeedTest) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keySpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan data speed ke database", err)
		logger.Error("Gagal menyimpan data speed ke database, (InsertCheck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (s *speedTestDao) RetrieveSpeed() (dto.SpeedTestList, rest_err.APIError) {
	coll := db.DB.Collection(keySpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		logger.Error("Gagal mendapatkan daftar check dari database (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.SpeedTestList{}, apiErr
	}

	speedList := dto.SpeedTestList{}
	if err = cursor.All(ctx, &speedList); err != nil {
		logger.Error("Gagal decode speedList cursor ke objek slice (FindCheck)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.SpeedTestList{}, apiErr
	}

	return speedList, nil
}
