package pendingreportdao

import (
	"context"
	"errors"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
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

	keyParticipantsID = "id" // id inner participant
)

func NewPR() PRAssumer {
	return &prDao{}
}

type PRAssumer interface {
	InsertPR(ctx context.Context, input dto.PendingReportModel) (*string, rest_err.APIError)
	EditPR(ctx context.Context, input dto.PendingReportEditModel) (*dto.PendingReportModel, rest_err.APIError)
	ChangeCompleteStatus(ctx context.Context, id primitive.ObjectID, completeStatus int, filterCompleteStatus int, filterBranch string) (*dto.PendingReportModel, rest_err.APIError)
	AddApprover(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError)
	AddParticipant(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError)
	EditParticipantApprover(ctx context.Context, input EditParticipantParams) (*dto.PendingReportModel, rest_err.APIError)
	RemoveApprover(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError)
	RemoveParticipant(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError)
	UploadImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReportModel, rest_err.APIError)
	DeleteImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReportModel, rest_err.APIError)
	GetPRByID(ctx context.Context, id primitive.ObjectID, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError)
	GetPRByNumber(ctx context.Context, number string, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError)
	FindDoc(ctx context.Context, inFilter dto.FilterFindPendingReport) ([]dto.PendingReportMin, rest_err.APIError)
}

type prDao struct{}

func (pd *prDao) InsertPR(ctx context.Context, input dto.PendingReportModel) (*string, rest_err.APIError) {
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

func (pd *prDao) EditPR(ctx context.Context, input dto.PendingReportEditModel) (*dto.PendingReportModel, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	input.NormalizeValue()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:             input.FilterID,
		keyBranch:         input.FilterBranch,
		keyUpdatedAt:      input.FilterTimestamp,
		keyCompleteStatus: enum.Draft,
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

	var res dto.PendingReportModel
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch timestamp complete_status")
		}

		logger.Error("Gagal mendapatkan doc dari database (EditPR)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) ChangeCompleteStatus(ctx context.Context, id primitive.ObjectID, completeStatus int, filterCompleteStatus int, filterBranch string) (*dto.PendingReportModel, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:             id,
		keyBranch:         filterBranch,
		keyCompleteStatus: filterCompleteStatus,
	}

	update := bson.M{
		"$set": bson.M{
			keyCompleteStatus: completeStatus,
			keyUpdatedAt:      time.Now().Unix(),
		},
	}

	var res dto.PendingReportModel
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch complete_status")
		}

		logger.Error("Gagal mendapatkan doc dari database (ChangeCompleteStatus)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) AddApprover(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError) {
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

	var res dto.PendingReportModel
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

func (pd *prDao) AddParticipant(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError) {
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

	var res dto.PendingReportModel
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal menambahkan approver ke database (AddParticipant)", err)
		apiErr := rest_err.NewInternalServerError("Gagal menambahkan approver ke database (AddParticipant)", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) RemoveApprover(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError) {
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
		"$pull": bson.M{
			keyApprovers: bson.M{
				keyParticipantsID: strings.ToUpper(input.Participant.ID),
			},
		},
	}

	var res dto.PendingReportModel
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mengahapus approver doc dari database (RemoveApprover)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mengahapus approver doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) RemoveParticipant(ctx context.Context, input ParticipantParams) (*dto.PendingReportModel, rest_err.APIError) {
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
		"$pull": bson.M{
			keyParticipants: bson.M{
				keyParticipantsID: strings.ToUpper(input.Participant.ID),
			},
		},
	}

	var res dto.PendingReportModel
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mengahapus partisipan doc dari database (RemoveParticipant)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mengahapus partisipan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) EditParticipantApprover(ctx context.Context, input EditParticipantParams) (*dto.PendingReportModel, rest_err.APIError) {
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
			keyParticipants: input.Participant,
			keyApprovers:    input.Approver,
			keyUpdatedByID:  input.UpdatedByID,
			keyUpdatedBy:    input.UpdatedBy,
			keyUpdatedAt:    input.UpdatedAt,
		},
	}

	var res dto.PendingReportModel
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch")
		}

		logger.Error("Gagal mendapatkan doc dari database (Sign)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) UploadImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReportModel, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:             id,
		keyBranch:         strings.ToUpper(filterBranch),
		keyCompleteStatus: enum.Draft,
	}

	update := bson.M{
		"$push": bson.M{
			keyImages: imagePath,
		},
	}

	var res dto.PendingReportModel
	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError("Doc tidak diupdate : validasi id branch complete_status")
		}

		logger.Error("Gagal mendapatkan doc dari database (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan doc dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) DeleteImage(ctx context.Context, id primitive.ObjectID, imagePath string, filterBranch string) (*dto.PendingReportModel, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyID:             id,
		keyBranch:         strings.ToUpper(filterBranch),
		keyCompleteStatus: enum.Draft,
	}

	var pr dto.PendingReportModel
	if err := coll.FindOne(ctxt, filter).Decode(&pr); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("delete image gagal, doc dengan id %s tidak ditemukan", id.Hex()))
		}

		logger.Error("Delete image doc dari db gagal, (DeleteImage)", err)
		apiErr := rest_err.NewInternalServerError("Delete image doc dari db gagal", err)
		return nil, apiErr
	}

	// mendelete dari data yang sudah ditemukan
	var finalImages []string
	for _, image := range pr.Images {
		if image != imagePath {
			finalImages = append(finalImages, image)
		}
	}

	// jika final images 0 maka isikan slice kosong, untuk menghindari nill di db
	if len(finalImages) == 0 {
		finalImages = []string{}
	}

	update := bson.M{
		"$set": bson.M{
			keyImages: finalImages,
		},
	}

	if err := coll.FindOneAndUpdate(ctxt, filter, update, opts).Decode(&pr); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, rest_err.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, doc dengan id %s tidak ditemukan", id.Hex()))
		}

		logger.Error("Memasukkan path image doc ke db gagal, (UploadImage)", err)
		apiErr := rest_err.NewInternalServerError("Memasukkan path image doc ke db gagal", err)
		return nil, apiErr
	}

	return &pr, nil
}

func (pd *prDao) GetPRByID(ctx context.Context, id primitive.ObjectID, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyID: id}
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var res dto.PendingReportModel
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

func (pd *prDao) GetPRByNumber(ctx context.Context, number string, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyNumber: strings.ToUpper(number)}
	if branchIfSpecific != "" {
		filter[keyBranch] = strings.ToUpper(branchIfSpecific)
	}

	var res dto.PendingReportModel
	if err := coll.FindOne(ctxt, filter).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := rest_err.NewNotFoundError(fmt.Sprintf("Data dengan nomer %s tidak ditemukan", number))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan data dari database (GetPRByNumber)", err)
		apiErr := rest_err.NewInternalServerError("Gagal mendapatkan data dari database", err)
		return nil, apiErr
	}

	return &res, nil
}

func (pd *prDao) FindDoc(ctx context.Context, inFilter dto.FilterFindPendingReport) ([]dto.PendingReportMin, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	inFilter.FilterBranch = strings.ToUpper(inFilter.FilterBranch)

	// filter
	filter := bson.M{}

	if inFilter.FilterBranch != "" {
		filter[keyBranch] = inFilter.FilterBranch
	}

	// complete status
	if inFilter.CompleteStatus != "" {
		// cek jika multi status (pisah dengan koma)
		if strings.Contains(inFilter.CompleteStatus, ",") {
			completeString := strings.Split(inFilter.CompleteStatus, ",")
			completeInts := make([]int, 0)
			for _, val := range completeString {
				completeInt, err := strconv.Atoi(val)
				if err != nil {
					continue
				}
				completeInts = append(completeInts, completeInt)
			}
			filter[keyCompleteStatus] = bson.M{"$in": completeInts}
		} else {
			completeStatus, err := strconv.Atoi(inFilter.CompleteStatus)
			if err == nil {
				// hanya jika error nil , pasang filter
				filter[keyCompleteStatus] = completeStatus
			}
		}
	}

	// filter title
	if inFilter.FilterTitle != "" {
		filter[keyTitle] = bson.M{
			"$regex": fmt.Sprintf(".*%s", inFilter.FilterTitle),
		}
	}

	if inFilter.Limit == 0 {
		inFilter.Limit = 100
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: keyID, Value: -1}})
	opts.SetLimit(inFilter.Limit)

	cursor, err := coll.Find(ctxt, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar document dari database (FindDoc)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.PendingReportMin{}, apiErr
	}

	docList := make([]dto.PendingReportMin, 0)
	if err = cursor.All(ctxt, &docList); err != nil {
		logger.Error("Gagal decode docList cursor ke objek slice (FindDoc)", err)
		apiErr := rest_err.NewInternalServerError("Database error", err)
		return []dto.PendingReportMin{}, apiErr
	}

	return docList, nil
}
