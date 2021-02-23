package user_dao

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
)

type UserDaoAssumer interface {
	InsertUser(user dto.UserRequest) (*string, rest_err.APIError)
	EditUser(userID string, userRequest dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError)
	DeleteUser(userID string) rest_err.APIError
	PutAvatar(userID string, avatar string) (*dto.UserResponse, rest_err.APIError)
	ChangePassword(data dto.UserChangePasswordRequest) rest_err.APIError

	GetUserByID(userID string) (*dto.UserResponse, rest_err.APIError)
	GetUserByIDWithPassword(userID string) (*dto.User, rest_err.APIError)
	FindUser() (dto.UserResponseList, rest_err.APIError)
	CheckIDAvailable(email string) (bool, rest_err.APIError)
}
