package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/dao/computerdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"sync"
	"time"
)

func NewComputerService(computerDao computerdao.ComputerDaoAssumer,
	histDao historydao.HistorySaver,
	genDao genunitdao.GenUnitDaoAssumer) ComputerServiceAssumer {
	return &computerService{
		daoC: computerDao,
		daoH: histDao,
		daoG: genDao,
	}
}

type computerService struct {
	daoC computerdao.ComputerDaoAssumer
	daoH historydao.HistorySaver
	daoG genunitdao.GenUnitDaoAssumer
}
type ComputerServiceAssumer interface {
	InsertComputer(ctx context.Context, user mjwt.CustomClaim, input dto.ComputerRequest) (*string, rest_err.APIError)
	EditComputer(ctx context.Context, user mjwt.CustomClaim, computerID string, input dto.ComputerEditRequest) (*dto.Computer, rest_err.APIError)
	DeleteComputer(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	DisableComputer(ctx context.Context, computerID string, user mjwt.CustomClaim, value bool) (*dto.Computer, rest_err.APIError)
	PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Computer, rest_err.APIError)

	GetComputerByID(ctx context.Context, computerID string, branchIfSpecific string) (*dto.Computer, rest_err.APIError)
	FindComputer(ctx context.Context, filter dto.FilterComputer) (dto.ComputerResponseMinList, dto.GenUnitResponseList, rest_err.APIError)
}

func (c *computerService) InsertComputer(ctx context.Context, user mjwt.CustomClaim, input dto.ComputerRequest) (*string, rest_err.APIError) {
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

	insertToComputer := func() {
		defer wg.Done()
		// Filling data
		timeNow := time.Now().Unix()
		data := dto.Computer{
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

			Hostname:       input.Hostname,
			Division:       input.Division,
			SeatManagement: input.SeatManagement,
			OS:             input.OS,
			Processor:      input.Processor,
			Ram:            input.Ram,
			Hardisk:        input.Hardisk,

			Date:  input.Date,
			Tag:   input.Tag,
			Image: "", // image empty saat inisialisasi
			Brand: input.Brand,
			Type:  input.Type,
			Note:  input.Note,
		}

		// DB
		insertedID, err := c.daoC.InsertPc(ctx, data)

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
				Category: category.PC,
				Name:     input.Name,
				IP:       input.IP,
				Branch:   user.Branch,
			})
		resultChan <- result{
			id:  insertedID,
			err: err,
		}
	}

	go insertToComputer()
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

func (c *computerService) EditComputer(ctx context.Context, user mjwt.CustomClaim, computerID string, input dto.ComputerEditRequest) (*dto.Computer, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(computerID)
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
	data := dto.ComputerEdit{ // Prefer filter dan datanya dipisah dengan dua struct
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
		Hostname:        input.Hostname,
		Division:        input.Division,
		SeatManagement:  input.SeatManagement,
		OS:              input.OS,
		Processor:       input.Processor,
		Ram:             input.Ram,
		Hardisk:         input.Hardisk,
		Date:            input.Date,
		Tag:             input.Tag,
		Brand:           input.Brand,
		Type:            input.Type,
		Note:            input.Note,
	}

	// DB
	computerEdited, err := c.daoC.EditPc(ctx, data)
	if err != nil {
		return nil, err
	}

	errChan := make(chan rest_err.APIError, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	editUnit := func() {
		defer wg.Done()
		// DB
		_, err = c.daoG.EditUnit(ctx, computerID, dto.GenUnitEditRequest{
			Category: category.PC,
			Name:     computerEdited.Name,
			IP:       computerEdited.IP,
			Branch:   computerEdited.Branch,
		})
		errChan <- err
	}

	insertHistory := func() {
		defer wg.Done()
		isVendor := sfunc.InSlice(roles.RoleVendor, user.Roles)
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
				Category:       category.PC,
				Branch:         user.Branch,
				ParentID:       computerID,
				ParentName:     computerEdited.Name,
				Status:         "Edit",
				Problem:        "Detail PC diubah",
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
		restErr := rest_err.NewInternalServerError(fmt.Sprintf("Mengedit computer berhasil namun : %s", errMessage), errorPartial)
		logger.Error(restErr.Message(), errorPartial)
		return nil, restErr
	}

	return computerEdited, nil
}

func (c *computerService) DeleteComputer(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := c.daoC.DeletePc(ctx, dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}

	// Delete unit_gen
	// DB
	err = c.daoG.DeleteUnit(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// DisableComputer if value true , computer will disabled
func (c *computerService) DisableComputer(ctx context.Context, computerID string, user mjwt.CustomClaim, value bool) (*dto.Computer, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(computerID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable computer
	computer, err := c.daoC.DisablePc(ctx, oid, user, value)
	if err != nil {
		return nil, err
	}

	// set disable enable gen_unit
	_, err = c.daoG.DisableUnit(ctx, oid.Hex(), value)
	if err != nil {
		return nil, err
	}

	return computer, nil
}

// PutImage memasukkan lokasi file (path) ke dalam database computer dengan mengecek kesesuaian branch
func (c *computerService) PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Computer, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	computer, err := c.daoC.UploadImage(ctx, oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return computer, nil
}

func (c *computerService) GetComputerByID(ctx context.Context, computerID string, branchIfSpecific string) (*dto.Computer, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(computerID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// kembalian dari golang channel
	type resultComputer struct {
		data *dto.Computer
		err  rest_err.APIError
	}

	type resultGeneral struct {
		data *dto.GenUnitResponse
		err  rest_err.APIError
	}

	resultComputerChan := make(chan resultComputer)
	resultGeneralChan := make(chan resultGeneral)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// DB
		computer, err := c.daoC.GetPcByID(ctx, oid, branchIfSpecific)
		resultComputerChan <- resultComputer{
			data: computer,
			err:  err,
		}
	}()

	go func() {
		defer wg.Done()
		// DB
		computer, err := c.daoG.GetUnitByID(ctx, computerID, branchIfSpecific)
		resultGeneralChan <- resultGeneral{
			data: computer,
			err:  err,
		}
	}()

	go func() {
		wg.Wait()
		close(resultComputerChan)
		close(resultGeneralChan)
	}()

	computerDetail := <-resultComputerChan
	if computerDetail.err != nil {
		return nil, computerDetail.err
	}

	computerExtraData := <-resultGeneralChan
	if computerExtraData.err != nil {
		return nil, computerExtraData.err
	}

	computerData := computerDetail.data
	computerData.Extra.Cases = computerExtraData.data.Cases
	computerData.Extra.LastPing = computerExtraData.data.LastPing
	computerData.Extra.CasesSize = computerExtraData.data.CasesSize
	computerData.Extra.PingsState = computerExtraData.data.PingsState

	return computerData, nil
}

func (c *computerService) FindComputer(ctx context.Context, filter dto.FilterComputer) (dto.ComputerResponseMinList, dto.GenUnitResponseList, rest_err.APIError) {
	// cek apakah ip address valid, jika valid maka set filter.FilterName ke kosong supaya pencarian berdasarkan IP
	if filter.FilterIP != "" {
		if net.ParseIP(filter.FilterIP) == nil {
			return nil, nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
		filter.FilterName = ""
	}

	// wrap golang channel
	type resultComputer struct {
		data dto.ComputerResponseMinList
		err  rest_err.APIError
	}

	// wrap golang channel
	type resultGeneral struct {
		data dto.GenUnitResponseList
		err  rest_err.APIError
	}

	computerListChan := make(chan resultComputer)
	generalListChan := make(chan resultGeneral)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		computerList, err := c.daoC.FindPc(ctx, filter)
		computerListChan <- resultComputer{
			data: computerList,
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
			Category: category.PC,
			Pings:    false,
			Name:     "",
		})
		generalListChan <- resultGeneral{
			data: generalList,
			err:  err,
		}
	}()

	computerList := <-computerListChan
	if computerList.err != nil {
		return nil, nil, computerList.err
	}

	generalList := <-generalListChan
	if generalList.err != nil {
		return nil, nil, generalList.err
	}

	filterGeneral(&generalList.data)
	return computerList.data, generalList.data, nil
}
