package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/crypt"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
	"time"
)

func TestUserService_GetUserByID(t *testing.T) {

	objectID := primitive.NewObjectID()

	m := new(dao.MockDao)
	m.On("GetUserByID", objectID).Return(&dto.UserResponse{
		ID:        objectID,
		Email:     "whois.muchlis@gmail.com",
		Name:      "muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Timestamp: 1610350965,
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	user, err := service.GetUser(objectID)

	assert.Nil(t, err)
	assert.Equal(t, "muchlis", user.Name)
	assert.Equal(t, "whois.muchlis@gmail.com", user.Email)
	assert.Equal(t, true, user.IsAdmin)
}

func TestUserService_GetUser_NoUserFound(t *testing.T) {
	objectID := primitive.NewObjectID()

	m := new(dao.MockDao)
	m.On("GetUserByID", objectID).Return(nil, rest_err.NewNotFoundError(fmt.Sprintf("User dengan ID %v tidak ditemukan", objectID.Hex())))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	user, err := service.GetUser(objectID)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("User dengan ID %v tidak ditemukan", objectID.Hex()), err.Message())
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestUserService_GetUserByEmail_Found(t *testing.T) {

	email := "whois.muchlis@gmail.com"

	m := new(dao.MockDao)
	m.On("GetUserByEmail", email).Return(&dto.UserResponse{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Name:      "Muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Timestamp: time.Now().Unix(),
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	user, err := service.GetUserByEmail(email)

	assert.Nil(t, err)
	assert.Equal(t, "Muchlis", user.Name)
	assert.Equal(t, "whois.muchlis@gmail.com", user.Email)
	assert.Equal(t, true, user.IsAdmin)
}

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {

	email := "whois.muchlis@gmail.com"

	m := new(dao.MockDao)
	m.On("GetUserByEmail", email).Return(nil, rest_err.NewNotFoundError(fmt.Sprintf("User dengan Email %s tidak ditemukan", email)))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	user, err := service.GetUserByEmail(email)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "User dengan Email whois.muchlis@gmail.com tidak ditemukan", err.Message())
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestUserService_FindUsers(t *testing.T) {

	m := new(dao.MockDao)
	m.On("FindUser").Return(dto.UserResponseList{
		dto.UserResponse{
			ID:        primitive.NewObjectID(),
			Email:     "whois.muchlis@gmail.com",
			Name:      "Muchlis",
			IsAdmin:   true,
			Avatar:    "",
			Timestamp: time.Now().Unix(),
		},
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	usersResult, err := service.FindUsers()

	assert.Nil(t, err)
	assert.Equal(t, "Muchlis", usersResult[0].Name)
	assert.Equal(t, "whois.muchlis@gmail.com", usersResult[0].Email)
}

func TestUserService_FindUsers_errorDatabase(t *testing.T) {
	m := new(dao.MockDao)
	m.On("FindUser").Return(dto.UserResponseList(nil), rest_err.NewInternalServerError("Database error", nil))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	usersResult, err := service.FindUsers()

	assert.NotNil(t, err)
	assert.Equal(t, dto.UserResponseList(nil), usersResult)
	assert.Equal(t, "Database error", err.Message())
	assert.Equal(t, http.StatusInternalServerError, err.Status())
}

func TestUserService_InsertUser_Success(t *testing.T) {
	userInput := dto.UserRequest{
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Password:  "password",
		Timestamp: time.Now().Unix(),
	}

	email := "whowho@gmail.com"
	idResp := "5f969f62259eae481fb0e856"

	m := new(dao.MockDao)
	m.On("InsertUser", mock.Anything).Return(&idResp, nil)
	m.On("CheckEmailAvailable", email).Return(true, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	insertedId, err := service.InsertUser(userInput)

	assert.Nil(t, err)
	assert.Equal(t, "5f969f62259eae481fb0e856", *insertedId)
}

func TestUserService_InsertUser_GenerateHashFailed(t *testing.T) {
	userInput := dto.UserRequest{
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Password:  "password",
		Timestamp: time.Now().Unix(),
	}

	email := "whowho@gmail.com"
	idResp := "5f969f62259eae481fb0e856"

	m := new(dao.MockDao)
	m.On("InsertUser", mock.Anything).Return(&idResp, nil)
	m.On("CheckEmailAvailable", email).Return(true, nil)
	c := new(crypt.MockBcrypt)
	c.On("GenerateHash", mock.Anything).Return("", rest_err.NewInternalServerError("Crypto error", nil))

	service := NewUserService(m, c, mjwt.NewJwt())

	insertedId, err := service.InsertUser(userInput)

	assert.Nil(t, insertedId)
	assert.NotNil(t, err)
	assert.Equal(t, 500, err.Status())
}

func TestUserService_InsertUser_EmailNotAvailable(t *testing.T) {
	userInput := dto.UserRequest{
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Password:  "password",
		Timestamp: time.Now().Unix(),
	}

	email := "whowho@gmail.com"
	idResp := "5f969f62259eae481fb0e856"

	m := new(dao.MockDao)
	m.On("InsertUser", mock.Anything).Return(&idResp, nil)
	m.On("CheckEmailAvailable", email).Return(false, rest_err.NewBadRequestError("Email tidak tersedia"))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	insertedId, err := service.InsertUser(userInput)

	assert.Nil(t, insertedId)
	assert.NotNil(t, err)
	assert.Equal(t, "Email tidak tersedia", err.Message())
	assert.Equal(t, 400, err.Status())
}

func TestUserService_InsertUser_DBError(t *testing.T) {
	userInput := dto.UserRequest{
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Password:  "password",
		Timestamp: time.Now().Unix(),
	}

	email := "whowho@gmail.com"

	m := new(dao.MockDao)
	m.On("InsertUser", mock.Anything).Return(nil, rest_err.NewInternalServerError("Gagal menyimpan user ke database", errors.New("db error")))
	m.On("CheckEmailAvailable", email).Return(true, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	insertedId, err := service.InsertUser(userInput)

	assert.Nil(t, insertedId)
	assert.NotNil(t, err)
	assert.Equal(t, "Gagal menyimpan user ke database", err.Message())
	assert.Equal(t, 500, err.Status())
}

func TestUserService_EditUser(t *testing.T) {
	email := "whowho@gmail.com"
	userInput := dto.UserEditRequest{
		Name:            "Muchlis",
		IsAdmin:         false,
		TimestampFilter: 0,
	}

	m := new(dao.MockDao)
	m.On("EditUser", email, userInput).Return(&dto.UserResponse{
		ID:        primitive.ObjectID{},
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "",
		Timestamp: 0,
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	userResponse, err := service.EditUser(email, userInput)

	assert.Nil(t, err)
	assert.Equal(t, "Muchlis", userResponse.Name)
}

func TestUserService_EditUser_TimeStampNotmatch(t *testing.T) {

	email := "whowho@gmail.com"
	userInput := dto.UserEditRequest{
		Name:            "Muchlis",
		IsAdmin:         false,
		TimestampFilter: 0,
	}

	m := new(dao.MockDao)
	m.On("EditUser", email, userInput).Return(nil, rest_err.NewBadRequestError("User tidak diupdate karena ID atau timestamp tidak valid"))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	userResponse, err := service.EditUser(email, userInput)

	assert.Nil(t, userResponse)
	assert.NotNil(t, err)
	assert.Equal(t, "User tidak diupdate karena ID atau timestamp tidak valid", err.Message())
	assert.Equal(t, 400, err.Status())
}

func TestUserService_DeleteUser(t *testing.T) {
	email := "whowho@gmail.com"

	m := new(dao.MockDao)
	m.On("DeleteUser", email).Return(nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	err := service.DeleteUser("whowho@gmail.com")

	assert.Nil(t, err)
}

func TestUserService_DeleteUser_Failed(t *testing.T) {
	email := "whowho@gmail.com"

	m := new(dao.MockDao)
	m.On("DeleteUser", email).Return(rest_err.NewBadRequestError("User gagal dihapus, dokumen tidak ditemukan"))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	err := service.DeleteUser("whowho@gmail.com")

	assert.NotNil(t, err)
	assert.Equal(t, "User gagal dihapus, dokumen tidak ditemukan", err.Message())
}

func TestUserService_Login(t *testing.T) {

	userRequest := dto.UserLoginRequest{
		Email:    "whowho@gmail.com",
		Password: "Password",
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", userRequest.Email).Return(&dto.User{
		ID:        primitive.ObjectID{},
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "",
		HashPw:    "$2a$04$N.8j0ys/1t8YBZuM051PQOq3B6p5hFNv2hzYr.1vooL65z9Bmb7fO",
		Timestamp: 0,
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	userResult, err := service.Login(userRequest)

	assert.Nil(t, err)
	assert.NotNil(t, userRequest)
	assert.Equal(t, "Muchlis", userResult.Name)
	assert.NotEmpty(t, userResult.AccessToken)
	assert.NotEmpty(t, userResult.RefreshToken)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	userRequest := dto.UserLoginRequest{
		Email:    "whowho@gmail.com",
		Password: "salahPassword",
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", userRequest.Email).Return(&dto.User{
		ID:        primitive.ObjectID{},
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "",
		HashPw:    "$2a$04$N.8j0ys/1t8YBZuM051PQOq3B6p5hFNv2hzYr.1vooL65z9Bmb7fO",
		Timestamp: 0,
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	userResult, err := service.Login(userRequest)

	assert.Nil(t, userResult)
	assert.NotNil(t, err)
	assert.Equal(t, "Username atau password tidak valid", err.Message())
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	userRequest := dto.UserLoginRequest{
		Email:    "notexist@gmail.com",
		Password: "salahPassword",
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", userRequest.Email).Return(nil, rest_err.NewUnauthorizedError("Username atau password tidak valid"))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	userResult, err := service.Login(userRequest)

	assert.Nil(t, userResult)
	assert.NotNil(t, err)
	assert.Equal(t, "Username atau password tidak valid", err.Message())
	assert.Equal(t, 401, err.Status())
}

func TestUserService_Login_GenerateTokenError(t *testing.T) {

	userRequest := dto.UserLoginRequest{
		Email:    "whowho@gmail.com",
		Password: "Password",
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", userRequest.Email).Return(&dto.User{
		ID:        primitive.ObjectID{},
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "",
		HashPw:    "$2a$04$N.8j0ys/1t8YBZuM051PQOq3B6p5hFNv2hzYr.1vooL65z9Bmb7fO",
		Timestamp: 0,
	}, nil)
	j := new(mjwt.MockJwt)
	j.On("GenerateToken", mock.Anything).Return("", rest_err.NewInternalServerError("gagal menandatangani token", nil))

	service := NewUserService(m, crypt.NewCrypto(), j)
	userResult, err := service.Login(userRequest)

	assert.Nil(t, userResult)
	assert.NotNil(t, err)
	assert.Equal(t, 500, err.Status())
}

func TestUserService_PutAvatar(t *testing.T) {

	email := "whowhos@gmail.com"
	filePath := "images/whowhos@gmail.com.jpg"

	m := new(dao.MockDao)
	m.On("PutAvatar", email, filePath).Return(&dto.UserResponse{
		ID:        primitive.ObjectID{},
		Email:     "whowhos@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "images/whowhos@gmail.com.jpg",
		Timestamp: 0,
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	userResult, err := service.PutAvatar("whowhos@gmail.com", "images/whowhos@gmail.com.jpg")

	assert.Nil(t, err)
	assert.Equal(t, "images/whowhos@gmail.com.jpg", userResult.Avatar)
}

func TestUserService_PutAvatar_UserNotFound(t *testing.T) {

	email := "whowhos@gmail.com"
	filePath := "images/whowhos@gmail.com.jpg"

	m := new(dao.MockDao)
	m.On("PutAvatar", email, filePath).Return(nil, rest_err.NewBadRequestError(fmt.Sprintf("User avatar gagal diupload, user dengan email %s tidak ditemukan", email)))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	userResult, err := service.PutAvatar("whowhos@gmail.com", "images/whowhos@gmail.com.jpg")

	assert.Nil(t, userResult)
	assert.NotNil(t, err)
	assert.Equal(t, "User avatar gagal diupload, user dengan email whowhos@gmail.com tidak ditemukan", err.Message())
}

func TestUserService_ChangePassword_Success(t *testing.T) {
	data := dto.UserChangePasswordRequest{
		Email:       "whowho@gmail.com",
		Password:    "Password",
		NewPassword: "NewPassword",
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", mock.Anything).Return(&dto.User{
		ID:        primitive.ObjectID{},
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "",
		HashPw:    "$2a$04$N.8j0ys/1t8YBZuM051PQOq3B6p5hFNv2hzYr.1vooL65z9Bmb7fO",
		Timestamp: 0,
	}, nil)
	m.On("ChangePassword", mock.Anything).Return(nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	err := service.ChangePassword(data)

	assert.Nil(t, err)
}

func TestUserService_ChangePassword_HashNewPasswordErr(t *testing.T) {
	data := dto.UserChangePasswordRequest{
		Email:       "whowho@gmail.com",
		Password:    "Password",
		NewPassword: "NewPassword",
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", mock.Anything).Return(&dto.User{
		ID:        primitive.ObjectID{},
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "",
		HashPw:    "$2a$04$N.8j0ys/1t8YBZuM051PQOq3B6p5hFNv2hzYr.1vooL65z9Bmb7fO",
		Timestamp: 0,
	}, nil)
	m.On("ChangePassword", mock.Anything).Return(nil)

	c := new(crypt.MockBcrypt)
	c.On("GenerateHash", mock.Anything).Return("", rest_err.NewInternalServerError("Crypto error", nil))
	c.On("IsPWAndHashPWMatch", mock.Anything, mock.Anything).Return(true)

	service := NewUserService(m, c, mjwt.NewJwt())
	err := service.ChangePassword(data)

	assert.NotNil(t, err)
	assert.Equal(t, 500, err.Status())
	assert.Equal(t, "Crypto error", err.Message())
}

func TestUserService_ChangePassword_FailPasswordSame(t *testing.T) {

	data := dto.UserChangePasswordRequest{
		Email:       "whowho@gmail.com",
		Password:    "Password",
		NewPassword: "Password",
	}
	m := new(dao.MockDao)
	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	err := service.ChangePassword(data)

	assert.NotNil(t, err)
	assert.Equal(t, "Gagal mengganti password, password tidak boleh sama dengan sebelumnya!", err.Message())
}

func TestUserService_ChangePassword_OldPasswordWrong(t *testing.T) {

	data := dto.UserChangePasswordRequest{
		Email:       "whowho@gmail.com",
		Password:    "salahPassword",
		NewPassword: "NewPassword",
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", mock.Anything).Return(&dto.User{
		ID:        primitive.ObjectID{},
		Email:     "whowho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   false,
		Avatar:    "",
		HashPw:    "$2a$04$N.8j0ys/1t8YBZuM051PQOq3B6p5hFNv2hzYr.1vooL65z9Bmb7fO",
		Timestamp: 0,
	}, nil)
	m.On("ChangePassword", mock.Anything).Return(nil)

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	err := service.ChangePassword(data)

	assert.NotNil(t, err)
	assert.Equal(t, "Gagal mengganti password, password salah!", err.Message())
}

func TestUserService_ChangePassword_EmailWrong(t *testing.T) {
	data := dto.UserChangePasswordRequest{
		Email:       "email.salah@gmail.com",
		Password:    "password",
		NewPassword: "Password",
	}
	m := new(dao.MockDao)
	m.On("GetUserByEmailWithPassword", mock.Anything).Return(nil, rest_err.NewUnauthorizedError("Username atau password tidak valid"))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	err := service.ChangePassword(data)

	assert.NotNil(t, err)
	assert.Equal(t, "Username atau password tidak valid", err.Message())
}

func TestUserService_ResetPassword(t *testing.T) {
	data := dto.UserChangePasswordRequest{
		Email:       "whowho@gmail.com",
		Password:    "",
		NewPassword: "PasswordBaru",
	}

	m := new(dao.MockDao)
	m.On("ChangePassword", mock.Anything).Return(nil)
	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())

	err := service.ResetPassword(data)

	assert.Nil(t, err)
}

func TestUserService_ResetPassword_EmailNotFound(t *testing.T) {
	data := dto.UserChangePasswordRequest{
		Email:       "emailsalah@gmail.com",
		Password:    "",
		NewPassword: "PasswordBaru",
	}
	m := new(dao.MockDao)
	m.On("ChangePassword", mock.Anything).Return(rest_err.NewBadRequestError("Penggantian password gagal, email salah"))

	service := NewUserService(m, crypt.NewCrypto(), mjwt.NewJwt())
	err := service.ResetPassword(data)

	assert.NotNil(t, err)
	assert.Equal(t, "Penggantian password gagal, email salah", err.Message())
}

func TestUserService_ResetPassword_GenerateHashFailed(t *testing.T) {
	data := dto.UserChangePasswordRequest{
		Email:       "whowho@gmail.com",
		Password:    "",
		NewPassword: "PasswordBaru",
	}

	m := new(dao.MockDao)
	m.On("ChangePassword", mock.Anything).Return(nil)
	c := new(crypt.MockBcrypt)
	c.On("GenerateHash", data.NewPassword).Return("", rest_err.NewInternalServerError("Crypto error", nil))

	service := NewUserService(m, c, mjwt.NewJwt())

	err := service.ResetPassword(data)

	assert.NotNil(t, err)
	assert.Equal(t, 500, err.Status())
}

func TestUserService_Refresh_Success(t *testing.T) {
	input := dto.UserRefreshTokenRequest{
		Limit: 5,
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmail", mock.Anything).Return(&dto.UserResponse{
		ID:        primitive.NewObjectID(),
		Email:     "whoswho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Timestamp: time.Now().Unix(),
	}, nil)
	j := new(mjwt.MockJwt)
	j.On("ValidateToken", mock.Anything).Return(&jwt.Token{}, nil)
	j.On("ReadToken", mock.Anything).Return(&mjwt.CustomClaim{
		Type: mjwt.Refresh,
	}, nil)

	j.On("GenerateToken", mock.Anything).Return("accessToken", nil)

	service := NewUserService(m, crypt.NewCrypto(), j)
	res, err := service.Refresh(input)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "accessToken", res.AccessToken)
	assert.Equal(t, time.Now().Add(time.Minute*time.Duration(input.Limit)).Unix(), res.Expired)

}

func TestUserService_Refresh_User_Not_Found(t *testing.T) {
	input := dto.UserRefreshTokenRequest{
		Limit: 5,
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmail", mock.Anything).Return(nil, rest_err.NewNotFoundError(fmt.Sprint("User dengan Email whoswho@gmail.com tidak ditemukan")))
	j := new(mjwt.MockJwt)
	j.On("ValidateToken", mock.Anything).Return(&jwt.Token{}, nil)
	j.On("ReadToken", mock.Anything).Return(&mjwt.CustomClaim{
		Type: mjwt.Refresh,
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), j)
	res, err := service.Refresh(input)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "User dengan Email whoswho@gmail.com tidak ditemukan", err.Message())

}

func TestUserService_Refresh_Token_Not_Valid(t *testing.T) {
	input := dto.UserRefreshTokenRequest{
		Limit: 5,
	}

	m := new(dao.MockDao)
	j := new(mjwt.MockJwt)
	j.On("ValidateToken", mock.Anything).Return(nil, rest_err.NewAPIError("Token signing method salah", http.StatusUnprocessableEntity, "jwt_error", nil))
	service := NewUserService(m, crypt.NewCrypto(), j)
	res, err := service.Refresh(input)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, err.Status())

}

func TestUserService_Refresh_Token_Read_Error(t *testing.T) {
	input := dto.UserRefreshTokenRequest{
		Limit: 5,
	}

	m := new(dao.MockDao)
	j := new(mjwt.MockJwt)
	j.On("ValidateToken", mock.Anything).Return(&jwt.Token{}, nil)
	j.On("ReadToken", mock.Anything).Return(nil, rest_err.NewInternalServerError("gagal mapping token", nil))

	service := NewUserService(m, crypt.NewCrypto(), j)
	res, err := service.Refresh(input)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}

func TestUserService_Refresh_Token_Not_Refresh_Token(t *testing.T) {
	input := dto.UserRefreshTokenRequest{
		Limit: 1000000,
	}

	m := new(dao.MockDao)
	j := new(mjwt.MockJwt)
	j.On("ValidateToken", mock.Anything).Return(&jwt.Token{}, nil)
	j.On("ReadToken", mock.Anything).Return(&mjwt.CustomClaim{
		Type: mjwt.Access,
	}, nil)

	service := NewUserService(m, crypt.NewCrypto(), j)
	res, err := service.Refresh(input)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, err.Status())

}

func TestUserService_Refresh_Token_Generate_Token_Error(t *testing.T) {
	input := dto.UserRefreshTokenRequest{
		Limit: 1000000,
	}

	m := new(dao.MockDao)
	m.On("GetUserByEmail", mock.Anything).Return(&dto.UserResponse{
		ID:        primitive.NewObjectID(),
		Email:     "whoswho@gmail.com",
		Name:      "Muchlis",
		IsAdmin:   true,
		Avatar:    "",
		Timestamp: time.Now().Unix(),
	}, nil)
	j := new(mjwt.MockJwt)
	j.On("ValidateToken", mock.Anything).Return(&jwt.Token{}, nil)
	j.On("ReadToken", mock.Anything).Return(&mjwt.CustomClaim{
		Type: mjwt.Refresh,
	}, nil)
	j.On("GenerateToken", mock.Anything).Return("", rest_err.NewInternalServerError("gagal menandatangani token", nil))

	service := NewUserService(m, crypt.NewCrypto(), j)
	res, err := service.Refresh(input)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status())

}
