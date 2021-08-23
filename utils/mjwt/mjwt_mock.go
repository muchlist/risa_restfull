package mjwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/stretchr/testify/mock"
)

type MockJwt struct {
	mock.Mock
}

func (m *MockJwt) GenerateToken(claims CustomClaim) (string, rest_err.APIError) {
	args := m.Called(claims)
	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return args.Get(0).(string), err
}

func (m *MockJwt) ValidateToken(tokenString string) (*jwt.Token, rest_err.APIError) {
	args := m.Called(tokenString)
	var res *jwt.Token
	if args.Get(0) != nil {
		res = args.Get(0).(*jwt.Token)
	}

	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return res, err
}

func (m *MockJwt) ReadToken(token *jwt.Token) (*CustomClaim, rest_err.APIError) {
	args := m.Called(token)
	var res *CustomClaim
	if args.Get(0) != nil {
		res = args.Get(0).(*CustomClaim)
	}

	var err rest_err.APIError
	if args.Get(1) != nil {
		err = args.Get(1).(rest_err.APIError)
	}

	return res, err
}
