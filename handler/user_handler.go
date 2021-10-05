package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"time"
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
func (usr *userHandler) Get(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	user, apiErr := usr.service.GetUser(c.Context(), userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// GetProfile mengembalikan user yang sedang login
func (usr *userHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	user, apiErr := usr.service.GetUserByID(c.Context(), claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// Register menambahkan user
func (usr *userHandler) Register(c *fiber.Ctx) error {
	var user dto.UserRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := usr.service.InsertUser(c.Context(), user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Register berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// Find menampilkan list user
func (usr *userHandler) Find(c *fiber.Ctx) error {
	userList, apiErr := usr.service.FindUsers(c.Context())
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userList})
}

// Edit mengedit user oleh admin
func (usr *userHandler) Edit(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	var user dto.UserEditRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := usr.service.EditUser(c.Context(), userID, user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// UpdateFcmToken mengupdateFCM token
func (usr *userHandler) UpdateFcmToken(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var fcmPayload dto.UserUpdateFcmRequest
	if err := c.BodyParser(&fcmPayload); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := fcmPayload.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := usr.service.EditFcm(c.Context(), claims.Identity, fcmPayload.FcmToken)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// Delete menghapus user, idealnya melalui middleware is_admin
func (usr *userHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	userIDParams := c.Params("user_id")

	if claims.Identity == userIDParams {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := usr.service.DeleteUser(c.Context(), userIDParams)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("user %s berhasil dihapus", userIDParams)})
}

// ChangePassword mengganti password pada user sendiri
func (usr *userHandler) ChangePassword(c *fiber.Ctx) error {
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

	apiErr := usr.service.ChangePassword(c.Context(), user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"error": apiErr, "data": "Password berhasil diubah!"})
}

// ResetPassword mengganti password oleh admin pada user tertentu
func (usr *userHandler) ResetPassword(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	data := dto.UserChangePasswordRequest{
		ID:          userID,
		NewPassword: "Password",
	}

	apiErr := usr.service.ResetPassword(c.Context(), data)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("Password user %s berhasil di reset!", c.Params("user_id"))})
}

// Login login
func (usr *userHandler) Login(c *fiber.Ctx) error {
	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := login.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | validate | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := usr.service.Login(c.Context(), login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// RefreshToken
func (usr *userHandler) RefreshToken(c *fiber.Ctx) error {
	var payload dto.UserRefreshTokenRequest
	if err := c.BodyParser(&payload); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := payload.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := usr.service.Refresh(c.Context(), payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
// sesuai authorisasi aktif. File disimpan di folder static/images dengan nama file == jwt.identity alias username
func (usr *userHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	randomName := fmt.Sprintf("%s%v", claims.Identity, time.Now().Unix())
	pathInDB, apiErr := saveImage(c, *claims, "avatar", randomName, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	usersResult, apiErr := usr.service.PutAvatar(c.Context(), claims.Identity, pathInDB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": usersResult})
}
