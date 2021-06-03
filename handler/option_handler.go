package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/constants/branches"
	"github.com/muchlist/risa_restfull/constants/checktype"
	"github.com/muchlist/risa_restfull/constants/hwlist"
	"github.com/muchlist/risa_restfull/constants/location"
	"github.com/muchlist/risa_restfull/constants/stocklist"
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

// OptCreateStock mengembalikan stock category
func (o *optionHandler) OptCreateStock(c *fiber.Ctx) error {
	stockCategory := stocklist.GetStockCategoryAvailable()
	options := fiber.Map{
		"category": stockCategory,
	}
	return c.JSON(options)
}

// OptCreateCctv mengembalikan cctv type dan lokasi tersedia
func (o *optionHandler) OptCreateCctv(c *fiber.Ctx) error {
	branch := c.Query("branch")
	var optLocation []string

	if branch != "" {
		optLocation = location.GetLocationAvailableFrom(strings.ToUpper(branch))
	} else {
		optLocation = location.GetLocationAvailable()
	}

	cctvType := hwlist.GetCctvTypeAvailable()
	options := fiber.Map{
		"location": optLocation,
		"type":     cctvType,
	}
	return c.JSON(options)
}

// OptCreateCctv mengembalikan cctv type dan lokasi tersedia
func (o *optionHandler) OptCreateComputer(c *fiber.Ctx) error {
	branch := c.Query("branch")
	var optLocation []string

	if branch != "" {
		optLocation = location.GetLocationAvailableFrom(strings.ToUpper(branch))
	} else {
		optLocation = location.GetLocationAvailable()
	}
	division := location.GetDivisionAvailable()
	computerType := hwlist.GetComputerTypeAvailable()
	os := hwlist.GetPCOSAvailable()
	processor := hwlist.GetPCProcessor()
	hardisk := hwlist.GetPCHDD()
	ram := hwlist.GetPCRam()

	options := fiber.Map{
		"location":  optLocation,
		"division":  division,
		"type":      computerType,
		"os":        os,
		"processor": processor,
		"hardisk":   hardisk,
		"ram":       ram,
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
