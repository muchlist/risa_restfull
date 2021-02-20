package dto

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User struct lengkap dari document user di Mongodb
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Name      string             `json:"name" bson:"name"`
	IsAdmin   bool               `json:"is_admin" bson:"is_admin"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	HashPw    string             `json:"hash_pw,omitempty" bson:"hash_pw,omitempty"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
}

//UserResponseList tipe slice dari UserResponse
type UserResponseList []UserResponse

//UserResponse struct kembalian dari MongoDB dengan menghilangkan hashPassword
type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Name      string             `json:"name" bson:"name"`
	IsAdmin   bool               `json:"is_admin" bson:"is_admin"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
}

//UserRequest input JSON untuk keperluan register, timestamp dapat diabaikan
type UserRequest struct {
	Email     string `json:"email" bson:"email"`
	Name      string `json:"name" bson:"name"`
	IsAdmin   bool   `json:"is_admin" bson:"is_admin"`
	Avatar    string `json:"avatar" bson:"avatar"`
	Password  string `json:"password" bson:"password"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}

//Validate input
func (u UserRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	)
}

//UserEditRequest input JSON oleh admin untuk mengedit user
type UserEditRequest struct {
	Name            string `json:"name" bson:"name"`
	IsAdmin         bool   `json:"is_admin" bson:"is_admin"`
	TimestampFilter int64  `json:"timestamp_filter" bson:"timestamp"`
}

func (u UserEditRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.TimestampFilter, validation.Required),
	)
}

//UserLoginRequest input JSON oleh client untuk keperluan login
type UserLoginRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Limit    int    `json:"limit"`
}

//Validate input
func (u UserLoginRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	)
}

//UserChangePasswordRequest struck untuk keperluan change password dan reset password
//pada reset password hanya menggunakan NewPassword dan mengabaikan Password
type UserChangePasswordRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

//Validate input
func (u UserChangePasswordRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
		validation.Field(&u.NewPassword, validation.Required, validation.Length(3, 20)),
	)
}

//UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginResponse struct {
	Email        string `json:"email" bson:"email"`
	Name         string `json:"name" bson:"name"`
	IsAdmin      bool   `json:"is_admin" bson:"is_admin"`
	Avatar       string `json:"avatar" bson:"avatar"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expired      int64  `json:"expired"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Limit        int    `json:"limit"`
}

//Validate input
func (u UserRefreshTokenRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.RefreshToken, validation.Required),
	)
}

//UserRefreshTokenResponse mengembalikan token dengan claims yang
//sama dengan token sebelumnya dengan expired yang baru
type UserRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expired     int64  `json:"expired"`
}
