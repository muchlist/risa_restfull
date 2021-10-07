package app

import (
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/dao/altaicheckdao"
	"github.com/muchlist/risa_restfull/dao/altaiphycheckdao"
	"github.com/muchlist/risa_restfull/dao/cctvdao"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/checkitemdao"
	"github.com/muchlist/risa_restfull/dao/computerdao"
	"github.com/muchlist/risa_restfull/dao/configcheckdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/improvedao"
	"github.com/muchlist/risa_restfull/dao/otherdao"
	"github.com/muchlist/risa_restfull/dao/pendingreportdao"
	"github.com/muchlist/risa_restfull/dao/reportdao"
	"github.com/muchlist/risa_restfull/dao/speedtestdao"
	"github.com/muchlist/risa_restfull/dao/stockdao"
	"github.com/muchlist/risa_restfull/dao/userdao"
	"github.com/muchlist/risa_restfull/dao/vendorcheckdao"
	"github.com/muchlist/risa_restfull/dao/venphycheckdao"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/crypt"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

var (
	userService          service.UserServiceAssumer
	genUnitService       service.GenUnitServiceAssumer
	historyService       service.HistoryServiceAssumer
	cctvService          service.CctvServiceAssumer
	stockService         service.StockServiceAssumer
	checkItemService     service.CheckItemServiceAssumer
	checkService         service.CheckServiceAssumer
	improveService       service.ImproveServiceAssumer
	computerService      service.ComputerServiceAssumer
	otherService         service.OtherServiceAssumer
	vendorCheckService   service.VendorCheckServiceAssumer
	altaiCheckService    service.AltaiCheckServiceAssumer
	venPhyCheckService   service.VenPhyCheckServiceAssumer
	altaiPhyCheckService service.AltaiPhyCheckServiceAssumer
	configCheckService   service.ConfigCheckServiceAssumer
	speedService         service.SpeedTestServiceAssumer
	reportService        service.ReportServiceAssumer
	prService            service.PRServiceAssumer
)

func setupDependency() {
	// Utils
	cryptoUtils := crypt.NewCrypto()
	jwt := mjwt.NewJwt()

	// Dao
	userDao := userdao.NewUserDao()
	genUnitDao := genunitdao.NewGenUnitDao()
	historyDao := historydao.NewHistoryDao()
	cctvDao := cctvdao.NewCctvDao()
	stockDao := stockdao.NewStockDao()
	checkItemDao := checkitemdao.NewCheckItemDao()
	checkDao := checkdao.NewCheckDao()
	improveDao := improvedao.NewImproveDao()
	computerDao := computerdao.NewComputerDao()
	otherDao := otherdao.NewOtherDao()
	vendorCheckDao := vendorcheckdao.NewVendorCheckDao()
	altaiCheckDao := altaicheckdao.NewAltaiCheckDao()
	venPhyCheckDao := venphycheckdao.NewVenPhyCheckDao()
	altaiPhyCheckDao := altaiphycheckdao.NewAltaiPhyCheckDao()
	configCheckDao := configcheckdao.NewConfigCheckDao()
	speedDao := speedtestdao.NewSpeedTestDao()
	pdfDao := reportdao.NewPdfDao()
	prDao := pendingreportdao.NewPR()

	// api client
	fcmClient := fcm.NewFcmClient()

	// Service
	userService = service.NewUserService(userDao, cryptoUtils, jwt)
	genUnitService = service.NewGenUnitService(genUnitDao, userDao, fcmClient)
	historyService = service.NewHistoryService(historyDao, genUnitDao, userDao, fcmClient)
	cctvService = service.NewCctvService(cctvDao, historyDao, genUnitDao)
	stockService = service.NewStockService(stockDao, historyDao, userDao, fcmClient)
	checkItemService = service.NewCheckItemService(checkItemDao)
	checkService = service.NewCheckService(checkDao, checkItemDao, genUnitDao, historyService)
	improveService = service.NewImproveService(improveDao)
	computerService = service.NewComputerService(computerDao, historyDao, genUnitDao)
	otherService = service.NewOtherService(otherDao, historyDao, genUnitDao)
	vendorCheckService = service.NewVendorCheckService(vendorCheckDao, genUnitDao, cctvDao, historyService)
	altaiCheckService = service.NewAltaiCheckService(altaiCheckDao, genUnitDao, otherDao, historyService)
	venPhyCheckService = service.NewVenPhyCheckService(venPhyCheckDao, genUnitDao, cctvDao, historyService)
	altaiPhyCheckService = service.NewAltaiPhyCheckService(altaiPhyCheckDao, genUnitDao, otherDao, historyService)
	configCheckService = service.NewConfigCheckService(configCheckDao, genUnitDao, otherDao, historyService)
	speedService = service.NewSpeedTestService(speedDao)
	prService = service.NewPRService(prDao, genUnitDao)
	reportService = service.NewReportService(service.ReportParams{
		History:       historyDao,
		CheckIT:       checkDao,
		CheckCCTV:     vendorCheckDao,
		CheckCCTVPhy:  venPhyCheckDao,
		CheckAltai:    altaiCheckDao,
		CheckAltaiPhy: altaiPhyCheckDao,
		CheckConfig:   configCheckDao,
		Stock:         stockDao,
		Pdf:           pdfDao,
	})
}
