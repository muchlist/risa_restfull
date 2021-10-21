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
	"github.com/muchlist/risa_restfull/dao/cctvdao"
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
	InsertCctv(ctx context.Context, user mjwt.CustomClaim, input dto.CctvRequest) (*string, rest_err.APIError)
	EditCctv(ctx context.Context, user mjwt.CustomClaim, cctvID string, input dto.CctvEditRequest) (*dto.Cctv, rest_err.APIError)
	DeleteCctv(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	DisableCctv(ctx context.Context, cctvID string, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError)
	PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Cctv, rest_err.APIError)

	GetCctvByID(ctx context.Context, cctvID string, branchIfSpecific string) (*dto.Cctv, rest_err.APIError)
	FindCctv(ctx context.Context, filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, dto.GenUnitResponseList, rest_err.APIError)
	MergeCctv(ctx context.Context, cctvID1, cctvID2 string) (*string, rest_err.APIError)
}

func (c *cctvService) InsertCctv(ctx context.Context, user mjwt.CustomClaim, input dto.CctvRequest) (*string, rest_err.APIError) {
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
			DisVendor:       input.DisVendor,
		}

		// DB
		insertedID, err := c.daoC.InsertCctv(ctx, data)

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

func (c *cctvService) EditCctv(ctx context.Context, user mjwt.CustomClaim, cctvID string, input dto.CctvEditRequest) (*dto.Cctv, rest_err.APIError) {
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
	cctvEdited, err := c.daoC.EditCctv(ctx, data)
	if err != nil {
		return nil, err
	}

	errChan := make(chan rest_err.APIError, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	editUnit := func() {
		defer wg.Done()
		// DB
		_, err = c.daoG.EditUnit(ctx, cctvID, dto.GenUnitEditRequest{
			Category: category.Cctv,
			Name:     cctvEdited.Name,
			IP:       cctvEdited.IP,
			Branch:   cctvEdited.Branch,
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
				Category:       category.Cctv,
				Branch:         user.Branch,
				ParentID:       cctvID,
				ParentName:     cctvEdited.Name,
				Status:         "Edit",
				Problem:        "Detail Cctv diubah",
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
		restErr := rest_err.NewInternalServerError(fmt.Sprintf("Mengedit cctv berhasil namun : %s", errMessage), errorPartial)
		logger.Error(restErr.Message(), errorPartial)
		return nil, restErr
	}

	return cctvEdited, nil
}

func (c *cctvService) DeleteCctv(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := c.daoC.DeleteCctv(ctx, dto.FilterIDBranchCreateGte{
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

// DisableCctv if value true , cctv will disabled
func (c *cctvService) DisableCctv(ctx context.Context, cctvID string, user mjwt.CustomClaim, value bool) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(cctvID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable cctv
	cctv, err := c.daoC.DisableCctv(ctx, oid, user, value)
	if err != nil {
		return nil, err
	}

	// set disable enable gen_unit
	_, err = c.daoG.DisableUnit(ctx, oid.Hex(), value)
	if err != nil {
		return nil, err
	}

	return cctv, nil
}

// PutImage memasukkan lokasi file (path) ke dalam database cctv dengan mengecek kesesuaian branch
func (c *cctvService) PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	cctv, err := c.daoC.UploadImage(ctx, oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return cctv, nil
}

func (c *cctvService) GetCctvByID(ctx context.Context, cctvID string, branchIfSpecific string) (*dto.Cctv, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(cctvID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// kembalian dari golang channel
	type resultCctv struct {
		data *dto.Cctv
		err  rest_err.APIError
	}

	type resultGeneral struct {
		data *dto.GenUnitResponse
		err  rest_err.APIError
	}

	resultCctvChan := make(chan resultCctv)
	resultGeneralChan := make(chan resultGeneral)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// DB
		cctv, err := c.daoC.GetCctvByID(ctx, oid, branchIfSpecific)
		resultCctvChan <- resultCctv{
			data: cctv,
			err:  err,
		}
	}()

	go func() {
		defer wg.Done()
		// DB
		cctv, err := c.daoG.GetUnitByID(ctx, cctvID, branchIfSpecific)
		resultGeneralChan <- resultGeneral{
			data: cctv,
			err:  err,
		}
	}()

	go func() {
		wg.Wait()
		close(resultCctvChan)
		close(resultGeneralChan)
	}()

	cctvDetail := <-resultCctvChan
	if cctvDetail.err != nil {
		return nil, cctvDetail.err
	}

	cctvExtraData := <-resultGeneralChan
	if cctvExtraData.err != nil {
		return nil, cctvExtraData.err
	}

	cctvData := cctvDetail.data
	cctvData.Extra.Cases = cctvExtraData.data.Cases
	cctvData.Extra.LastPing = cctvExtraData.data.LastPing
	cctvData.Extra.CasesSize = cctvExtraData.data.CasesSize
	cctvData.Extra.PingsState = cctvExtraData.data.PingsState

	return cctvData, nil
}

func (c *cctvService) FindCctv(ctx context.Context, filter dto.FilterBranchLocIPNameDisable) (dto.CctvResponseMinList, dto.GenUnitResponseList, rest_err.APIError) {
	// cek apakah ip address valid, jika valid maka set filter.FilterName ke kosong supaya pencarian berdasarkan IP
	if filter.FilterIP != "" {
		if net.ParseIP(filter.FilterIP) == nil {
			return nil, nil, rest_err.NewBadRequestError("IP Address tidak valid")
		}
		filter.FilterName = ""
	}

	// wrap golang channel
	type resultCctv struct {
		data dto.CctvResponseMinList
		err  rest_err.APIError
	}

	// wrap golang channel
	type resultGeneral struct {
		data dto.GenUnitResponseList
		err  rest_err.APIError
	}

	cctvListChan := make(chan resultCctv)
	generalListChan := make(chan resultGeneral)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		cctvList, err := c.daoC.FindCctv(ctx, filter)
		cctvListChan <- resultCctv{
			data: cctvList,
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
			Category: category.Cctv,
			Pings:    true,
			Name:     "",
		})
		generalListChan <- resultGeneral{
			data: generalList,
			err:  err,
		}
	}()

	cctvList := <-cctvListChan
	if cctvList.err != nil {
		return nil, nil, cctvList.err
	}

	generalList := <-generalListChan
	if generalList.err != nil {
		return nil, nil, generalList.err
	}

	filterGeneralList(&generalList.data)
	return cctvList.data, generalList.data, nil
}

// MergeCctv akan menggabungkan cctv 1 ke cctv 2 dan menghapus cctv 1
// mungkin memerlukan bantuan developer untuk menghapus keberadaan cctv 1 pada checklist bulanan dan triwulanan
func (c *cctvService) MergeCctv(ctx context.Context, cctvID1, cctvID2 string) (*string, rest_err.APIError) {
	oid1, errT := primitive.ObjectIDFromHex(cctvID1)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	oid2, errT := primitive.ObjectIDFromHex(cctvID2)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// mendapatkan detail cctvID2
	cctv1Detail, apiErr := c.daoC.GetCctvByID(ctx, oid1, "")
	if apiErr != nil {
		return nil, apiErr
	}

	// mendapatkan detail cctvID2
	cctv2Detail, apiErr := c.daoC.GetCctvByID(ctx, oid2, "")
	if apiErr != nil {
		return nil, apiErr
	}

	// cek history in cctvID1 , untuk setiap historynya maka masukkan ke cctvID2
	history1, apiErr := c.daoH.FindHistoryForParent(ctx, cctvID1)
	if apiErr != nil {
		return nil, apiErr
	}

	if len(history1) > 0 {
		for _, history := range history1 {
			generateHistID := primitive.NewObjectID()
			_, apiErr = c.daoH.InsertHistory(ctx, dto.History{
				Version:        history.Version,
				ID:             generateHistID,
				CreatedAt:      history.CreatedAt,
				CreatedBy:      history.CreatedBy,
				UpdatedAt:      history.UpdatedAt,
				UpdatedBy:      history.UpdatedBy,
				Category:       history.Category,
				Branch:         history.Branch,
				ParentID:       cctvID2,
				ParentName:     cctv2Detail.Name,
				Status:         history.Status,
				Problem:        history.Problem,
				ProblemResolve: history.ProblemResolve,
				CompleteStatus: history.CompleteStatus,
				DateStart:      history.DateStart,
				DateEnd:        history.DateEnd,
				Tag:            history.Tag,
				Image:          history.Image,
				Updates:        history.Updates,
			}, false)

			if apiErr != nil {
				return nil, apiErr
			}

			historyIsComplete := history.CompleteStatus == enum.HComplete
			historyIsInfo := history.CompleteStatus == enum.HInfo
			historyIsDataInfo := history.CompleteStatus == enum.HDataInfo
			if !(historyIsComplete || historyIsInfo || historyIsDataInfo) {
				// DB
				_, apiErr = c.daoG.InsertCase(ctx, dto.GenUnitCaseRequest{
					UnitID:       cctvID2,
					FilterBranch: cctv2Detail.Branch,
					CaseID:       generateHistID.Hex(), // gunakan History id sebagai caseID
					CaseNote:     fmt.Sprintf("#%s# %s : %s", enum.GetProgressString(history.CompleteStatus), history.Status, history.Problem),
				})
				if apiErr != nil {
					return nil, apiErr
				}
			}
		}
	}

	// hapus cctvID 1
	// hapus gen Unit cctvID 1
	_, apiErr = c.daoC.DeleteCctv(ctx, dto.FilterIDBranchCreateGte{
		FilterID:        oid1,
		FilterBranch:    cctv1Detail.Branch,
		FilterCreateGTE: 0,
	})
	if apiErr != nil {
		return nil, apiErr
	}

	apiErr = c.daoG.DeleteUnit(ctx, cctvID1)
	if apiErr != nil {
		return nil, apiErr
	}

	msg := fmt.Sprintf(
		"Cctv dengan id %s berhasil digabungkan dengan id %s, cctv %s telah dihapus dari database, mohon melakukan pengecekan ulang pada checklist bulanan dan triwulanan",
		cctvID1, cctvID2, cctvID1)
	return &msg, nil
}

// hanya mereturn unit yang memiliki case atau sedang down.
func filterGeneralList(data *dto.GenUnitResponseList) {
	temp := dto.GenUnitResponseList{}
	for _, gen := range *data {
		if gen.CasesSize > 0 || gen.LastPing == "DOWN" {
			temp = append(temp, gen)
		}
	}
	*data = temp
}
