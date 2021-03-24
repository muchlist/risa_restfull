package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

func NewUserHandler(userService service.UserServiceAssumer) *userHandler {
	return &userHandler{
		service: userService,
	}
}

type userHandler struct {
	service service.UserServiceAssumer
}

// Get menampilkan user berdasarkan ID (bukan email)
func (u *userHandler) Get(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	user, apiErr := u.service.GetUser(userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// GetProfile mengembalikan user yang sedang login
func (u *userHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	user, apiErr := u.service.GetUserByID(claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// Register menambahkan user
func (u *userHandler) Register(c *fiber.Ctx) error {
	var user dto.UserRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := u.service.InsertUser(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Register berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// Find menampilkan list user
func (u *userHandler) Find(c *fiber.Ctx) error {
	userList, apiErr := u.service.FindUsers()
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userList})
}

// Edit mengedit user oleh admin
func (u *userHandler) Edit(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	//if err := validation.Validate(userID,
	//	is.Email,
	//); err != nil {
	//	apiErr := rest_err.NewBadRequestError(err.Error())
	//	return c.Status(apiErr.Status()).JSON(apiErr)
	//}

	var user dto.UserEditRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := u.service.EditUser(userID, user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// Delete menghapus user, idealnya melalui middleware is_admin
func (u *userHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	userIDParams := c.Params("user_id")

	if claims.Identity == userIDParams {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := u.service.DeleteUser(userIDParams)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("user %s berhasil dihapus", userIDParams)})
}

// ChangePassword mengganti password pada user sendiri
func (u *userHandler) ChangePassword(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var user dto.UserChangePasswordRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	//mengganti user id dengan user aktif
	user.ID = claims.Identity

	apiErr := u.service.ChangePassword(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"error": apiErr, "data": "Password berhasil diubah!"})
}

// ResetPassword mengganti password oleh admin pada user tertentu
func (u *userHandler) ResetPassword(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	data := dto.UserChangePasswordRequest{
		ID:          userID,
		NewPassword: "Password",
	}

	apiErr := u.service.ResetPassword(data)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("Password user %s berhasil di reset!", c.Params("user_id"))})
}

// Login login
func (u *userHandler) Login(c *fiber.Ctx) error {
	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := login.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: - | validate | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := u.service.Login(login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// RefreshToken
func (u *userHandler) RefreshToken(c *fiber.Ctx) error {
	var payload dto.UserRefreshTokenRequest
	if err := c.BodyParser(&payload); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := payload.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := u.service.Refresh(payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
// sesuai authorisasi aktif. File disimpan di folder static/images dengan nama file == jwt.identity alias username
func (u *userHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	pathInDB, apiErr := saveImage(c, *claims, "avatar", claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	usersResult, apiErr := u.service.PutAvatar(claims.Identity, pathInDB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": usersResult})
}
