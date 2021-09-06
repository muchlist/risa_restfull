package serverfiledao

import (
	"context"
	"errors"
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
	keyCollection  = "serverConf"

	keyID        = "_id"
	keyTitle     = "title"
	keyUpdatedAt = "updated_at"
	keyUpdatedBy = "updated_by"
	keyBranch    = "branch"
	keyDiff      = "diff"
	keyImage     = "image"
	keyNote      = "note"
)

func NewServerFileDao() ServerFileAssumer {
	return &serverFileDao{}
}

type serverFileDao struct {
}

type ServerFileAssumer interface {
	Insert(input dto.ServerFile) (*string, rest_err.APIError)
	Delete(input dto.FilterIDBranchCreateGte) (*dto.ServerFile, rest_err.APIError)
	UploadImage(itemID primitive.ObjectID, imagePath string, filterBranch string) (*dto.ServerFile, rest_err.APIError)
	GetByID(itemID primitive.ObjectID, branchIfSpecific string) (*dto.ServerFile, rest_err.APIError)
	Find(branch string, start int64, end int64) ([]dto.ServerFile, rest_err.APIError)
}

func (s *serverFileDao) Insert(input dto.ServerFile) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Title = strings.ToUpper(input.Title)
	input.Branch = strings.ToUpper(input.Branch)

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan item ke database", err)
		logger.Error("Gagal menyimpan item ke database, (Insert)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (s *serverFileDao) Delete(input dto.FilterIDBranchCreateGte) (*dto.ServerFile, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyID:        input.FilterID,
		keyBranch:    input.FilterBranch,
		keyUpdatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var sf dto.ServerFile
	err := coll.FindOneAndDelete(ctx, filter).Decode(&sf)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Item tidak diupdate : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus item dari database (Delete)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan item dari database", err)
		return nil, apiErr
	}

	return &sf, nil
}

func (s *serverFileDao) UploadImage(fileID primitive.ObjectID, imagePath string, filterBranch string) (*dto.ServerFile, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:     fileID,
		keyBranch: strings.ToUpper(filterBranch),
	}
	update := bson.M{
		"$set": bson.M{
			keyImage: imagePath,
		},
	}

	var sf dto.ServerFile
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&sf); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, item dengan id %s tidak ditemukan", fileID.Hex()))
		}

		logger.Error("Memasukkan path image utem ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image item ke db gagal", err)
		return nil, apiErr
	}

	return &sf, nil
}

func (s *serverFileDao) GetByID(fileID primitive.ObjectID, branchIfSpecific string) (*dto.ServerFile, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: fileID}
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var sf dto.ServerFile
	if err := coll.FindOne(ctx, filter).Decode(&sf); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("item dengan ID %s tidak ditemukan", fileID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan item dari database (GetByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan item dari database", err)
		return nil, apiErr
	}

	return &sf, nil
}

func (s *serverFileDao) Find(branch string, start int64, end int64) ([]dto.ServerFile, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyBranch: strings.ToUpper(branch),
	}

	// option range
	if start != 0 {
		filter[keyUpdatedAt] = bson.M{"$gte": start}
	}
	if end != 0 {
		filter[keyUpdatedAt] = bson.M{"$lte": end}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyID, Value: -1}})
	opts.SetLimit(100)

	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		logger.Error("Gagal mendapatkan daftar item dari database (Find)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.ServerFile{}, apiErr
	}

	serverFiles := make([]dto.ServerFile, 0)
	if err = cursor.All(ctx, &serverFiles); err != nil {
		logger.Error("Gagal decode serverFile cursor ke objek slice (Find)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.ServerFile{}, apiErr
	}

	return serverFiles, nil
}
