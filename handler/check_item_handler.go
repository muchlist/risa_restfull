package handler

import "github.com/muchlist/risa_restfull/service"

func NewCheckItemHandler(checkItemService service.CheckItemServiceAssumer) *checkItemHandler {
	return &checkItemHandler{
		service: checkItemService,
	}
}

type checkItemHandler struct {
	service service.CheckItemServiceAssumer
}
