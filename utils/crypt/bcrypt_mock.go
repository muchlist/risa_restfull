package crypt

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/stretchr/testify/mock"
)

type MockBcrypt struct {
	mock.Mock
}

func (m *MockBcrypt) GenerateHash(password string) (string, rest_err.APIError) {
	args := m.Called(password)
	var err rest_err.APIError = nil
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return args.Get(0).(string), err
}

func (m *MockBcrypt) IsPWAndHashPWMatch(password string, hashPass string) bool {
	args := m.Called(password, hashPass)
	return args.Bool(0)
}
