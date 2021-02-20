package crypt

import (
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"golang.org/x/crypto/bcrypt"
)

func NewCrypto() BcryptAssumer {
	return &cryptoObj{}
}

type BcryptAssumer interface {
	GenerateHash(password string) (string, rest_err.APIError)
	IsPWAndHashPWMatch(password string, hashPass string) bool
}

type cryptoObj struct {
}

//GenerateHash membuat hashpassword, hash password 1 dengan yang lainnya akan berbeda meskipun
//inputannya sama, sehingga untuk membandingkan hashpassword memerlukan method lain IsPWAndHashPWMatch
func (c *cryptoObj) GenerateHash(password string) (string, rest_err.APIError) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.Error("Error pada kriptograpi (GenerateHash)", err)
		restErr := rest_err.NewInternalServerError("Crypto error", err)
		return "", restErr
	}
	return string(passwordHash), nil
}

//IsPWAndHashPWMatch return true jika inputan password dan hashpassword sesuai
func (c *cryptoObj) IsPWAndHashPWMatch(password string, hashPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	return err == nil
}
