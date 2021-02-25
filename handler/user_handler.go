package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"net/http"
	"path/filepath"
)

func NewUserHandler(userService service.UserServiceAssumer) *userHandler {
	return &userHandler{
		service: userService,
	}
}

type userHandler struct {
	service service.UserServiceAssumer
}

//Get menampilkan user berdasarkan ID (bukan email)
func (u *userHandler) Get(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	//userID, err := primitive.ObjectIDFromHex(c.Params("user_id"))
	//if err != nil {
	//	apiErr := rest_err.NewBadRequestError("Format userID salah")
	//	return c.Status(apiErr.Status()).JSON(apiErr)
	//}

	user, apiErr := u.service.GetUser(userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(user)
}

//GetProfile mengembalikan user yang sedang login
func (u *userHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	user, apiErr := u.service.GetUserByID(claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(user)
}

//Register menambahkan user
func (u *userHandler) Register(c *fiber.Ctx) error {

	var user dto.UserRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	insertID, apiErr := u.service.InsertUser(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Register berhasil, ID: %s", *insertID)}
	return c.JSON(res)
}

//Find menampilkan list user
func (u *userHandler) Find(c *fiber.Ctx) error {

	userList, apiErr := u.service.FindUsers()
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"users": userList})
}

//Edit mengedit user oleh admin
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
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	userEdited, apiErr := u.service.EditUser(userID, user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(userEdited)
}

//Delete menghapus user, idealnya melalui middleware is_admin
func (u *userHandler) Delete(c *fiber.Ctx) error {

	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	userIDParams := c.Params("user_id")

	if claims.Identity == userIDParams {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	apiErr := u.service.DeleteUser(userIDParams)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("user %s berhasil dihapus", userIDParams)})
}

//ChangePassword mengganti password pada user sendiri
func (u *userHandler) ChangePassword(c *fiber.Ctx) error {

	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var user dto.UserChangePasswordRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	//mengganti user id dengan user aktif
	user.ID = claims.Identity

	apiErr := u.service.ChangePassword(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": "Password berhasil diubah!"})
}

//ResetPassword mengganti password oleh admin pada user tertentu
func (u *userHandler) ResetPassword(c *fiber.Ctx) error {

	userID := c.Params("user_id")

	data := dto.UserChangePasswordRequest{
		ID:          userID,
		NewPassword: "Password",
	}

	apiErr := u.service.ResetPassword(data)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("Password user %s berhasil di reset!", c.Params("user_id"))})
}

//Login login
func (u *userHandler) Login(c *fiber.Ctx) error {

	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := login.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	response, apiErr := u.service.Login(login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(response)
}

//RefreshToken
func (u *userHandler) RefreshToken(c *fiber.Ctx) error {

	var payload dto.UserRefreshTokenRequest
	if err := c.BodyParser(&payload); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := payload.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	response, apiErr := u.service.Refresh(payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(response)
}

//UploadImage melakukan pengambilan file menggunakan form "avatar" mengecek ekstensi dan memasukkannya ke database
//sesuai authorisasi aktif. File disimpan di folder static/images dengan nama file == jwt.identity alias email
func (u *userHandler) UploadImage(c *fiber.Ctx) error {

	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	file, err := c.FormFile("image")
	if err != nil {
		apiErr := rest_err.NewAPIError("File gagal di upload", http.StatusBadRequest, "bad_request", []interface{}{err.Error()})
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	fileName := file.Filename
	fileExtension := filepath.Ext(fileName)
	if !(fileExtension == ".jpg" || fileExtension == ".png" || fileExtension == ".jpeg") {
		apiErr := rest_err.NewBadRequestError("Ektensi file tidak di support")
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if file.Size > 1*1024*1024 { // 1 MB
		apiErr := rest_err.NewBadRequestError("Ukuran file tidak dapat melebihi 1MB")
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	// rename image to claims.Identity
	path := "static/image/avatar/" + claims.Identity + fileExtension
	pathInDb := "image/avatar/" + claims.Identity + fileExtension

	err = c.SaveFile(file, path)
	if err != nil {
		logger.Error(fmt.Sprintf("%s gagal mengupload file", claims.Identity), err)
		apiErr := rest_err.NewInternalServerError("File gagal di upload", err)
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	usersResult, apiErr := u.service.PutAvatar(claims.Identity, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(usersResult)
}
