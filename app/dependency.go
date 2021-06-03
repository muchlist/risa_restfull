package app

import (
	"github.com/muchlist/risa_restfull/dao/cctvdao"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/checkitemdao"
	"github.com/muchlist/risa_restfull/dao/computerdao"
	"github.com/muchlist/risa_restfull/dao/genunitdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/improvedao"
	"github.com/muchlist/risa_restfull/dao/stockdao"
	"github.com/muchlist/risa_restfull/dao/userdao"
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
	userDao      = userdao.NewUserDao()
	genUnitDao   = genunitdao.NewGenUnitDao()
	historyDao   = historydao.NewHistoryDao()
	cctvDao      = cctvdao.NewCctvDao()
	stockDao     = stockdao.NewStockDao()
	checkItemDao = checkitemdao.NewCheckItemDao()
	checkDao     = checkdao.NewCheckDao()
	improveDao   = improvedao.NewImproveDao()
	computerDao  = computerdao.NewComputerDao()

	// Service
	userService      = service.NewUserService(userDao, cryptoUtils, jwt)
	genUnitService   = service.NewGenUnitService(genUnitDao)
	historyService   = service.NewHistoryService(historyDao, genUnitDao)
	cctvService      = service.NewCctvService(cctvDao, historyDao, genUnitDao)
	stockService     = service.NewStockService(stockDao, historyDao)
	checkItemService = service.NewCheckItemService(checkItemDao)
	checkService     = service.NewCheckService(checkDao, checkItemDao, genUnitDao, historyService)
	improveService   = service.NewImproveService(improveDao)
	computerService  = service.NewComputerService(computerDao, historyDao, genUnitDao)

	// Controller or Handler
	pingHandler      = handler.NewPingHandler()
	optionHandler    = handler.NewOptionHandler()
	userHandler      = handler.NewUserHandler(userService)
	genUnitHandler   = handler.NewGenUnitHandler(genUnitService)
	historyHandler   = handler.NewHistoryHandler(historyService)
	cctvHandler      = handler.NewCctvHandler(cctvService)
	stockHandler     = handler.NewStockHandler(stockService)
	checkItemHandler = handler.NewCheckItemHandler(checkItemService)
	checkHandler     = handler.NewCheckHandler(checkService)
	improveHandler   = handler.NewImproveHandler(improveService)
	computerHandler  = handler.NewComputerHandler(computerService)
)
