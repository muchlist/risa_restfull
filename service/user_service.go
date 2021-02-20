package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/crypt"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

func NewUserService(dao dao.UserDaoAssumer, crypto crypt.BcryptAssumer, jwt mjwt.JWTAssumer) UserServiceAssumer {
	return &userService{
		dao:    dao,
		crypto: crypto,
		jwt:    jwt,
	}
}

type userService struct {
	dao    dao.UserDaoAssumer
	crypto crypt.BcryptAssumer
	jwt    mjwt.JWTAssumer
}

type UserServiceAssumer interface {
	GetUser(primitive.ObjectID) (*dto.UserResponse, rest_err.APIError)
	GetUserByEmail(email string) (*dto.UserResponse, rest_err.APIError)
	InsertUser(dto.UserRequest) (*string, rest_err.APIError)
	FindUsers() (dto.UserResponseList, rest_err.APIError)
	EditUser(email string, userEdit dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError)
	DeleteUser(email string) rest_err.APIError
	Login(dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError)
	Refresh(login dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
	PutAvatar(email string, fileLocation string) (*dto.UserResponse, rest_err.APIError)
	ChangePassword(data dto.UserChangePasswordRequest) rest_err.APIError
	ResetPassword(data dto.UserChangePasswordRequest) rest_err.APIError
}

//GetUser mendapatkan user dari database
func (u *userService) GetUser(userID primitive.ObjectID) (*dto.UserResponse, rest_err.APIError) {
	user, err := u.dao.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetUserByEmail mendapatkan user berdasarkan email
func (u *userService) GetUserByEmail(email string) (*dto.UserResponse, rest_err.APIError) {
	user, err := u.dao.GetUserByEmail(strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	return user, nil
}

//FindUsers
func (u *userService) FindUsers() (dto.UserResponseList, rest_err.APIError) {
	userList, err := u.dao.FindUser()
	if err != nil {
		return nil, err
	}
	return userList, nil
}

//InsertUser melakukan register user semua email yang di registrasikan diubah menjadi lowercase di tahap ini
func (u *userService) InsertUser(user dto.UserRequest) (*string, rest_err.APIError) {

	user.Email = strings.ToLower(user.Email)
	user.IsAdmin = false

	// cek ketersediaan email
	_, err := u.dao.CheckEmailAvailable(user.Email)
	if err != nil {
		return nil, err
	}
	// END cek ketersediaan email

	hashPassword, err := u.crypto.GenerateHash(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPassword
	user.Timestamp = time.Now().Unix()

	insertedID, err := u.dao.InsertUser(user)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
}

//EditUser
func (u *userService) EditUser(email string, request dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError) {
	result, err := u.dao.EditUser(strings.ToLower(email), request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//DeleteUser
func (u *userService) DeleteUser(email string) rest_err.APIError {
	err := u.dao.DeleteUser(email)
	if err != nil {
		return err
	}

	return nil
}

//Login
func (u *userService) Login(login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError) {

	login.Email = strings.ToLower(login.Email)

	user, err := u.dao.GetUserByEmailWithPassword(login.Email)
	if err != nil {
		return nil, err
	}

	if !u.crypto.IsPWAndHashPWMatch(login.Password, user.HashPw) {
		return nil, rest_err.NewUnauthorizedError("Username atau password tidak valid")
	}

	if login.Limit == 0 || login.Limit > 10080 { // 10080 minute = 7 day
		login.Limit = 10080
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.Email,
		Name:        user.Name,
		IsAdmin:     user.IsAdmin,
		ExtraMinute: time.Duration(login.Limit),
		Type:        mjwt.Access,
		Fresh:       true,
	}

	RefreshClaims := mjwt.CustomClaim{
		Identity:    user.Email,
		Name:        user.Name,
		IsAdmin:     user.IsAdmin,
		ExtraMinute: 10080, // 7 days
		Type:        mjwt.Refresh,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	refreshToken, err := u.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserLoginResponse{
		Name:         user.Name,
		Email:        user.Email,
		IsAdmin:      user.IsAdmin,
		Avatar:       user.Avatar,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expired:      time.Now().Add(time.Minute * time.Duration(login.Limit)).Unix(),
	}

	return &userResponse, nil

}

//Refresh token
func (u *userService) Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError) {

	token, apiErr := u.jwt.ValidateToken(payload.RefreshToken)
	if apiErr != nil {
		return nil, apiErr
	}
	claims, apiErr := u.jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}

	// cek apakah tipe claims token yang dikirim adalah tipe refresh (1)
	if claims.Type != mjwt.Refresh {
		return nil, rest_err.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{"not a refresh token"})
	}

	// mendapatkan data terbaru dari user
	user, apiErr := u.dao.GetUserByEmail(strings.ToLower(claims.Identity))
	if apiErr != nil {
		return nil, apiErr
	}

	if payload.Limit == 0 || payload.Limit > 10080 { // 10080 minute = 7 day
		payload.Limit = 10080
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.Email,
		Name:        user.Name,
		IsAdmin:     user.IsAdmin,
		ExtraMinute: time.Duration(payload.Limit),
		Type:        mjwt.Access,
		Fresh:       false,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}

	userRefreshTokenResponse := dto.UserRefreshTokenResponse{
		AccessToken: accessToken,
		Expired:     time.Now().Add(time.Minute * time.Duration(payload.Limit)).Unix(),
	}

	return &userRefreshTokenResponse, nil
}

//PutAvatar memasukkan lokasi file (path) ke dalam database user
func (u *userService) PutAvatar(email string, fileLocation string) (*dto.UserResponse, rest_err.APIError) {

	email = strings.ToLower(email)

	user, err := u.dao.PutAvatar(email, fileLocation)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//ChangePassword melakukan perbandingan hashpassword lama dan memasukkan hashpassword baru ke database
func (u *userService) ChangePassword(data dto.UserChangePasswordRequest) rest_err.APIError {

	if data.Password == data.NewPassword {
		return rest_err.NewBadRequestError("Gagal mengganti password, password tidak boleh sama dengan sebelumnya!")
	}

	userResult, err := u.dao.GetUserByEmailWithPassword(data.Email)
	if err != nil {
		return err
	}

	if !u.crypto.IsPWAndHashPWMatch(data.Password, userResult.HashPw) {
		return rest_err.NewBadRequestError("Gagal mengganti password, password salah!")
	}

	newPasswordHash, err := u.crypto.GenerateHash(data.NewPassword)
	if err != nil {
		return err
	}
	data.NewPassword = newPasswordHash

	_ = u.dao.ChangePassword(data)

	return nil
}

//ResetPassword . inputan password berada di level handler
//hanya memproses field newPassword, mengabaikan field password
func (u *userService) ResetPassword(data dto.UserChangePasswordRequest) rest_err.APIError {

	data.Email = strings.ToLower(data.Email)

	newPasswordHash, err := u.crypto.GenerateHash(data.NewPassword)
	if err != nil {
		return err
	}
	data.NewPassword = newPasswordHash

	err = u.dao.ChangePassword(data)
	if err != nil {
		return err
	}

	return nil
}
