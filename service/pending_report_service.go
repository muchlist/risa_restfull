package service

import (
	"context"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/constants/enum"
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
	fcmClient fcm.ClientAssumer,
) PRServiceAssumer {
	return &prService{
		daoP: prDao,
		daoG: genDao,
		daoU: userDao,
		fcm:  fcmClient,
	}
}

type prService struct {
	daoP pendingreportdao.PRAssumer
	daoG genunitdao.GenUnitLoader
	daoU userdao.UserLoader
	fcm  fcm.ClientAssumer
}

// hapus participant dan approver
// geser2 level complete status
// bikin pdf

type PRServiceAssumer interface {
	InsertPR(ctx context.Context, user mjwt.CustomClaim, input dto.PendingReportRequest) (*string, rest_err.APIError)
	AddParticipant(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError)
	AddApprover(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError)
	RemoveParticipant(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError)
	RemoveApprover(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError)
	SendToSigningMode(ctx context.Context, user mjwt.CustomClaim, id string) (*dto.PendingReportModel, rest_err.APIError)
	SendToDraftMode(ctx context.Context, user mjwt.CustomClaim, id string) (*dto.PendingReportModel, rest_err.APIError)
	SignDocument(ctx context.Context, user mjwt.CustomClaim, id string) (*dto.PendingReportModel, rest_err.APIError)
	EditPR(ctx context.Context, user mjwt.CustomClaim, id string, input dto.PendingReportEditRequest) (*dto.PendingReportModel, rest_err.APIError)
	DeleteImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.PendingReportModel, rest_err.APIError)
	PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.PendingReportModel, rest_err.APIError)
	GetPRByID(ctx context.Context, id string, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError)
	FindDocs(ctx context.Context, user mjwt.CustomClaim, filter dto.FilterFindPendingReport) ([]dto.PendingReportMin, rest_err.APIError)
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

func (ps *prService) RemoveParticipant(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.RemoveParticipant(ctx, pendingreportdao.ParticipantParams{
		ID: oid,
		Participant: dto.Participant{
			ID: userID,
		},
		FilterBranch: user.Branch,
		UpdatedAt:    time.Now().Unix(),
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
	})
}

func (ps *prService) RemoveApprover(ctx context.Context, user mjwt.CustomClaim, id string, userID string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.RemoveApprover(ctx, pendingreportdao.ParticipantParams{
		ID: oid,
		Participant: dto.Participant{
			ID: userID,
		},
		FilterBranch: user.Branch,
		UpdatedAt:    time.Now().Unix(),
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
	})
}

func (ps *prService) SignDocument(ctx context.Context, user mjwt.CustomClaim, id string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// get document dan cek apakah levelnya memenuhi syarat
	// dan user ada didalam participant atau approver
	doc, restErr := ps.daoP.GetPRByID(ctx, oid, "")
	if restErr != nil {
		return nil, restErr
	}

	if doc.CompleteStatus != enum.NeedSign {
		if doc.CompleteStatus == enum.CompletedSign {
			return nil, rest_err.NewBadRequestError("Dokumen yang sudah selesai tidak dapat ditandatangani")
		}
		return nil, rest_err.NewBadRequestError("Dokumen masih dalam status draft")
	}

	var userExist bool

	for index, val := range doc.Participants {
		if val.UserID == user.Identity {
			userExist = true
			doc.Participants[index].Sign = "SIGNED"
			doc.Participants[index].SignAt = time.Now().Unix()
		}
	}

	for index, val := range doc.Approvers {
		if val.UserID == user.Identity {
			userExist = true
			doc.Approvers[index].Sign = "SIGNED"
			doc.Approvers[index].SignAt = time.Now().Unix()
		}
	}

	if !userExist {
		return nil, rest_err.NewBadRequestError("User tidak termasuk kedalam dokumen")
	}

	// update pending report dengan data approvers dan participant terbaru
	doc, restErr = ps.daoP.EditParticipantApprover(ctx, pendingreportdao.EditParticipantParams{
		ID:           oid,
		FilterBranch: user.Branch,
		Participant:  doc.Participants,
		Approver:     doc.Approvers,
		UpdatedAt:    time.Now().Unix(),
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
	})
	if restErr != nil {
		return nil, restErr
	}

	// cek apakah participant sudah ttd semua
	// jika iya kirim notif ke approver
	completeParticipantSign := true
	completeApproverSign := true
	for _, val := range doc.Participants {
		if val.Sign == "" {
			completeParticipantSign = false
		}
	}

	if completeParticipantSign {
		// create map string approver id
		approverMaps := make(map[string]struct{}, 0)
		for _, apr := range doc.Approvers {
			if apr.Sign == "" {
				approverMaps[apr.UserID] = struct{}{}
				completeApproverSign = false
			}
		}

		go func(aprMap map[string]struct{}) {
			users, err := ps.daoU.FindUser(ctx, user.Branch)
			if err != nil {
				logger.Error("mendapatkan user gagal saat menambahkan fcm (SignDocument)", err)
			}
			var tokens []string
			for _, u := range users {
				// skip jika user == user yang mengirim
				if u.ID == user.Identity {
					continue
				}
				_, exist := aprMap[u.ID]
				if exist {
					tokens = append(tokens, u.FcmToken)
				}
			}
			// firebase
			ps.fcm.SendMessage(fcm.Payload{
				Title:          fmt.Sprint("Dokumen memerlukan persetujuan"),
				Message:        fmt.Sprintf("Dokumen dengan judul %s memerlukan persetujuan", doc.Title),
				ReceiverTokens: tokens,
			})
		}(approverMaps)
	}

	if completeApproverSign {
		doc, restErr = ps.daoP.ChangeCompleteStatus(ctx, oid, enum.CompletedSign, enum.NeedSign, user.Branch)
		if restErr != nil {
			logger.Error(fmt.Sprintf("gagal melakukan complete status document berita acara dengan oid : %s", id), restErr)
		}
	}

	return doc, restErr
}

func (ps *prService) SendToSigningMode(ctx context.Context, user mjwt.CustomClaim, id string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	doc, restErr := ps.daoP.ChangeCompleteStatus(ctx, oid, enum.NeedSign, enum.Draft, user.Branch)
	if restErr != nil {
		return nil, restErr
	}

	// create map string approver id
	participantMaps := make(map[string]struct{}, 0)
	for _, apr := range doc.Participants {
		participantMaps[apr.UserID] = struct{}{}
	}

	go func(partyMap map[string]struct{}) {
		users, err := ps.daoU.FindUser(ctx, user.Branch)
		if err != nil {
			logger.Error("mendapatkan user gagal saat menambahkan fcm (SendToSigningMode)", err)
		}
		var tokens []string
		for _, u := range users {
			// skip jika user == user yang mengirim
			if u.ID == user.Identity {
				continue
			}
			_, exist := partyMap[u.ID]
			if exist {
				tokens = append(tokens, u.FcmToken)
			}
		}
		// firebase
		ps.fcm.SendMessage(fcm.Payload{
			Title:          fmt.Sprint("Dokumen memerlukan tanda tangan"),
			Message:        fmt.Sprintf("Dokumen dengan judul %s memerlukan tanda tangan", doc.Title),
			ReceiverTokens: tokens,
		})
	}(participantMaps)

	// kirim notif ke participant

	return doc, restErr
}

func (ps *prService) SendToDraftMode(ctx context.Context, user mjwt.CustomClaim, id string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	doc, restErr := ps.daoP.ChangeCompleteStatus(ctx, oid, enum.Draft, enum.NeedSign, user.Branch)
	if restErr != nil {
		return nil, restErr
	}

	return doc, restErr
}

// PutImage memasukkan lokasi file (path) ke dalam database violation dengan mengecek kesesuaian branch
func (ps *prService) PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.UploadImage(ctx, oid, imagePath, user.Branch)
}

// DeleteImage menghapus lokasi file (path) ke dalam database violation dengan mengecek kesesuaian branch
func (ps *prService) DeleteImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.DeleteImage(ctx, oid, imagePath, user.Branch)
}

func (ps *prService) GetPRByID(ctx context.Context, id string, branchIfSpecific string) (*dto.PendingReportModel, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	return ps.daoP.GetPRByID(ctx, oid, branchIfSpecific)
}

func (ps *prService) FindDocs(ctx context.Context, user mjwt.CustomClaim, filter dto.FilterFindPendingReport) ([]dto.PendingReportMin, rest_err.APIError) {
	if filter.FilterBranch == "" {
		filter.FilterBranch = user.Branch
	}
	return ps.daoP.FindDoc(ctx, filter)
}
