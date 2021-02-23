package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/gen_unit_dao"
	"github.com/muchlist/risa_restfull/dto"
)

func NewGenUnitService(dao gen_unit_dao.GenUnitDaoAssumer) GenUnitServiceAssumer {
	return &genUnitService{
		dao: dao,
	}
}

type genUnitService struct {
	dao gen_unit_dao.GenUnitDaoAssumer
}
type GenUnitServiceAssumer interface {
	FindUnit(filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError)
	InsertCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
	DeleteCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError)
}

func (g *genUnitService) FindUnit(filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError) {
	unitList, err := g.dao.FindUnit(filter)
	if err != nil {
		return nil, err
	}
	return unitList, nil
}

func (g *genUnitService) InsertCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError) {
	unit, err := g.dao.InsertCase(payload)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func (g *genUnitService) DeleteCase(payload dto.GenUnitCaseRequest) (*dto.GenUnitResponse, rest_err.APIError) {
	unit, err := g.dao.DeleteCase(payload)
	if err != nil {
		return nil, err
	}
	return unit, nil
}
