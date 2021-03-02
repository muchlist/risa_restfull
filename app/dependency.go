package app

import (
	"github.com/muchlist/risa_restfull/dao/cctv_dao"
	"github.com/muchlist/risa_restfull/dao/check_dao"
	"github.com/muchlist/risa_restfull/dao/check_item_dao"
	"github.com/muchlist/risa_restfull/dao/gen_unit_dao"
	"github.com/muchlist/risa_restfull/dao/history_dao"
	"github.com/muchlist/risa_restfull/dao/stock_dao"
	"github.com/muchlist/risa_restfull/dao/user_dao"
	"github.com/muchlist/risa_restfull/handler"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/crypt"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

var (
	//Utils
	cryptoUtils = crypt.NewCrypto()
	jwt         = mjwt.NewJwt()

	//Dao
	userDao      = user_dao.NewUserDao()
	genUnitDao   = gen_unit_dao.NewGenUnitDao()
	historyDao   = history_dao.NewHistoryDao()
	cctvDao      = cctv_dao.NewCctvDao()
	stockDao     = stock_dao.NewStockDao()
	checkItemDao = check_item_dao.NewCheckItemDao()
	checkDao     = check_dao.NewCheckDao()

	//Service
	userService      = service.NewUserService(userDao, cryptoUtils, jwt)
	genUnitService   = service.NewGenUnitService(genUnitDao)
	historyService   = service.NewHistoryService(historyDao, genUnitDao)
	cctvService      = service.NewCctvService(cctvDao, historyDao, genUnitDao)
	stockService     = service.NewStockService(stockDao, historyDao)
	checkItemService = service.NewCheckItemService(checkItemDao)
	checkService     = service.NewCheckService(checkDao, checkItemDao, genUnitDao, historyService)

	//Controller or Handler
	pingHandler      = handler.NewPingHandler()
	userHandler      = handler.NewUserHandler(userService)
	genUnitHandler   = handler.NewGenUnitHandler(genUnitService)
	historyHandler   = handler.NewHistoryHandler(historyService)
	cctvHandler      = handler.NewCctvHandler(cctvService)
	stockHandler     = handler.NewStockHandler(stockService)
	checkItemHandler = handler.NewCheckItemHandler(checkItemService)
	checkHandler     = handler.NewCheckHandler(checkService)
)
