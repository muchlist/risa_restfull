package userdao

import (
	"context"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/stretchr/testify/mock"
)

type MockDao struct {
	mock.Mock
}

func (m *MockDao) InsertUser(_ context.Context, user dto.UserRequest) (*string, rest_err.APIError) {
	args := m.Called(user)

	var res *string
	if args.Get(0) != nil {
		res = args.Get(0).(*string)
	}

	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return res, err
}

func (m *MockDao) GetUserByIDWithPassword(_ context.Context, email string) (*dto.User, rest_err.APIError) {
	args := m.Called(email)
	var res *dto.User
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.User)
	}
	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return res, err
}

func (m *MockDao) CheckIDAvailable(_ context.Context, email string) (bool, rest_err.APIError) {
	args := m.Called(email)
	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return args.Get(0).(bool), err
}

func (m *MockDao) EditUser(_ context.Context, userEmail string, userRequest dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError) {
	args := m.Called(userEmail, userRequest)
	var res *dto.UserResponse
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.UserResponse)
	}
	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return res, err
}

func (m *MockDao) EditFcm(_ context.Context, userID string, fcmToken string) (*dto.UserResponse, rest_err.APIError) {
	args := m.Called(userID, fcmToken)
	var res *dto.UserResponse
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.UserResponse)
	}
	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return res, err
}

func (m *MockDao) DeleteUser(_ context.Context, userEmail string) rest_err.APIError {
	args := m.Called(userEmail)
	var err rest_err.APIError
	if args.Get(0) != nil {
		err = args.Get(0).(rest_err.APIError)
	}
	return err
}

func (m *MockDao) PutAvatar(_ context.Context, email string, avatar string) (*dto.UserResponse, rest_err.APIError) {
	args := m.Called(email, avatar)
	var res *dto.UserResponse
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.UserResponse)
	}
	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return res, err
}

func (m *MockDao) ChangePassword(_ context.Context, data dto.UserChangePasswordRequest) rest_err.APIError {
	args := m.Called(data)
	var err rest_err.APIError
	if args.Get(0) != nil {
		err = args.Get(0).(rest_err.APIError)
	}
	return err
}

func (m *MockDao) GetUserByID(_ context.Context, userID string) (*dto.UserResponse, rest_err.APIError) {
	args := m.Called(userID)

	var res *dto.UserResponse
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.UserResponse)
	}

	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return res, err
}

func (m *MockDao) FindUser(_ context.Context, branch string) (dto.UserResponseList, rest_err.APIError) {
	args := m.Called(branch)

	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return args.Get(0).(dto.UserResponseList), err
}
