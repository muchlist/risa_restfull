package service

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/dao/improvedao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewImproveService(improveDao improvedao.ImproveDaoAssumer) ImproveServiceAssumer {
	return &improveService{
		daoS: improveDao,
	}
}

type improveService struct {
	daoS improvedao.ImproveDaoAssumer
}
type ImproveServiceAssumer interface {
	InsertImprove(ctx context.Context, user mjwt.CustomClaim, input dto.ImproveRequest) (*string, rest_err.APIError)
	EditImprove(ctx context.Context, user mjwt.CustomClaim, improveID string, input dto.ImproveEditRequest) (*dto.Improve, rest_err.APIError)
	ActivateImprove(ctx context.Context, improveID string, user mjwt.CustomClaim, isEnable bool) (*dto.Improve, rest_err.APIError)
	ChangeImprove(ctx context.Context, user mjwt.CustomClaim, improveID string, data dto.ImproveChangeRequest) (*dto.Improve, rest_err.APIError)
	GetImproveByID(ctx context.Context, improveID string, branchIfSpecific string) (*dto.Improve, rest_err.APIError)
	DeleteImprove(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	FindImprove(ctx context.Context, filter dto.FilterBranchCompleteTimeRangeLimit) (dto.ImproveResponseMinList, rest_err.APIError)
}

func (s *improveService) InsertImprove(ctx context.Context, user mjwt.CustomClaim, input dto.ImproveRequest) (*string, rest_err.APIError) {
	// Filling data
	// Ketika membuat improve juga menambahkan increment field untuk pertama kali
	timeNow := time.Now().Unix()
	var improveAccepted bool
	// Jika role Approver ada di dalam user roles nya maka isActive otomatis true
	if sfunc.InSlice(roles.RoleApprove, user.Roles) {
		improveAccepted = true
	}

	oidGenerated := primitive.NewObjectID()
	data := dto.Improve{
		ID:             oidGenerated,
		CreatedAt:      timeNow,
		CreatedBy:      user.Name,
		CreatedByID:    user.Identity,
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Branch:         user.Branch,
		Title:          input.Title,
		Description:    input.Description,
		Goal:           input.Goal,
		CompleteStatus: input.CompleteStatus,
		IsActive:       improveAccepted,
	}

	// DB
	insertedID, err := s.daoS.InsertImprove(ctx, data)
	if err != nil {
		return nil, err
	}

	// IMPROVEMENT jika improveAccepted false memberikan notifikasi ke boss
	// agar bisa di approve

	return insertedID, nil
}

func (s *improveService) EditImprove(ctx context.Context, user mjwt.CustomClaim, improveID string, input dto.ImproveEditRequest) (*dto.Improve, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(improveID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.ImproveEdit{
		FilterIDBranchTimestamp: dto.FilterIDBranchTimestamp{
			FilterID:        oid,
			FilterBranch:    user.Branch,
			FilterTimestamp: input.FilterTimestamp,
		},
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Title:          input.Title,
		Description:    input.Description,
		Goal:           input.Goal,
		CompleteStatus: input.CompleteStatus,
	}

	// DB
	improveEdited, err := s.daoS.EditImprove(ctx, data)
	if err != nil {
		return nil, err
	}

	return improveEdited, nil
}

// DisableImprove if value true , improve will disabled
func (s *improveService) ActivateImprove(ctx context.Context, improveID string, user mjwt.CustomClaim, isEnable bool) (*dto.Improve, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(improveID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable improve
	improve, err := s.daoS.ActivateImprove(ctx, oid, user, isEnable)
	if err != nil {
		return nil, err
	}

	return improve, nil
}

func (s *improveService) ChangeImprove(ctx context.Context, user mjwt.CustomClaim, improveID string, data dto.ImproveChangeRequest) (*dto.Improve, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(improveID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Filling data=
	timeNow := time.Now().Unix()
	incDec := dto.ImproveChange{
		DummyID:   time.Now().UnixNano(),
		Author:    user.Name,
		Increment: data.Increment,
		Note:      data.Note,
		Time:      data.Time,
	}
	if data.Time == 0 {
		incDec.Time = timeNow
	}

	filter := dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	}

	// DB
	improveEdited, err := s.daoS.ChangeImprove(ctx, filter, incDec)
	if err != nil {
		return nil, err
	}

	return improveEdited, nil
}

func (s *improveService) DeleteImprove(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := s.daoS.DeleteImprove(ctx, dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *improveService) GetImproveByID(ctx context.Context, improveID string, branchIfSpecific string) (*dto.Improve, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(improveID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	improve, err := s.daoS.GetImproveByID(ctx, oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return improve, nil
}

func (s *improveService) FindImprove(ctx context.Context, filter dto.FilterBranchCompleteTimeRangeLimit) (dto.ImproveResponseMinList, rest_err.APIError) {
	improveList, err := s.daoS.FindImprove(ctx, filter)
	if err != nil {
		return nil, err
	}
	return improveList, nil
}
