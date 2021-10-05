package service

import (
	"context"
	"fmt"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/stockdao"
	"github.com/muchlist/risa_restfull/dao/userdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewStockService(stockDao stockdao.StockDaoAssumer,
	histDao historydao.HistorySaver, userDao userdao.UserDaoAssumer,
	fcmClient fcm.ClientAssumer) StockServiceAssumer {
	return &stockService{
		daoS:      stockDao,
		daoH:      histDao,
		daoU:      userDao,
		fcmClient: fcmClient,
	}
}

type stockService struct {
	daoS      stockdao.StockDaoAssumer
	daoH      historydao.HistorySaver
	daoU      userdao.UserDaoAssumer
	fcmClient fcm.ClientAssumer
}
type StockServiceAssumer interface {
	InsertStock(ctx context.Context, user mjwt.CustomClaim, input dto.StockRequest) (*string, rest_err.APIError)
	EditStock(ctx context.Context, user mjwt.CustomClaim, stockID string, input dto.StockEditRequest) (*dto.Stock, rest_err.APIError)
	DeleteStock(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError
	DisableStock(ctx context.Context, stockID string, user mjwt.CustomClaim, value bool) (*dto.Stock, rest_err.APIError)
	PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Stock, rest_err.APIError)
	ChangeQtyStock(ctx context.Context, user mjwt.CustomClaim, stockID string, data dto.StockChangeRequest) (*dto.Stock, rest_err.APIError)
	GetStockByID(ctx context.Context, stockID string, branchIfSpecific string) (*dto.Stock, rest_err.APIError)
	FindStock(ctx context.Context, filter dto.FilterBranchNameCatDisable) (dto.StockResponseMinList, rest_err.APIError)
	FindNeedReStock(ctx context.Context, branch string) (dto.StockResponseMinList, rest_err.APIError)
	FindNeedReStock2(ctx context.Context, filter dto.FilterBranchCatDisable) ([]dto.Stock, rest_err.APIError)
}

func (s *stockService) InsertStock(ctx context.Context, user mjwt.CustomClaim, input dto.StockRequest) (*string, rest_err.APIError) {
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
	insertedID, err := s.daoS.InsertStock(ctx, data)
	if err != nil {
		return nil, err
	}
	isVendor := sfunc.InSlice(roles.RoleVendor, user.Roles)
	// DB
	_, err = s.daoH.InsertHistory(ctx, dto.History{
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
	}, isVendor)
	if err != nil {
		logger.Error("Berhasil membuat stock namun gagal membuat History (InsertStock)", err)
		errPlus := rest_err.NewInternalServerError(fmt.Sprintf("galat : stock berhasil ditambahkan -> %s", err.Message()), err)
		return nil, errPlus
	}

	return insertedID, nil
}

func (s *stockService) EditStock(ctx context.Context, user mjwt.CustomClaim, stockID string, input dto.StockEditRequest) (*dto.Stock, rest_err.APIError) {
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
	stockEdited, err := s.daoS.EditStock(ctx, data)
	if err != nil {
		return nil, err
	}

	return stockEdited, nil
}

func (s *stockService) DeleteStock(ctx context.Context, user mjwt.CustomClaim, id string) rest_err.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Dokumen yang dibuat sehari sebelumnya masih bisa dihapus
	timeMinusOneDay := time.Now().AddDate(0, 0, -1)
	// DB
	_, err := s.daoS.DeleteStock(ctx, dto.FilterIDBranchCreateGte{
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
func (s *stockService) DisableStock(ctx context.Context, stockID string, user mjwt.CustomClaim, isDisable bool) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// set disable enable stock
	stock, err := s.daoS.DisableStock(ctx, oid, user, isDisable)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

// PutImage memasukkan lokasi file (path) ke dalam database stock dengan mengecek kesesuaian branch
func (s *stockService) PutImage(ctx context.Context, user mjwt.CustomClaim, id string, imagePath string) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	stock, err := s.daoS.UploadImage(ctx, oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *stockService) ChangeQtyStock(ctx context.Context, user mjwt.CustomClaim, stockID string, data dto.StockChangeRequest) (*dto.Stock, rest_err.APIError) {
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
	stockEdited, err := s.daoS.ChangeQtyStock(ctx, filter, incDec)
	if err != nil {
		return nil, err
	}

	// Filling Data History
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
				-data.Qty,
				stockEdited.Unit,
				data.Note,
				stockEdited.Qty,
				stockEdited.Unit,
			)
		} else {
			history.Problem = fmt.Sprintf("Mengurangi stok (%d) %s : %s - sisa stok %d %s",
				-data.Qty,
				stockEdited.Unit,
				data.Note,
				stockEdited.Qty,
				stockEdited.Unit,
			)
		}
	}

	isVendor := sfunc.InSlice(roles.RoleVendor, user.Roles)
	// DB
	_, err = s.daoH.InsertHistory(ctx, history, isVendor)
	if err != nil {
		logger.Error("Berhasil membuat stock namun gagal membuat History (ChangeQtyStock)", err)
		errPlus := rest_err.NewInternalServerError(fmt.Sprintf("galat : stock berhasil diuubah -> %s", err.Message()), err)
		return nil, errPlus
	}

	go func() {
		users, err := s.daoU.FindUser(ctx, user.Branch)
		if err != nil {
			logger.Error("mendapatkan user gagal saat menambahkan fcm (EDIT STOCK)", err)
		}

		var tokens []string
		for _, u := range users {
			if u.ID == user.Identity {
				continue
			}
			// tidak dikirimkan ke user vendor
			if sfunc.InSlice(roles.RoleVendor, u.Roles) {
				continue
			}
			tokens = append(tokens, u.FcmToken)
		}
		// firebase
		s.fcmClient.SendMessage(fcm.Payload{
			Title:          fmt.Sprintf("Stok %s berubah", stockEdited.Name),
			Message:        history.Problem,
			ReceiverTokens: tokens,
		})
	}()

	return stockEdited, nil
}

func (s *stockService) GetStockByID(ctx context.Context, stockID string, branchIfSpecific string) (*dto.Stock, rest_err.APIError) {
	oid, errT := primitive.ObjectIDFromHex(stockID)
	if errT != nil {
		return nil, rest_err.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	stock, err := s.daoS.GetStockByID(ctx, oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *stockService) FindStock(ctx context.Context, filter dto.FilterBranchNameCatDisable) (dto.StockResponseMinList, rest_err.APIError) {
	stockList, err := s.daoS.FindStock(ctx, filter)
	if err != nil {
		return nil, err
	}
	return stockList, nil
}

func (s *stockService) FindNeedReStock(ctx context.Context, branch string) (dto.StockResponseMinList, rest_err.APIError) {
	stockList, err := s.daoS.FindStock(ctx, dto.FilterBranchNameCatDisable{
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

func (s *stockService) FindNeedReStock2(ctx context.Context, filter dto.FilterBranchCatDisable) ([]dto.Stock, rest_err.APIError) {
	stockList, err := s.daoS.FindStockNeedRestock(ctx, filter)
	if err != nil {
		return nil, err
	}
	return stockList, nil
}
