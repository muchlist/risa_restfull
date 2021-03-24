package service

import (
	"errors"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dao/cctvdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"sync"
	"time"
)

func NewCctvService(cctvDao cctvdao.CctvDaoAssumer,
	histDao historydao.HistoryDaoAssumer,
	genDao genunitdao.GenUnitDaoAssumer) CctvServiceAssumer {
	return &cctvService{
		daoC: cctvDao,
		daoH: histDao,
		daoG: genDao,
	}
}

type cctvService struct {
	daoC cctvdao.CctvDaoAssumer
	daoH historydao.HistoryDaoAssumer
	daoG genunitdao.GenUnitDaoAssumer
}
type CctvServiceAssumer interface {
	InsertCctv(user mjwt.CustomClaim, input dto.CctvRequest) (*string, rest_err.APIError)
	EditCctv(user mjwt.CustomClaim, cctvID string, input dto.CctvEditRequest) (*dto.Cctv, rest_err.APIError)
	DeleteCctv(user mjwt.CustomClaim, id string) rest_err.APIError
	DisableCctv(cctvID string, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError)
	PutImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Cctv, rest_err.APIError)

	GetCctvByID(cctvID string, branchIfSpecific string) (*dto.Cctv, rest_err.APIError)
	FindCctv(filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError)
}

func (c *cctvService) InsertCctv(user mjwt.CustomClaim, input dto.CctvRequest) (*string, rest_err.APIError) {
	// FilterID digunakan juga oleh gen_unit_dao sehingga dibuat disini, bukan di database
	idGenerated := primitive.NewObjectID()

	// Default value
	// jika ip address tidak kosong cek apakah ip address valid, jika kosong maka isikan nilai default
	if input.IP != "" {
		if net.ParseIP(input.IP) == nil {
			return nil, rest_err.NewBadRequestError("FilterIP Address tidak valid")
		}
	} else {
		input.IP = "0.0.0.0"
	}

	// kembalian dari golang channel
	type result struct {
		id  *string
		err rest_err.APIError
	}

	resultChan := make(chan result, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	insertToCctv := func() {
		defer wg.Done()
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

		// DB
		insertedID, err := c.daoC.InsertCctv(data)

		resultChan <- result{
			id:  insertedID,
			err: err,
		}
	}

	insertToGenUnit := func() {
		defer wg.Done()
		// Menambahkan juga General Unit dengan ID yang sama
		// DB
		insertedID, err := c.daoG.InsertUnit(
			dto.GenUnit{
				ID:       idGenerated.Hex(),
				Category: category.Cctv,
				Name:     input.Name,
				IP:       input.IP,
				Branch:   user.Branch,
			})
		resultChan <- result{
			id:  insertedID,
			err: err,
		}
	}

	go insertToCctv()
	go insertToGenUnit()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var resultID *string
	var errString string
	for v := range resultChan {
		if v.err != nil {
			errString = errString + v.err.Message() + ". "
		}
		if resultID == nil {
			resultID = v.id
		}
	}

	if errString != "" {
		return nil, rest_err.NewBadRequestError(errString)
	}

	return resultID, nil
}

func (c *cctvService) EditCctv(user mjwt.CustomClaim, cctvID string, input dto.CctvEditRequest) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(cctvID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// cek apakah ip address valid, jika kosong masukkan ip default 0.0.0.0
	if input.IP != "" {
		if net.ParseIP(input.IP) == nil {
			return nil, rest_err.NewBadRequestError("FilterIP Address tidak valid")
		}
	} else {
		input.IP = "0.0.0.0"
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.CctvEdit{ // Prefer filter dan datanya dipisah dengan dua struct
		ID:              oid,
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		Name:            input.Name,
		IP:              input.IP,
		InventoryNumber: input.InventoryNumber,
		Location:        input.Location,
		LocationLat:     input.LocationLat,
		LocationLon:     input.LocationLon,
		Date:            input.Date,
		Tag:             input.Tag,
		Brand:           input.Brand,
		Type:            input.Type,
		Note:            input.Note,
	}

	// DB
	cctvEdited, err := c.daoC.EditCctv(data)
	if err != nil {
		return nil, err
	}

	errChan := make(chan rest_err.APIError, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	editUnit := func() {
		defer wg.Done()
		// DB
		_, err = c.daoG.EditUnit(cctvID, dto.GenUnitEditRequest{
			Category: category.Cctv,
			Name:     cctvEdited.Name,
			IP:       cctvEdited.IP,
			Branch:   cctvEdited.Branch,
		})
		errChan <- err
	}

	insertHistory := func() {
		defer wg.Done()
		// DB
		_, err = c.daoH.InsertHistory(
			dto.History{
				ID:             primitive.NewObjectID(),
				CreatedAt:      timeNow,
				CreatedBy:      user.Name,
				CreatedByID:    user.Identity,
				UpdatedAt:      timeNow,
				UpdatedBy:      user.Name,
				UpdatedByID:    user.Identity,
				Category:       category.Cctv,
				Branch:         user.Branch,
				ParentID:       cctvID,
				ParentName:     cctvEdited.Name,
				Status:         "Edit",
				Problem:        "Detail Cctv dirubah",
				ProblemResolve: "",
				CompleteStatus: enum.HInfo,
				DateStart:      timeNow,
				DateEnd:        timeNow,
				Tag:            []string{},
				Image:          "",
			})
		errChan <- err
	}

	go editUnit()
	go insertHistory()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errMessage string
	for err := range errChan {
		if err != nil {
			errMessage += ". "
		}
	}

	if errMessage != "" {
		errorPartial := errors.New("partial error")
		restErr := rest_err.NewInternalServerError(fmt.Sprintf("Mengedit cctv berhasil namun : %s", errMessage), errorPartial)
		logger.Error(restErr.Message(), errorPartial)
		return nil, restErr
	}

	return cctvEdited, nil
}

func (c *cctvService) DeleteCctv(user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := c.daoC.DeleteCctv(dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}

	// Delete unit_gen
	// DB
	err = c.daoG.DeleteUnit(id)
	if err != nil {
		return err
	}

	return nil
}

// DisableCctv if value true , cctv will disabled
func (c *cctvService) DisableCctv(cctvID string, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(cctvID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable cctv
	cctv, err := c.daoC.DisableCctv(oid, user, value)
	if err != nil {
		return nil, err
	}

	// set disable enable gen_unit
	_, err = c.daoG.DisableUnit(oid.Hex(), value)
	if err != nil {
		return nil, err
	}

	return cctv, nil
}

// PutImage memasukkan lokasi file (path) ke dalam database cctv dengan mengecek kesesuaian branch
func (c *cctvService) PutImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	cctv, err := c.daoC.UploadImage(oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return cctv, nil
}

func (c *cctvService) GetCctvByID(cctvID string, branchIfSpecific string) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(cctvID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	cctv, err := c.daoC.GetCctvByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return cctv, nil
}

func (c *cctvService) FindCctv(filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, rest_err.APIError) {
	// cek apakah ip address valid, jika valid maka set filter.FilterName ke kosong supaya pencarian berdasarkan IP
	if filter.FilterIP != "" {
		if net.ParseIP(filter.FilterIP) == nil {
			return nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
		filter.FilterName = ""
	}

	cctvList, err := c.daoC.FindCctv(filter)
	if err != nil {
		return nil, err
	}
	return cctvList, nil
}
