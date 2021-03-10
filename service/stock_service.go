package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dao/history_dao"
	"github.com/muchlist/risa_restfull/dao/stock_dao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewStockService(stockDao stock_dao.StockDaoAssumer,
	histDao history_dao.HistoryDaoAssumer) StockServiceAssumer {
	return &stockService{
		daoS: stockDao,
		daoH: histDao,
	}
}

type stockService struct {
	daoS stock_dao.StockDaoAssumer
	daoH history_dao.HistoryDaoAssumer
}
type StockServiceAssumer interface {
	InsertStock(user mjwt.CustomClaim, input dto.StockRequest) (*string, rest_err.APIError)
	EditStock(user mjwt.CustomClaim, stockID string, input dto.StockEditRequest) (*dto.Stock, rest_err.APIError)
	DeleteStock(user mjwt.CustomClaim, id string) rest_err.APIError
	DisableStock(stockID string, user mjwt.CustomClaim, value bool) (*dto.Stock, rest_err.APIError)
	PutImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Stock, rest_err.APIError)
	ChangeQtyStock(user mjwt.CustomClaim, stockID string, data dto.StockChangeRequest) (*dto.Stock, rest_err.APIError)
	GetStockByID(stockID string, branchIfSpecific string) (*dto.Stock, rest_err.APIError)
	FindStock(filter dto.FilterBranchNameCatDisable) (dto.StockResponseMinList, rest_err.APIError)
	FindNeedReStock(branch string) (dto.StockResponseMinList, rest_err.APIError)
}

func (s *stockService) InsertStock(user mjwt.CustomClaim, input dto.StockRequest) (*string, rest_err.APIError) {
	// Filling data
	// Ketika membuat stock juga menambahkan increment field untuk pertama kali
	timeNow := time.Now().Unix()
	inc := dto.StockChange{
		DummyID:  time.Now().UnixNano(),
		Author:   user.Name,
		Qty:      input.Qty,
		BaNumber: "",
		Note:     "inisiasi",
		Time:     timeNow,
	}

	oidGenerated := primitive.NewObjectID()
	data := dto.Stock{
		ID:            oidGenerated,
		CreatedAt:     timeNow,
		CreatedBy:     user.Name,
		CreatedByID:   user.Identity,
		UpdatedAt:     timeNow,
		UpdatedBy:     user.Name,
		UpdatedByID:   user.Identity,
		Branch:        user.Branch,
		Name:          input.Name,
		Disable:       false,
		StockCategory: input.StockCategory,
		Unit:          input.Unit,
		Qty:           input.Qty,
		Location:      input.Location,
		Threshold:     input.Threshold,
		Increment:     []dto.StockChange{inc},
		Decrement:     nil,
		Tag:           input.Tag,
		Image:         "",
		Note:          input.Note,
	}

	// DB
	insertedID, err := s.daoS.InsertStock(data)
	if err != nil {
		return nil, err
	}

	// DB
	_, err = s.daoH.InsertHistory(dto.History{
		CreatedAt:      timeNow,
		CreatedBy:      user.Name,
		CreatedByID:    user.Identity,
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Category:       category.Stock,
		Branch:         user.Branch,
		ID:             primitive.NewObjectID(),
		ParentID:       oidGenerated.Hex(),
		ParentName:     input.Name,
		Status:         "Change",
		Problem:        fmt.Sprintf("Menambahkan stok : %d %s", input.Qty, input.Unit),
		ProblemResolve: "",
		CompleteStatus: enum.HInfo,
		DateStart:      timeNow,
		DateEnd:        timeNow,
		Tag:            []string{},
		Image:          "",
	})
	if err != nil {
		logger.Error("Berhasil membuat stock namun gagal membuat history (InsertStock)", err)
		errPlus := rest_err.NewInternalServerError(fmt.Sprintf("galat : stock berhasil ditambahkan -> %s", err.Message()), err)
		return nil, errPlus
	}

	return insertedID, nil
}

func (s *stockService) EditStock(user mjwt.CustomClaim, stockID string, input dto.StockEditRequest) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.StockEdit{
		ID:              oid,
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		StockCategory:   input.StockCategory,
		Unit:            input.Unit,
		Threshold:       input.Threshold,
		Name:            input.Name,
		Location:        input.Location,
		Tag:             input.Tag,
		Note:            input.Note,
	}

	// DB
	stockEdited, err := s.daoS.EditStock(data)
	if err != nil {
		return nil, err
	}

	return stockEdited, nil
}

func (s *stockService) DeleteStock(user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := s.daoS.DeleteStock(dto.FilterIDBranchCreateGte{
		FilterID:        oid,
		FilterBranch:    user.Branch,
		FilterCreateGTE: timeMinusOneDay.Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// DisableStock if value true , stock will disabled
func (s *stockService) DisableStock(stockID string, user mjwt.CustomClaim, isDisable bool) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable stock
	stock, err := s.daoS.DisableStock(oid, user, isDisable)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

// PutImage memasukkan lokasi file (path) ke dalam database stock dengan mengecek kesesuaian branch
func (s *stockService) PutImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	stock, err := s.daoS.UploadImage(oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *stockService) ChangeQtyStock(user mjwt.CustomClaim, stockID string, data dto.StockChangeRequest) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Filling data=
	timeNow := time.Now().Unix()
	incDec := dto.StockChange{
		DummyID:  time.Now().UnixNano(),
		Author:   user.Name,
		Qty:      data.Qty,
		BaNumber: data.BaNumber,
		Note:     data.Note,
	}
	if data.Time == 0 {
		incDec.Time = timeNow
	}

	filter := dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	}

	// DB
	stockEdited, err := s.daoS.ChangeQtyStock(filter, incDec)
	if err != nil {
		return nil, err
	}

	// Filling Data history
	history := dto.History{
		ID:             primitive.NewObjectID(),
		CreatedAt:      timeNow,
		CreatedBy:      user.Name,
		CreatedByID:    user.Identity,
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Category:       category.Stock,
		Branch:         user.Branch,
		ParentID:       stockID,
		ParentName:     stockEdited.Name,
		Status:         "Change",
		ProblemResolve: "",
		CompleteStatus: enum.HInfo,
		DateStart:      timeNow,
		DateEnd:        timeNow,
		Tag:            []string{},
		Image:          "",
	}
	if data.Qty > 0 {
		history.Problem = fmt.Sprintf("Menambahkan stok %d %s : %s", data.Qty, stockEdited.Unit, data.Note)
	} else {
		if stockEdited.Qty <= stockEdited.Threshold {
			history.Problem = fmt.Sprintf("Mengurangi stok (%d) %s : %s - sisa stok %d %s (perlu restock)",
				data.Qty,
				stockEdited.Unit,
				data.Note,
				stockEdited.Qty,
				stockEdited.Unit,
			)
		} else {
			history.Problem = fmt.Sprintf("Mengurangi stok (%d) %s : %s - sisa stok %d %s",
				data.Qty,
				stockEdited.Unit,
				data.Note,
				stockEdited.Qty,
				stockEdited.Unit,
			)
		}
	}

	// DB
	_, err = s.daoH.InsertHistory(history)
	if err != nil {
		logger.Error("Berhasil membuat stock namun gagal membuat history (ChangeQtyStock)", err)
		errPlus := rest_err.NewInternalServerError(fmt.Sprintf("galat : stock berhasil dirubah -> %s", err.Message()), err)
		return nil, errPlus
	}

	// IMPROVEMENT
	// jika stockEdited qty lebih kurang atau semadengan threshold
	// send notification to firebase

	return stockEdited, nil
}

func (s *stockService) GetStockByID(stockID string, branchIfSpecific string) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	stock, err := s.daoS.GetStockByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *stockService) FindStock(filter dto.FilterBranchNameCatDisable) (dto.StockResponseMinList, rest_err.APIError) {
	stockList, err := s.daoS.FindStock(filter)
	if err != nil {
		return nil, err
	}
	return stockList, nil
}

func (s *stockService) FindNeedReStock(branch string) (dto.StockResponseMinList, rest_err.APIError) {
	stockList, err := s.daoS.FindStock(dto.FilterBranchNameCatDisable{
		FilterBranch: branch,
	})

	if err != nil {
		return nil, err
	}

	var needRestockList dto.StockResponseMinList
	for _, v := range stockList {
		qtyIsZero := v.Qty == 0
		qtyNeedRestock := v.Qty <= v.Threshold

		if qtyIsZero || qtyNeedRestock {
			needRestockList = append(needRestockList, v)
		}
	}

	return needRestockList, nil
}
