package service

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/pendingreportdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewPRService(
	prDao pendingreportdao.PRAssumer,
	genDao genunitdao.GenUnitLoader,
) PRServiceAssumer {
	return &prService{
		daoP: prDao,
		daoG: genDao,
	}
}

type prService struct {
	daoP pendingreportdao.PRAssumer
	daoG genunitdao.GenUnitLoader
}

type PRServiceAssumer interface {
	InsertPR(ctx context.Context, user mjwt.CustomClaim, input dto.PendingReportRequest) (*string, rest_err.APIError)
	GetPRByID(ctx context.Context, id string, branchIfSpecific string) (*dto.PendingReport, rest_err.APIError)
}

func (ps *prService) InsertPR(ctx context.Context, user mjwt.CustomClaim, input dto.PendingReportRequest) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	if input.Branch == "" {
		input.Branch = user.Branch
	}

	if input.Date == 0 {
		input.Date = timeNow
	}

	res, err := ps.daoP.InsertPR(ctx, dto.PendingReport{
		ID:             primitive.NewObjectID(),
		CreatedAt:      timeNow,
		CreatedBy:      user.Name,
		CreatedByID:    user.Identity,
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Branch:         input.Branch,
		Number:         input.Number,
		Title:          input.Title,
		Descriptions:   input.Descriptions,
		Date:           input.Date,
		Participants:   nil,
		Approvers:      nil,
		Equipments:     input.Equipments,
		CompleteStatus: 0,
		Location:       input.Location,
		Images:         nil,
	})

	return res, err
}

func (ps *prService) GetPRByID(ctx context.Context, id string, branchIfSpecific string) (*dto.PendingReport, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.GetPRByID(ctx, oid, branchIfSpecific)
}
