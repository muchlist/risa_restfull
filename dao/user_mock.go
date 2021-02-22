package dao

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/stretchr/testify/mock"
)

type MockDao struct {
	mock.Mock
}

func (m *MockDao) InsertUser(user dto.UserRequest) (*string, rest_err.APIError) {
	args := m.Called(user)

	var res *string = nil
	if args.Get(0) != nil {
		res = args.Get(0).(*string)
	}

	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return res, err
}

func (m *MockDao) GetUserByIDWithPassword(email string) (*dto.User, rest_err.APIError) {
	args := m.Called(email)
	var res *dto.User = nil
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.User)
	}
	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return res, err
}

func (m *MockDao) CheckIDAvailable(email string) (bool, rest_err.APIError) {
	args := m.Called(email)
	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return args.Get(0).(bool), err
}

func (m *MockDao) EditUser(userEmail string, userRequest dto.UserEditRequest) (*dto.UserResponse, rest_err.APIError) {
	args := m.Called(userEmail, userRequest)
	var res *dto.UserResponse = nil
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.UserResponse)
	}
	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return res, err
}

func (m *MockDao) DeleteUser(userEmail string) rest_err.APIError {
	args := m.Called(userEmail)
	var err rest_err.APIError = nil
	if args.Get(0) != nil {
		err = args.Get(0).(rest_err.APIError)
	}
	return err
}

func (m *MockDao) PutAvatar(email string, avatar string) (*dto.UserResponse, rest_err.APIError) {
	args := m.Called(email, avatar)
	var res *dto.UserResponse = nil
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.UserResponse)
	}
	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}
	return res, err
}

func (m *MockDao) ChangePassword(data dto.UserChangePasswordRequest) rest_err.APIError {
	args := m.Called(data)
	var err rest_err.APIError = nil
	if args.Get(0) != nil {
		err = args.Get(0).(rest_err.APIError)
	}
	return err
}

func (m *MockDao) GetUserByID(userID string) (*dto.UserResponse, rest_err.APIError) {
	args := m.Called(userID)

	var res *dto.UserResponse = nil
	if args.Get(0) != nil {
		res = args.Get(0).(*dto.UserResponse)
	}

	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return res, err
}

func (m *MockDao) FindUser() (dto.UserResponseList, rest_err.APIError) {
	args := m.Called()

	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return args.Get(0).(dto.UserResponseList), err
}
