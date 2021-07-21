package app

import (
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/dao/cctvdao"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/checkitemdao"
	"github.com/muchlist/risa_restfull/dao/computerdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/improvedao"
	"github.com/muchlist/risa_restfull/dao/otherdao"
	"github.com/muchlist/risa_restfull/dao/reportdao"
	"github.com/muchlist/risa_restfull/dao/speedtestdao"
	"github.com/muchlist/risa_restfull/dao/stockdao"
	"github.com/muchlist/risa_restfull/dao/userdao"
	"github.com/muchlist/risa_restfull/dao/vendorcheckdao"
	"github.com/muchlist/risa_restfull/handler"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/crypt"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

var (
	// Utils
	cryptoUtils = crypt.NewCrypto()
	jwt         = mjwt.NewJwt()

	// Dao
	userDao        = userdao.NewUserDao()
	genUnitDao     = genunitdao.NewGenUnitDao()
	historyDao     = historydao.NewHistoryDao()
	cctvDao        = cctvdao.NewCctvDao()
	stockDao       = stockdao.NewStockDao()
	checkItemDao   = checkitemdao.NewCheckItemDao()
	checkDao       = checkdao.NewCheckDao()
	improveDao     = improvedao.NewImproveDao()
	computerDao    = computerdao.NewComputerDao()
	otherDao       = otherdao.NewOtherDao()
	vendorCheckDao = vendorcheckdao.NewVendorCheckDao()
	speedDao       = speedtestdao.NewSpeedTestDao()
	pdfDao         = reportdao.NewPdfDao()

	// api client
	fcmClient = fcm.NewFcmClient()

	// Service
	userService        = service.NewUserService(userDao, cryptoUtils, jwt)
	genUnitService     = service.NewGenUnitService(genUnitDao, userDao, fcmClient)
	historyService     = service.NewHistoryService(historyDao, genUnitDao, userDao, fcmClient)
	cctvService        = service.NewCctvService(cctvDao, historyDao, genUnitDao)
	stockService       = service.NewStockService(stockDao, historyDao, userDao, fcmClient)
	checkItemService   = service.NewCheckItemService(checkItemDao)
	checkService       = service.NewCheckService(checkDao, checkItemDao, genUnitDao, historyService)
	improveService     = service.NewImproveService(improveDao)
	computerService    = service.NewComputerService(computerDao, historyDao, genUnitDao)
	otherService       = service.NewOtherService(otherDao, historyDao, genUnitDao)
	vendorCheckService = service.NewVendorCheckService(vendorCheckDao, genUnitDao, cctvDao, historyService)
	speedService       = service.NewSpeedTestService(speedDao)
	reportService      = service.NewReportService(historyDao, checkDao, pdfDao)

	// Controller or Handler
	pingHandler        = handler.NewPingHandler()
	optionHandler      = handler.NewOptionHandler()
	userHandler        = handler.NewUserHandler(userService)
	genUnitHandler     = handler.NewGenUnitHandler(genUnitService)
	historyHandler     = handler.NewHistoryHandler(historyService)
	cctvHandler        = handler.NewCctvHandler(cctvService)
	stockHandler       = handler.NewStockHandler(stockService)
	checkItemHandler   = handler.NewCheckItemHandler(checkItemService)
	checkHandler       = handler.NewCheckHandler(checkService)
	improveHandler     = handler.NewImproveHandler(improveService)
	computerHandler    = handler.NewComputerHandler(computerService)
	otherHandler       = handler.NewOtherHandler(otherService)
	vendorCheckHandler = handler.NewVendorCheckHandler(vendorCheckService)
	speedHandler       = handler.NewSpeedHandler(speedService)
	reportHandler      = handler.NewReportHandler(reportService)
)
