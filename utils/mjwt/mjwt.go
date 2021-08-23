package mjwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/spf13/cast"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	CLAIMS    = "claims"
	secretKey = "SECRET_KEY"

	identityKey  = "identity"
	nameKey      = "name"
	rolesKey     = "roles"
	branchKey    = "branch"
	tokenTypeKey = "type"
	expKey       = "exp"
	freshKey     = "fresh"
)

var (
	secret []byte
)

func NewJwt() JWTAssumer {
	return &jwtUtils{}
}

func Init() {
	secret = []byte(os.Getenv(secretKey))
	if string(secret) == "" {
		log.Fatal("Secret key tidak boleh kosong, ENV : SECRET_KEY")
	}
}

type JWTAssumer interface {
	GenerateToken(claims CustomClaim) (string, rest_err.APIError)
	ValidateToken(tokenString string) (*jwt.Token, rest_err.APIError)
	ReadToken(token *jwt.Token) (*CustomClaim, rest_err.APIError)
}

type jwtUtils struct {
}

// GenerateToken membuat token jwt untuk login header, untuk menguji nilai payloadnya
// dapat menggunakan situs jwt.io
func (j *jwtUtils) GenerateToken(claims CustomClaim) (string, rest_err.APIError) {
	expired := time.Now().Add(time.Minute * claims.ExtraMinute).Unix()

	jwtClaim := jwt.MapClaims{}
	jwtClaim[identityKey] = claims.Identity
	jwtClaim[nameKey] = claims.Name
	jwtClaim[rolesKey] = claims.Roles
	jwtClaim[branchKey] = claims.Branch
	jwtClaim[expKey] = expired
	jwtClaim[tokenTypeKey] = claims.Type
	jwtClaim[freshKey] = claims.Fresh

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		logger.Error("gagal menandatangani token", err)
		return "", rest_err.NewInternalServerError("gagal menandatangani token", err)
	}

	return signedToken, nil
}

// ReadToken membaca inputan token dan menghasilkan pointer struct CustomClaim
// struct CustomClaim digunakan untuk nilai passing antar middleware
func (j *jwtUtils) ReadToken(token *jwt.Token) (*CustomClaim, rest_err.APIError) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Error("gagal mapping token atau token tidak valid", nil)
		return nil, rest_err.NewInternalServerError("gagal mapping token", nil)
	}

	customClaim := CustomClaim{
		Identity: cast.ToString(claims[identityKey]),
		Name:     cast.ToString(claims[nameKey]),
		Exp:      cast.ToInt64(claims[expKey]),
		Roles:    cast.ToStringSlice(claims[rolesKey]),
		Branch:   cast.ToString(claims[branchKey]),
		Type:     cast.ToInt(claims[tokenTypeKey]),
		Fresh:    cast.ToBool(claims[freshKey]),
	}
	return &customClaim, nil
}

//func iToSliceString(assumedSliceInterface interface{}) []string {
//	sliceInterface := assumedSliceInterface.([]interface{})
//	sliceString := make([]string, len(sliceInterface))
//	for i, v := range sliceInterface {
//		sliceString[i] = v.(string)
//	}
//
//	return sliceString
//}

// ValidateToken memvalidasi apakah token string masukan valid, termasuk memvalidasi apabila field exp nya kadaluarsa
func (j *jwtUtils) ValidateToken(tokenString string) (*jwt.Token, rest_err.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, rest_err.NewAPIError("Token signing method salah", http.StatusUnprocessableEntity, "jwt_error", nil)
		}
		return secret, nil
	})

	// Jika expired akan muncul disini asalkan ada claims exp
	if err != nil {
		return nil, rest_err.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{err.Error()})
	}

	return token, nil
}
