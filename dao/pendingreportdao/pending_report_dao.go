package pendingreportdao

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
	keyCollection  = "pendingReport"

	keyID             = "_id"
	keyCreatedAt      = "created_at"
	keyCreatedBy      = "created_by"
	keyCreatedByID    = "created_by_id"
	keyUpdatedAt      = "updated_at"
	keyUpdatedBy      = "updated_by"
	keyUpdatedByID    = "updated_by_id"
	keyBranch         = "branch"
	keyNumber         = "number"
	keyTitle          = "title"
	keyDescriptions   = "descriptions"
	keyDate           = "date"
	keyParticipants   = "participants"
	keyApprovers      = "approvers"
	keyEquipments     = "equipments"
	keyCompleteStatus = "complete_status"
	keyLocation       = "location"
	keyImages         = "images"
)

func NewPR() PRAssumer {
	return &prDao{}
}

type PRAssumer interface {
	InsertPR(ctx context.Context, input dto.PendingReport) (*string, rest_err.APIError)
	EditPR(ctx context.Context, input dto.PendingReportEdit) (*dto.PendingReport, rest_err.APIError)
	GetPRByID(ctx context.Context, id primitive.ObjectID, branchIfSpecific string) (*dto.PendingReport, rest_err.APIError)
	ChangeCompleteStatus(ctx context.Context, id primitive.ObjectID, completeStatus int, filterBranch string) (*dto.PendingReport, rest_err.APIError)
	AddApprover(ctx context.Context, input AddParticipantInput) (*dto.PendingReport, rest_err.APIError)
	AddParticipant(ctx context.Context, input AddParticipantInput) (*dto.PendingReport, rest_err.APIError)
	RemoveApprover(ctx context.Context, id primitive.ObjectID, participantID string, filterBranch string) (*dto.PendingReport, rest_err.APIError)
	RemoveParticipant(ctx context.Context, id primitive.ObjectID, participantID string, filterBranch string) (*dto.PendingReport, rest_err.APIError)
	UploadImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReport, rest_err.APIError)
	DeleteImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReport, rest_err.APIError)
}

type prDao struct{}

func (pd *prDao) EditPR(ctx context.Context, input dto.PendingReportEdit) (*dto.PendingReport, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	input.NormalizeValue()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:        input.FilterID,
		keyBranch:    input.FilterBranch,
		keyUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyUpdatedAt:    input.UpdatedAt,
			keyUpdatedBy:    input.UpdatedBy,
			keyUpdatedByID:  input.UpdatedByID,
			keyNumber:       input.Number,
			keyTitle:        input.Title,
			keyDescriptions: input.Descriptions,
			keyEquipments:   input.Equipments,
			keyDate:         input.Date,
			keyLocation:     input.Location,
		},
	}

	var res dto.PendingReport
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch timestamp")
		}

		logger.Error("Gagal mendapatkan doc dari database (EditPR)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) ChangeCompleteStatus(ctx context.Context, id primitive.ObjectID, completeStatus int, filterBranch string) (*dto.PendingReport, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:     id,
		keyBranch: filterBranch,
	}

	update := bson.M{
		"$set": bson.M{
			keyCompleteStatus: completeStatus,
		},
	}

	var res dto.PendingReport
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendapatkan doc dari database (EditPR)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) AddApprover(ctx context.Context, input AddParticipantInput) (*dto.PendingReport, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:     input.ID,
		keyBranch: input.FilterBranch,
	}

	update := bson.M{
		"$set": bson.M{
			keyUpdatedBy:   input.UpdatedBy,
			keyUpdatedByID: input.UpdatedByID,
			keyUpdatedAt:   input.UpdatedAt,
		},
		"$push": bson.M{
			keyApprovers: input.Participant,
		},
	}

	var res dto.PendingReport
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendapatkan doc dari database (EditPR)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) AddParticipant(ctx context.Context, input AddParticipantInput) (*dto.PendingReport, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:     input.ID,
		keyBranch: input.FilterBranch,
	}

	update := bson.M{
		"$set": bson.M{
			keyUpdatedBy:   input.UpdatedBy,
			keyUpdatedByID: input.UpdatedByID,
			keyUpdatedAt:   input.UpdatedAt,
		},
		"$push": bson.M{
			keyParticipants: input.Participant,
		},
	}

	var res dto.PendingReport
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendapatkan doc dari database (EditPR)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) RemoveApprover(ctx context.Context, id primitive.ObjectID, participantID string, filterBranch string) (*dto.PendingReport, rest_err.APIError) {
	panic("implement me")
}

func (pd *prDao) RemoveParticipant(ctx context.Context, id primitive.ObjectID, participantID string, filterBranch string) (*dto.PendingReport, rest_err.APIError) {
	panic("implement me")
}

func (pd *prDao) UploadImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReport, rest_err.APIError) {
	panic("implement me")
}

func (pd *prDao) DeleteImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReport, rest_err.APIError) {
	panic("implement me")
}

func (pd *prDao) InsertPR(ctx context.Context, input dto.PendingReport) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	input.NormalizeValue()

	result, err := coll.InsertOne(ctxt, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan doc ke database", err)
		logger.Error("Gagal menyimpan cctv ke database, (InsertPR)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (pd *prDao) GetPRByID(ctx context.Context, id primitive.ObjectID, branchIfSpecific string) (*dto.PendingReport, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: id}
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var res dto.PendingReport
	if err := coll.FindOne(ctxt, filter).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Data dengan ID %s tidak ditemukan", id.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan data dari database (GetPRByID)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan data dari database", err)
		return nil, apiErr
	}

	return &res, nil
}
