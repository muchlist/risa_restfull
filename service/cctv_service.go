package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/dao/cctv_dao"
	"github.com/muchlist/risa_restfull/dao/gen_unit_dao"
	"github.com/muchlist/risa_restfull/dao/history_dao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"time"
)

func NewCctvService(cctvDao cctv_dao.CctvDaoAssumer,
	histDao history_dao.HistoryDaoAssumer,
	genDao gen_unit_dao.GenUnitDaoAssumer) CctvServiceAssumer {
	return &cctvService{
		daoC: cctvDao,
		daoH: histDao,
		daoG: genDao,
	}
}

type cctvService struct {
	daoC cctv_dao.CctvDaoAssumer
	daoH history_dao.HistoryDaoAssumer
	daoG gen_unit_dao.GenUnitDaoAssumer
}
type CctvServiceAssumer interface {
	InsertCctv(user mjwt.CustomClaim, input dto.CctvRequest) (*string, rest_err.APIError)
	GetCctvByID(cctvID string) (*dto.Cctv, rest_err.APIError)
	FindCctv(filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError)
}

func (c *cctvService) InsertCctv(user mjwt.CustomClaim, input dto.CctvRequest) (*string, rest_err.APIError) {

	// ID digunakan juga oleh gen_unit_dao sehingga dibuat disini, bukan di database
	idGenerated := primitive.NewObjectID()

	// Default value
	// jika ip address tidak kosong cek apakah ip address valid, jika kosong maka isikan nilai default
	if input.IP != "" {
		if net.ParseIP(input.IP) == nil {
			return nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
	} else {
		input.IP = "0.0.0.0"
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.Cctv{
		CreatedAt:       timeNow,
		CreatedBy:       user.Name,
		CreatedByID:     user.Identity,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		Branch:          user.Branch,
		ID:              idGenerated,
		Name:            input.Name,
		IP:              input.IP,
		InventoryNumber: input.InventoryNumber,
		Location:        input.Location,
		LocationLat:     input.LocationLat,
		LocationLon:     input.LocationLon,
		Date:            input.Date,
		Tag:             input.Tag,
		Image:           "", // image empty saat inisialisasi
		Brand:           input.Brand,
		Type:            input.Type,
		Note:            input.Note,
	}

	//DB
	insertedID, err := c.daoC.InsertCctv(data)
	if err != nil {
		return nil, err
	}

	// Menambahkan juga General Unit dengan ID yang sama
	//DB
	_, err = c.daoG.InsertUnit(
		dto.GenUnitRequest{
			ID:       idGenerated.Hex(),
			Category: category.Cctv,
			Name:     input.Name,
			IP:       input.IP,
			Branch:   user.Branch,
		})
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (c *cctvService) GetCctvByID(cctvID string) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(cctvID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	history, err := c.daoC.GetCctvByID(oid)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (c *cctvService) FindCctv(filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError) {

	// cek apakah ip address valid, jika valid maka set filter.Name ke kosong supaya pencarian berdasarkan IP
	if filter.IP != "" {
		if net.ParseIP(filter.IP) == nil {
			return nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
		filter.Name = ""
	}

	historyList, err := c.daoC.FindCctv(filter)
	if err != nil {
		return nil, err
	}
	return historyList, nil
}
