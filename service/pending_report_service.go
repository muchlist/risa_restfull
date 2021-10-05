package service

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/pendingreportdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
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
}

func (ps *prService) InsertPR(ctx context.Context, user mjwt.CustomClaim, input dto.PendingReportRequest) (*string, rest_err.APIError) {
	return nil, nil
}
