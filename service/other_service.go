package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/otherdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"strings"
	"sync"
	"time"
)

func NewOtherService(otherDao otherdao.OtherDaoAssumer,
	histDao historydao.HistorySaver,
	genDao genunitdao.GenUnitDaoAssumer) OtherServiceAssumer {
	return &otherService{
		daoO: otherDao,
		daoH: histDao,
		daoG: genDao,
	}
}

type otherService struct {
	daoO otherdao.OtherDaoAssumer
	daoH historydao.HistorySaver
	daoG genunitdao.GenUnitDaoAssumer
}
type OtherServiceAssumer interface {
	InsertOther(ctx context.Context, user mjwt.CustomClaim, input dto.OtherRequest) (*string, rest_err.APIError)
	EditOther(ctx context.Context, user mjwt.CustomClaim, otherID string, input dto.OtherEditRequest) (*dto.Other, rest_err.APIError)
	DeleteOther(ctx context.Context, user mjwt.CustomClaim, subCategory string, id string) rest_err.APIError
	DisableOther(ctx context.Context, otherID string, user mjwt.CustomClaim, subCategory string, value bool) (*dto.Other, rest_err.APIError)
	PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Other, rest_err.APIError)

	GetOtherByID(ctx context.Context, otherID string, branchIfSpecific string) (*dto.Other, rest_err.APIError)
	FindOther(ctx context.Context, filter dto.FilterOther) (dto.OtherResponseMinList, dto.GenUnitResponseList, rest_err.APIError)
}

func (c *otherService) InsertOther(ctx context.Context, user mjwt.CustomClaim, input dto.OtherRequest) (*string, rest_err.APIError) {
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

	subCategory := strings.ToUpper(input.SubCategory)

	// kembalian dari golang channel
	type result struct {
		id  *string
		err rest_err.APIError
	}

	resultChan := make(chan result, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	insertToOther := func() {
		defer wg.Done()
		// Filling data
		timeNow := time.Now().Unix()
		data := dto.Other{
			CreatedAt:       timeNow,
			CreatedBy:       user.Name,
			CreatedByID:     user.Identity,
			UpdatedAt:       timeNow,
			UpdatedBy:       user.Name,
			UpdatedByID:     user.Identity,
			Branch:          user.Branch,
			SubCategory:     subCategory,
			ID:              idGenerated,
			Name:            input.Name,
			IP:              input.IP,
			InventoryNumber: input.InventoryNumber,
			Location:        input.Location,
			LocationLat:     input.LocationLat,
			LocationLon:     input.LocationLon,

			Division: input.Division,
			Detail:   input.Detail,

			Date:  input.Date,
			Tag:   input.Tag,
			Image: "", // image empty saat inisialisasi
			Brand: input.Brand,
			Type:  input.Type,
			Note:  input.Note,
		}

		// DB
		insertedID, err := c.daoO.InsertOther(ctx, data)

		resultChan <- result{
			id:  insertedID,
			err: err,
		}
	}

	insertToGenUnit := func() {
		defer wg.Done()
		// Menambahkan juga General Unit dengan ID yang sama
		// DB
		insertedID, err := c.daoG.InsertUnit(ctx,
			dto.GenUnit{
				ID:       idGenerated.Hex(),
				Category: subCategory,
				Name:     input.Name,
				IP:       input.IP,
				Branch:   user.Branch,
			})
		resultChan <- result{
			id:  insertedID,
			err: err,
		}
	}

	go insertToOther()
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

func (c *otherService) EditOther(ctx context.Context, user mjwt.CustomClaim, otherID string, input dto.OtherEditRequest) (*dto.Other, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(otherID)
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

	subCategory := strings.ToUpper(input.FilterSubCategory)

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.OtherEdit{ // Prefer filter dan datanya dipisah dengan dua struct
		ID:                oid,
		FilterBranch:      user.Branch,
		FilterTimestamp:   input.FilterTimestamp,
		FilterSubCategory: subCategory,
		UpdatedAt:         timeNow,
		UpdatedBy:         user.Name,
		UpdatedByID:       user.Identity,
		Name:              input.Name,
		IP:                input.IP,
		InventoryNumber:   input.InventoryNumber,
		Location:          input.Location,
		LocationLat:       input.LocationLat,
		LocationLon:       input.LocationLon,
		Division:          input.Division,
		Detail:            input.Detail,
		Date:              input.Date,
		Tag:               input.Tag,
		Brand:             input.Brand,
		Type:              input.Type,
		Note:              input.Note,
	}

	// DB
	otherEdited, err := c.daoO.EditOther(ctx, data)
	if err != nil {
		return nil, err
	}

	errChan := make(chan rest_err.APIError, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	editUnit := func() {
		defer wg.Done()
		// DB
		_, err = c.daoG.EditUnit(ctx, otherID, dto.GenUnitEditRequest{
			Category: subCategory,
			Name:     otherEdited.Name,
			IP:       otherEdited.IP,
			Branch:   otherEdited.Branch,
		})
		errChan <- err
	}

	insertHistory := func() {
		isVendor := sfunc.InSlice(roles.RoleVendor, user.Roles)
		defer wg.Done()
		// DB
		_, err = c.daoH.InsertHistory(ctx,
			dto.History{
				ID:             primitive.NewObjectID(),
				CreatedAt:      timeNow,
				CreatedBy:      user.Name,
				CreatedByID:    user.Identity,
				UpdatedAt:      timeNow,
				UpdatedBy:      user.Name,
				UpdatedByID:    user.Identity,
				Category:       subCategory,
				Branch:         user.Branch,
				ParentID:       otherID,
				ParentName:     otherEdited.Name,
				Status:         "Edit",
				Problem:        fmt.Sprintf("Detail %s diubah", subCategory),
				ProblemResolve: "",
				CompleteStatus: enum.HDataInfo,
				DateStart:      timeNow,
				DateEnd:        timeNow,
				Tag:            []string{},
				Image:          "",
			}, isVendor)
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
		restErr := rest_err.NewInternalServerError(fmt.Sprintf("Mengedit %s berhasil namun : %s", subCategory, errMessage), errorPartial)
		logger.Error(restErr.Message(), errorPartial)
		return nil, restErr
	}

	return otherEdited, nil
}

func (c *otherService) DeleteOther(ctx context.Context, user mjwt.CustomClaim, subCategory string, otherID string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(otherID)
	if errT != nil {
		return rest_err.NewBadRequestError(errT.Error())
		//return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := c.daoO.DeleteOther(ctx, dto.FilterIDBranchCategoryCreateGte{
		FilterID:          oid,
		FilterBranch:      user.Branch,
		FilterSubCategory: subCategory,
		FilterCreateGTE:   timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}

	// Delete unit_gen
	// DB
	err = c.daoG.DeleteUnit(ctx, otherID)
	if err != nil {
		return err
	}

	return nil
}

// DisableOther if value true , other will disabled
func (c *otherService) DisableOther(ctx context.Context, otherID string, user mjwt.CustomClaim, subCategory string, value bool) (*dto.Other, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(otherID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable other
	other, err := c.daoO.DisableOther(ctx, oid, user, subCategory, value)
	if err != nil {
		return nil, err
	}

	// set disable enable gen_unit
	_, err = c.daoG.DisableUnit(ctx, oid.Hex(), value)
	if err != nil {
		return nil, err
	}

	return other, nil
}

// PutImage memasukkan lokasi file (path) ke dalam database other dengan mengecek kesesuaian branch
func (c *otherService) PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Other, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	other, err := c.daoO.UploadImage(ctx, oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return other, nil
}

func (c *otherService) GetOtherByID(ctx context.Context, otherID string, branchIfSpecific string) (*dto.Other, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(otherID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// kembalian dari golang channel
	type resultOther struct {
		data *dto.Other
		err  rest_err.APIError
	}

	type resultGeneral struct {
		data *dto.GenUnitResponse
		err  rest_err.APIError
	}

	resultOtherChan := make(chan resultOther)
	resultGeneralChan := make(chan resultGeneral)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// DB
		other, err := c.daoO.GetOtherByID(ctx, oid, branchIfSpecific)
		resultOtherChan <- resultOther{
			data: other,
			err:  err,
		}
	}()

	go func() {
		defer wg.Done()
		// DB
		other, err := c.daoG.GetUnitByID(ctx, otherID, branchIfSpecific)
		resultGeneralChan <- resultGeneral{
			data: other,
			err:  err,
		}
	}()

	go func() {
		wg.Wait()
		close(resultOtherChan)
		close(resultGeneralChan)
	}()

	otherDetail := <-resultOtherChan
	if otherDetail.err != nil {
		return nil, otherDetail.err
	}

	otherExtraData := <-resultGeneralChan
	if otherExtraData.err != nil {
		return nil, otherExtraData.err
	}

	otherData := otherDetail.data
	otherData.Extra.Cases = otherExtraData.data.Cases
	otherData.Extra.LastPing = otherExtraData.data.LastPing
	otherData.Extra.CasesSize = otherExtraData.data.CasesSize
	otherData.Extra.PingsState = otherExtraData.data.PingsState

	return otherData, nil
}

func (c *otherService) FindOther(ctx context.Context, filter dto.FilterOther) (dto.OtherResponseMinList, dto.GenUnitResponseList, rest_err.APIError) {
	// cek apakah ip address valid, jika valid maka set filter.FilterName ke kosong supaya pencarian berdasarkan IP
	if filter.FilterIP != "" {
		if net.ParseIP(filter.FilterIP) == nil {
			return nil, nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
		filter.FilterName = ""
	}

	// wrap golang channel
	type resultOther struct {
		data dto.OtherResponseMinList
		err  rest_err.APIError
	}

	// wrap golang channel
	type resultGeneral struct {
		data dto.GenUnitResponseList
		err  rest_err.APIError
	}

	otherListChan := make(chan resultOther)
	generalListChan := make(chan resultGeneral)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		otherList, err := c.daoO.FindOther(ctx, filter)
		otherListChan <- resultOther{
			data: otherList,
			err:  err,
		}
	}()

	go func() {
		defer wg.Done()

		// jika ada query pencarian, informasi ini tidak akan dimuat
		if filter.FilterIP != "" || filter.FilterName != "" {
			generalListChan <- resultGeneral{
				data: dto.GenUnitResponseList{},
				err:  nil,
			}
			return
		}

		generalList, err := c.daoG.FindUnit(ctx, dto.GenUnitFilter{
			Branch:   filter.FilterBranch,
			Category: filter.FilterSubCategory,
			Pings:    false,
			Name:     "",
		})
		generalListChan <- resultGeneral{
			data: generalList,
			err:  err,
		}
	}()

	otherList := <-otherListChan
	if otherList.err != nil {
		return nil, nil, otherList.err
	}

	generalList := <-generalListChan
	if generalList.err != nil {
		return nil, nil, generalList.err
	}

	filterGeneral(&generalList.data)
	return otherList.data, generalList.data, nil
}
