package handler

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

func NewPingHandler() *pingHandler {
	return &pingHandler{}
}

type pingHandler struct{}

//Ping mengembalikan pong untuk keperluan pengecekan ketersediaan server
func (p *pingHandler) Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"msg": "PONG!", "time": time.Now()})
}
