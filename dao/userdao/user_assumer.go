package userdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
)

type UserDaoAssumer interface {
	UserSaver
	UserLoader
}

type UserSaver interface {
	InsertUser(ctx context.Context, user dto.UserRequest) (*string, rest_err.APIError)
	EditUser(ctx context.Context, userID string, userRequest dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError)
	EditFcm(ctx context.Context, userID string, fcmToken string) (*dto.UserResponse, rest_err.APIError)
	DeleteUser(ctx context.Context, userID string) rest_err.APIError
	PutAvatar(ctx context.Context, userID string, avatar string) (*dto.UserResponse, rest_err.APIError)
	ChangePassword(ctx context.Context, data dto.UserChangePasswordRequest) rest_err.APIError
}

type UserLoader interface {
	GetUserByID(ctx context.Context, userID string) (*dto.UserResponse, rest_err.APIError)
	GetUserByIDWithPassword(ctx context.Context, userID string) (*dto.User, rest_err.APIError)
	FindUser(ctx context.Context, branch string) (dto.UserResponseList, rest_err.APIError)
	CheckIDAvailable(ctx context.Context, email string) (bool, rest_err.APIError)
}
