package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/constants/branches"
	"github.com/muchlist/risa_restfull/constants/checktype"
	"github.com/muchlist/risa_restfull/constants/location"
	"strings"
)

func NewOptionHandler() *optionHandler {
	return &optionHandler{}
}

type optionHandler struct{}

// OptCreateCheckItem mengembalikan location, dan type
func (o *optionHandler) OptCreateCheckItem(c *fiber.Ctx) error {
	branch := c.Query("branch")
	var optLocation []string

	if branch != "" {
		optLocation = location.GetLocationAvailableFrom(strings.ToUpper(branch))
	} else {
		optLocation = location.GetLocationAvailable()
	}

	optType := checktype.GetCheckTypeAvailable()
	options := fiber.Map{
		"location": optLocation,
		"type":     optType,
	}
	return c.JSON(options)
}

// OptBranch mengembalikan branch yang tersedia
func (o *optionHandler) OptBranch(c *fiber.Ctx) error {
	optBranch := branches.GetBranchesAvailable()
	options := fiber.Map{
		"branch": optBranch,
	}
	return c.JSON(options)
}
