package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/gen_unit_dao"
	"github.com/muchlist/risa_restfull/dto"
	"net"
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
	GetIPList(branchIfSpecific string, category string) ([]string, rest_err.APIError)
}

func (g *genUnitService) FindUnit(filter dto.GenUnitFilter) (dto.GenUnitResponseList, rest_err.APIError) {

	// cek apakah ip address valid, jika valid maka set filter.Name ke kosong supaya pencarian berdasarkan IP
	if filter.IP != "" {
		if net.ParseIP(filter.IP) == nil {
			return nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
		filter.Name = ""
	}

	// DB
	unitList, err := g.dao.FindUnit(filter)
	if err != nil {
		return nil, err
	}
	return unitList, nil
}

func (g *genUnitService) GetIPList(branchIfSpecific string, category string) ([]string, rest_err.APIError) {
	// DB
	ipAddressList, err := g.dao.GetIPList(branchIfSpecific, category)
	if err != nil {
		return nil, err
	}
	return ipAddressList, nil
}
