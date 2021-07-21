package improvedao

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
	keyImpCollection = "improve"

	keyImpID          = "_id"
	keyImpCreatedAt   = "created_at"
	keyImpUpdatedAt   = "updated_at"
	keyImpUpdatedBy   = "updated_by"
	keyImpUpdatedByID = "updated_by_id"
	keyImpBranch      = "branch"

	keyImpTitle          = "title"
	keyImpDescription    = "description"
	keyImpGoal           = "goal"
	keyImpGoalsAchieved  = "goals_achieved"
	keyImpIsActive       = "is_active"
	keyImpCompleteStatus = "complete_status"
	keyImpImproveChanges = "improve_changes"
)

func NewImproveDao() ImproveDaoAssumer {
	return &improveDao{}
}

type improveDao struct {
}

type ImproveDaoAssumer interface {
	InsertImprove(input dto.Improve) (*string, rest_err.APIError)
	EditImprove(input dto.ImproveEdit) (*dto.Improve, rest_err.APIError)
	ChangeImprove(filterA dto.FilterIDBranch, data dto.ImproveChange) (*dto.Improve, rest_err.APIError)
	DeleteImprove(input dto.FilterIDBranchCreateGte) (*dto.Improve, rest_err.APIError)
	ActivateImprove(improveID primitive.ObjectID, user mjwt.CustomClaim, isEnable bool) (*dto.Improve, rest_err.APIError)

	GetImproveByID(improveID primitive.ObjectID, branchIfSpecific string) (*dto.Improve, rest_err.APIError)
	FindImprove(filterA dto.FilterBranchCompleteTimeRangeLimit) (dto.ImproveResponseMinList, rest_err.APIError)
}

func (s *improveDao) InsertImprove(input dto.Improve) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyImpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Title = strings.ToUpper(input.Title)
	input.Branch = strings.ToUpper(input.Branch)
	input.ImproveChanges = []dto.ImproveChange{}

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan improve ke database", err)
		logger.Error("Gagal menyimpan improve ke database, (InsertImprove)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (s *improveDao) EditImprove(input dto.ImproveEdit) (*dto.Improve, rest_err.APIError) {
	coll := db.DB.Collection(keyImpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Title = strings.ToUpper(input.Title)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyImpID:        input.FilterID,
		keyImpBranch:    input.FilterBranch,
		keyImpUpdatedAt: input.FilterTimestamp,
	}

	update := bson.D{
		{"$set", bson.M{ //nolint:govet
			keyImpTitle:          input.Title,
			keyImpUpdatedAt:      input.UpdatedAt,
			keyImpUpdatedBy:      input.UpdatedBy,
			keyImpUpdatedByID:    input.UpdatedByID,
			keyImpDescription:    input.Description,
			keyImpGoal:           input.Goal,
			keyImpCompleteStatus: input.CompleteStatus,
		}},
	}

	var improve dto.Improve
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&improve); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Improve tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan improve dari database (EditImprove)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan improve dari database", err)
		return nil, apiErr
	}

	return &improve, nil
}

func (s *improveDao) ChangeImprove(filterA dto.FilterIDBranch, data dto.ImproveChange) (*dto.Improve, rest_err.APIError) {
	coll := db.DB.Collection(keyImpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyImpID:       filterA.FilterID,
		keyImpBranch:   strings.ToUpper(filterA.FilterBranch),
		keyImpIsActive: true,
	}

	update := bson.D{
		{"$set", bson.M{keyImpUpdatedAt: time.Now().Unix()}},  //nolint:govet
		{"$inc", bson.M{keyImpGoalsAchieved: data.Increment}}, //nolint:govet
		{"$push", bson.M{keyImpImproveChanges: data}},         //nolint:govet
	}

	var improve dto.Improve
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&improve); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Improve tidak diupdate : validasi id branch active")
		}

		logger.Error("Merubah nilai improve gagal, (ChangeImprove)", err)
		apiErr := rest_err.NewInternalServerError("Merubah jumlah improve gagal", err)
		return nil, apiErr
	}

	return &improve, nil
}

func (s *improveDao) DeleteImprove(input dto.FilterIDBranchCreateGte) (*dto.Improve, rest_err.APIError) {
	coll := db.DB.Collection(keyImpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyImpID:        input.FilterID,
		keyImpBranch:    input.FilterBranch,
		keyImpCreatedAt: bson.M{"$gte": input.FilterCreateGTE},
	}

	var improve dto.Improve
	err := coll.FindOneAndDelete(ctx, filter).Decode(&improve)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Improve tidak diupdate : validasi id branch time_reach")
		}

		logger.Error("Gagal menghapus improve dari database (DeleteImprove)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan improve dari database", err)
		return nil, apiErr
	}

	return &improve, nil
}

func (s *improveDao) ActivateImprove(improveID primitive.ObjectID, user mjwt.CustomClaim, isEnable bool) (*dto.Improve, rest_err.APIError) {
	coll := db.DB.Collection(keyImpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyImpID:     improveID,
		keyImpBranch: user.Branch,
	}

	update := bson.M{
		"$set": bson.M{
			keyImpIsActive:  isEnable,
			keyImpUpdatedAt: time.Now().Unix(),
		},
	}

	var improve dto.Improve
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&improve); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Improve tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mengupdate improve di database (ActivateImprove)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mengupdate improve di database", err)
		return nil, apiErr
	}

	return &improve, nil
}

func (s *improveDao) GetImproveByID(improveID primitive.ObjectID, branchIfSpecific string) (*dto.Improve, rest_err.APIError) {
	coll := db.DB.Collection(keyImpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyImpID: improveID}
	if branchIfSpecific != "" {
		filter[keyImpBranch] = strings.ToUpper(branchIfSpecific)
	}

	var improve dto.Improve
	if err := coll.FindOne(ctx, filter).Decode(&improve); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Improve dengan ID %s tidak ditemukan", improveID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan improve dari database (GetImproveByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan improve dari database", err)
		return nil, apiErr
	}

	return &improve, nil
}

func (s *improveDao) FindImprove(filterA dto.FilterBranchCompleteTimeRangeLimit) (dto.ImproveResponseMinList, rest_err.APIError) {
	coll := db.DB.Collection(keyImpCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyImpBranch] = strings.ToUpper(filterA.FilterBranch)
	}
	if filterA.FilterCompleteStatus != 0 {
		filter[keyImpCompleteStatus] = filterA.FilterCompleteStatus
	}
	if filterA.Limit == 0 {
		filterA.Limit = 300
	}

	// option range
	if filterA.FilterStart != 0 {
		filter[keyImpCreatedAt] = bson.M{"$gte": filterA.FilterStart}
	}
	if filterA.FilterEnd != 0 {
		filter[keyImpCreatedAt] = bson.M{"$lte": filterA.FilterEnd}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyImpID, Value: -1}})
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar improve dari database (FindImprove)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.ImproveResponseMinList{}, apiErr
	}

	improveList := dto.ImproveResponseMinList{}
	if err = cursor.All(ctx, &improveList); err != nil {
		logger.Error("Gagal decode improveList cursor ke objek slice (FindImprove)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return dto.ImproveResponseMinList{}, apiErr
	}

	return improveList, nil
}
