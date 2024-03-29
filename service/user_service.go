package service

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/userdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/crypt"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"net/http"
	"time"
)

func NewUserService(dao userdao.UserDaoAssumer, crypto crypt.BcryptAssumer, jwt mjwt.JWTAssumer) UserServiceAssumer {
	return &userService{
		dao:    dao,
		crypto: crypto,
		jwt:    jwt,
	}
}

type userService struct {
	dao    userdao.UserDaoAssumer
	crypto crypt.BcryptAssumer
	jwt    mjwt.JWTAssumer
}

type UserServiceAssumer interface {
	GetUser(ctx context.Context, userID string) (*dto.UserResponse, rest_err.APIError)
	GetUserByID(ctx context.Context, email string) (*dto.UserResponse, rest_err.APIError)
	InsertUser(ctx context.Context, user dto.UserRequest) (*string, rest_err.APIError)
	FindUsers(ctx context.Context) (dto.UserResponseList, rest_err.APIError)
	EditUser(ctx context.Context, userID string, userEdit dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError)
	EditFcm(ctx context.Context, userID string, fcmToken string) (*dto.UserResponse, rest_err.APIError)
	DeleteUser(ctx context.Context, userID string) rest_err.APIError
	Login(ctx context.Context, login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError)
	Refresh(ctx context.Context, login dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
	PutAvatar(ctx context.Context, userID string, fileLocation string) (*dto.UserResponse, rest_err.APIError)
	ChangePassword(ctx context.Context, data dto.UserChangePasswordRequest) rest_err.APIError
	ResetPassword(ctx context.Context, data dto.UserChangePasswordRequest) rest_err.APIError
}

// GetUser mendapatkan user dari database
func (u *userService) GetUser(ctx context.Context, userID string) (*dto.UserResponse, rest_err.APIError) {
	user, err := u.dao.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail mendapatkan user berdasarkan email
func (u *userService) GetUserByID(ctx context.Context, userID string) (*dto.UserResponse, rest_err.APIError) {
	user, err := u.dao.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUsers
func (u *userService) FindUsers(ctx context.Context) (dto.UserResponseList, rest_err.APIError) {
	userList, err := u.dao.FindUser(ctx, "")
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// InsertUser melakukan register user
func (u *userService) InsertUser(ctx context.Context, user dto.UserRequest) (*string, rest_err.APIError) {
	// cek ketersediaan id
	_, err := u.dao.CheckIDAvailable(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	// END cek ketersediaan id

	hashPassword, err := u.crypto.GenerateHash(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPassword
	user.Timestamp = time.Now().Unix()

	insertedID, err := u.dao.InsertUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
}

// EditUser
func (u *userService) EditUser(ctx context.Context, userID string, request dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError) {
	result, err := u.dao.EditUser(ctx, userID, request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userService) EditFcm(ctx context.Context, userID string, fcmToken string) (*dto.UserResponse, rest_err.APIError) {
	result, err := u.dao.EditFcm(ctx, userID, fcmToken)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteUser
func (u *userService) DeleteUser(ctx context.Context, userID string) rest_err.APIError {
	err := u.dao.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// Login
func (u *userService) Login(ctx context.Context, login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError) {
	user, err := u.dao.GetUserByIDWithPassword(ctx, login.ID)
	if err != nil {
		return nil, err
	}

	if !u.crypto.IsPWAndHashPWMatch(login.Password, user.HashPw) {
		return nil, rest_err.NewUnauthorizedError("Username atau password tidak valid")
	}

	if login.Limit == 0 || login.Limit > 60*24*30 { // 30 days
		login.Limit = 60 * 24 * 30
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        user.Name,
		Roles:       user.Roles,
		Branch:      user.Branch,
		ExtraMinute: time.Duration(login.Limit),
		Type:        mjwt.Access,
		Fresh:       true,
	}

	RefreshClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        user.Name,
		Roles:       user.Roles,
		Branch:      user.Branch,
		ExtraMinute: 60 * 24 * 90, // 90 days
		Type:        mjwt.Refresh,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserLoginResponse{
		ID:           user.ID,
		Name:         user.Name,
		Branch:       user.Branch,
		Email:        user.Email,
		Roles:        user.Roles,
		Avatar:       user.Avatar,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expired:      time.Now().Add(time.Minute * time.Duration(login.Limit)).Unix(),
	}

	return &userResponse, nil
}

// Refresh token
func (u *userService) Refresh(ctx context.Context, payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError) {
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
	user, apiErr := u.dao.GetUserByID(ctx, claims.Identity)
	if apiErr != nil {
		return nil, apiErr
	}

	if payload.Limit == 0 || payload.Limit > 60*24*30 { // 30 day
		payload.Limit = 60 * 24 * 30
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        user.Name,
		Roles:       user.Roles,
		Branch:      user.Branch,
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

// PutAvatar memasukkan lokasi file (path) ke dalam database user
func (u *userService) PutAvatar(ctx context.Context, userID string, fileLocation string) (*dto.UserResponse, rest_err.APIError) {
	user, err := u.dao.PutAvatar(ctx, userID, fileLocation)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword melakukan perbandingan hashpassword lama dan memasukkan hashpassword baru ke database
func (u *userService) ChangePassword(ctx context.Context, data dto.UserChangePasswordRequest) rest_err.APIError {
	if data.Password == data.NewPassword {
		return rest_err.NewBadRequestError("Gagal mengganti password, password tidak boleh sama dengan sebelumnya!")
	}

	userResult, err := u.dao.GetUserByIDWithPassword(ctx, data.ID)
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

	_ = u.dao.ChangePassword(ctx, data)

	return nil
}

// ResetPassword . inputan password berada di level handler
// hanya memproses field newPassword, mengabaikan field password
func (u *userService) ResetPassword(ctx context.Context, data dto.UserChangePasswordRequest) rest_err.APIError {
	newPasswordHash, err := u.crypto.GenerateHash(data.NewPassword)
	if err != nil {
		return err
	}
	data.NewPassword = newPasswordHash

	err = u.dao.ChangePassword(ctx, data)
	if err != nil {
		return err
	}

	return nil
}
