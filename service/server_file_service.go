package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/serverfiledao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func NewServerFileService(assumer serverfiledao.ServerFileAssumer) ServerFileServiceAssumer {
	return &serverFileService{
		daoS: assumer,
	}
}

type serverFileService struct {
	daoS serverfiledao.ServerFileAssumer
}

type ServerFileServiceAssumer interface {
	Insert(user mjwt.CustomClaim, input dto.ServerFileReq) (*string, rest_err.APIError)
	Delete(user mjwt.CustomClaim, inputID string) rest_err.APIError
	UploadImage(user mjwt.CustomClaim, inputID string, imagePath string) (*dto.ServerFile, rest_err.APIError)
	GetByID(stockID string, branchIfSpecific string) (*dto.ServerFile, rest_err.APIError)
	Find(branch string, start int64, end int64) ([]dto.ServerFile, rest_err.APIError)
}

func (s *serverFileService) Insert(user mjwt.CustomClaim, input dto.ServerFileReq) (*string, rest_err.APIError) {
	// Default value
	timeNow := time.Now().Unix()
	generatedID := primitive.NewObjectID()

	// DB
	insertedID, err := s.daoS.Insert(dto.ServerFile{
		ID:        generatedID,
		UpdatedAt: timeNow,
		UpdatedBy: user.Name,
		Branch:    user.Branch,
		Title:     strings.ToUpper(input.Title),
		Note:      input.Note,
		Diff:      input.Diff,
		Image:     input.Image,
	})
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (s *serverFileService) Delete(user mjwt.CustomClaim, inputID string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(inputID)
	if errT != nil {
		return rest_err.NewBadRequestError(errT.Error())
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := s.daoS.Delete(dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *serverFileService) UploadImage(user mjwt.CustomClaim, stockID string, imagePath string) (*dto.ServerFile, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	item, err := s.daoS.UploadImage(oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *serverFileService) GetByID(stockID string, branchIfSpecific string) (*dto.ServerFile, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	item, err := s.daoS.GetByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *serverFileService) Find(branch string, start int64, end int64) ([]dto.ServerFile, rest_err.APIError) {
	serverList, err := s.daoS.Find(branch, start, end)
	if err != nil {
		return nil, err
	}
	return serverList, nil
}
