package db

import (
	"context"
	"github.com/muchlist/erru_utils_go/logger"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout = 10
	mongoURLGetKey = "MONGO_DB_URL_RISA"
	databaseName   = "risa_go"
)

var (
	// Db objek sebagai database objek
	Db       *mongo.Database
	mongoURL = "mongodb://localhost:27017"
)

// Init menginisiasi database
// responsenya digunakan untuk memutus koneksi apabila main program dihentikan
func Init() (*mongo.Client, context.Context, context.CancelFunc) {

	if os.Getenv(mongoURLGetKey) != "" {
		mongoURL = os.Getenv(mongoURLGetKey)
	}

	Client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.Error("gagal membuat client mongodb", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = Client.Connect(ctx)
	if err != nil {
		logger.Error("gagal menghubungkan koneksi mongodb", err)
		panic(err)
	}

	logger.Info("database berhasil terkoneksi")
	Db = Client.Database(databaseName)

	return Client, ctx, cancel
}
