package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"time"
)

func NewPRHandler(prService service.PRServiceAssumer) *prHandler {
	return &prHandler{
		service: prService,
	}
}

type prHandler struct {
	service service.PRServiceAssumer
}

func (pr *prHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.PendingReportRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))

		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := pr.service.InsertPR(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	res := fmt.Sprintf("Menambahkan doc berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	var req dto.PendingReportEditRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res, apiErr := pr.service.EditPR(c.Context(), *claims, id, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) AddParticipant(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	userID := c.Params("userid")

	res, apiErr := pr.service.AddParticipant(c.Context(), *claims, id, userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) AddApprover(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	userID := c.Params("userid")

	res, apiErr := pr.service.AddApprover(c.Context(), *claims, id, userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) RemoveParticipant(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	userID := c.Params("userid")

	res, apiErr := pr.service.RemoveParticipant(c.Context(), *claims, id, userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) RemoveApprover(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	userID := c.Params("userid")

	res, apiErr := pr.service.RemoveApprover(c.Context(), *claims, id, userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) SendToSignMode(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	res, apiErr := pr.service.SendToSigningMode(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) SendToDraftMode(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	res, apiErr := pr.service.SendToDraftMode(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) SigningDoc(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	res, apiErr := pr.service.SignDocument(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID && branch ada
	_, apiErr := pr.service.GetPRByID(c.Context(), id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%s%v", id, time.Now().Unix())
	// simpan image
	pathInDB, apiErr := saveImage(c, *claims, "document", randomName, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	docResult, apiErr := pr.service.PutImage(c.Context(), *claims, id, pathInDB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": docResult})
}

func (pr *prHandler) DeleteImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	pathImage := fmt.Sprintf("image/document/%s", c.Params("image"))

	// cek apakah ID violation && branch ada
	violation, apiErr := pr.service.GetPRByID(c.Context(), id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// untuk sementara gambar tidak dihapus dari system
	imageExist := sfunc.InSlice(pathImage, violation.Images)
	if !imageExist {
		apiErr := rest_err.NewBadRequestError("gambar tidak ditemukan")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	docResult, apiErr := pr.service.DeleteImage(c.Context(), *claims, id, pathImage)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": docResult})
}

func (pr *prHandler) Get(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	branch := c.Query("branch")

	if branch == "" {
		branch = claims.Branch
	}

	res, apiErr := pr.service.GetPRByID(c.Context(), id, branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (pr *prHandler) Find(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	var filter dto.FilterFindPendingReport
	filter.FilterBranch = c.Query("branch")
	filter.Limit = int64(sfunc.StrToInt(c.Query("limit"), 100))
	filter.CompleteStatus = c.Query("complete")
	filter.FilterTitle = c.Query("title")

	res, apiErr := pr.service.FindDocs(c.Context(), *claims, filter)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": res})
}
