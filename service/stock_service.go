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

	//DB
	insertedID, err := s.daoS.InsertStock(data)
	if err != nil {
		return nil, err
	}

	//DB
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
		CompleteStatus: enum.HComplete,
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
