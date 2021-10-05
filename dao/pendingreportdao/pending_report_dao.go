package pendingreportdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	connectTimeout = 3
	keyCollection  = "pendingReport"

	keyID             = "id"
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
	InsertCctv(ctx context.Context, input dto.PendingReport) (*string, rest_err.APIError)
}

type prDao struct{}

func (pd *prDao) InsertCctv(ctx context.Context, input dto.PendingReport) (*string, rest_err.APIError) {
	coll := db.DB.Collection(keyCollection)
	ctxt, cancel := context.WithTimeout(ctx, connectTimeout*time.Second)
	defer cancel()

	input.NormalizeValue()

	result, err := coll.InsertOne(ctxt, input)
	if err != nil {
		apiErr := rest_err.NewInternalServerError("Gagal menyimpan cctv ke database", err)
		logger.Error("Gagal menyimpan cctv ke database, (InsertCctv)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}
