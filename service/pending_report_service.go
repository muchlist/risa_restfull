package service

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/pendingreportdao"
	"github.com/muchlist/risa_restfull/dao/userdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewPRService(
	prDao pendingreportdao.PRAssumer,
	genDao genunitdao.GenUnitLoader,
	userDao userdao.UserLoader,
) PRServiceAssumer {
	return &prService{
		daoP: prDao,
		daoG: genDao,
		daoU: userDao,
	}
}

type prService struct {
	daoP pendingreportdao.PRAssumer
	daoG genunitdao.GenUnitLoader
	daoU userdao.UserLoader
}

// TODO insertPR validation V
// TODO Insert participant itu gimana atuh, apakah dari user yang sudah ada atau input sendiri ?
// jika user ID ada maka ambil dari user, jika tidak ada maka inputan dari user. atau buat dua service yang berbeda ?
// TODO Description type ,
// TODO Complete status report ada berapa level
// TODO insert approver sama dengan insert participant
// hapus participant dan approver
// geser2 level complete status
// bikin pdf

type PRServiceAssumer interface {
	InsertPR(ctx context.Context, user mjwt.CustomClaim, input dto.PendingReportRequest) (*string, rest_err.APIError)
	AddParticipant(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError)
	AddApprover(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError)
	GetPRByID(ctx context.Context, id string, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError)
	EditPR(ctx context.Context, user mjwt.CustomClaim, id string, input dto.PendingReportEditRequest) (*dto.PendingReportModel, rest_err.APIError)
}

func (ps *prService) InsertPR(ctx context.Context, user mjwt.CustomClaim, input dto.PendingReportRequest) (*string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	if input.Branch == "" {
		input.Branch = user.Branch
	}

	if input.Date == 0 {
		input.Date = timeNow
	}

	res, err := ps.daoP.InsertPR(ctx, dto.PendingReportModel{
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

func (ps *prService) GetPRByID(ctx context.Context, id string, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.GetPRByID(ctx, oid, branchIfSpecific)
}

func (ps *prService) EditPR(ctx context.Context, user mjwt.CustomClaim, id string, input dto.PendingReportEditRequest) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.EditPR(ctx, dto.PendingReportEditModel{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		UpdatedAt:       time.Now().Unix(),
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		Number:          input.Number,
		Title:           input.Title,
		Descriptions:    input.Descriptions,
		Date:            input.Date,
		Equipments:      input.Equipments,
		Location:        input.Location,
	})
}

func (ps *prService) AddParticipant(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// cek ketersediaan user
	userResult, restErr := ps.daoU.GetUserByID(ctx, userID)
	if restErr != nil {
		return nil, rest_err.NewNotFoundError("user yang dimasukkan tidak tersedia")
	}

	// cek apakah user tersebut sudah ada didalam daftar participant
	pendingReport, restErr := ps.daoP.GetPRByID(ctx, oid, "")
	if restErr != nil {
		return nil, rest_err.NewNotFoundError("dokumen yang dimasukkan tidak tersedia")
	}
	if pendingReport.Participants != nil && len(pendingReport.Participants) != 0 {
		for _, val := range pendingReport.Participants {
			if val.ID == userID {
				return nil, rest_err.NewBadRequestError("Participant yang dimasukkan sudah ada pada dokumen eksisting")
			}
		}
	}
	if pendingReport.Approvers != nil && len(pendingReport.Approvers) != 0 {
		for _, val := range pendingReport.Approvers {
			if val.ID == userID {
				return nil, rest_err.NewBadRequestError("Participant yang dimasukkan sudah ada pada dokumen eksisting")
			}
		}
	}

	return ps.daoP.AddParticipant(ctx, pendingreportdao.ParticipantParams{
		ID: oid,
		Participant: dto.Participant{
			ID:       userResult.ID,
			Name:     userResult.Name,
			Position: userResult.Position,
			Division: userResult.Division,
			UserID:   userResult.ID,
			Sign:     "",
			SignAt:   0,
		},
		FilterBranch: user.Branch,
		UpdatedAt:    time.Now().Unix(),
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
	})
}

func (ps *prService) AddApprover(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// cek ketersediaan user
	userResult, restErr := ps.daoU.GetUserByID(ctx, userID)
	if restErr != nil {
		return nil, rest_err.NewNotFoundError("user yang dimasukkan tidak tersedia")
	}

	// cek apakah user tersebut sudah ada didalam daftar participant
	pendingReport, restErr := ps.daoP.GetPRByID(ctx, oid, "")
	if restErr != nil {
		return nil, rest_err.NewNotFoundError("dokumen yang dimasukkan tidak tersedia")
	}
	if pendingReport.Participants != nil && len(pendingReport.Participants) != 0 {
		for _, val := range pendingReport.Participants {
			if val.ID == userID {
				return nil, rest_err.NewBadRequestError("Approver yang dimasukkan sudah ada pada dokumen eksisting")
			}
		}
	}
	if pendingReport.Approvers != nil && len(pendingReport.Approvers) != 0 {
		for _, val := range pendingReport.Approvers {
			if val.ID == userID {
				return nil, rest_err.NewBadRequestError("Approver yang dimasukkan sudah ada pada dokumen eksisting")
			}
		}
	}

	return ps.daoP.AddApprover(ctx, pendingreportdao.ParticipantParams{
		ID: oid,
		Participant: dto.Participant{
			ID:       userResult.ID,
			Name:     userResult.Name,
			Position: userResult.Position,
			Division: userResult.Division,
			UserID:   userResult.ID,
			Sign:     "",
			SignAt:   0,
		},
		FilterBranch: user.Branch,
		UpdatedAt:    time.Now().Unix(),
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
	})
}
