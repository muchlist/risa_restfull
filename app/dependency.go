package app

import (
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
	userDao = user_dao.NewUserDao()

	//Service
	userService = service.NewUserService(userDao, cryptoUtils, jwt)

	//Controller or Handler
	pingHandler = handler.NewPingHandler()
	userHandler = handler.NewUserHandler(userService)
)
