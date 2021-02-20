package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"strings"
)

var (
	jwt = mjwt.NewJwt()
)

const (
	headerKey = "Authorization"
	bearerKey = "Bearer"
)

type role int

const (
	normalAuth role = iota
	freshAuth
	adminAuth
)

//AuthMiddleware memvalidasi token JWT, mengembalikan claims berupa pointer mjwt.CustomClaims
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get(headerKey)
	claims, err := authValidator(authHeader, normalAuth)
	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	c.Locals(mjwt.CLAIMS, claims)
	return c.Next()
}

//AuthMustFreshMiddleware memvalidasi token JWT, mengembalikan claims berupa pointer mjwt.CustomClaims.
//perbedaannya dengan AuthMidlleware adalah ini mengharuskan pengakses berstatus fresh true
func AuthMustFreshMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get(headerKey)
	claims, err := authValidator(authHeader, freshAuth)
	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	c.Locals(mjwt.CLAIMS, claims)
	return c.Next()
}

//AuthAdminMiddleware memvalidasi token JWT, mengembalikan claims berupa pointer mjwt.CustomClaims.
//perbedaannya dengan AuthMidlleware adalah ini mengharuskan pengakses berstatus is_admin true
func AuthAdminMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get(headerKey)
	claims, err := authValidator(authHeader, adminAuth)
	if err != nil {
		return c.Status(err.Status()).JSON(err)
	}

	c.Locals(mjwt.CLAIMS, claims)
	return c.Next()
}

func authValidator(authHeader string, role role) (*mjwt.CustomClaim, rest_err.APIError) {
	if !strings.Contains(authHeader, bearerKey) {
		apiErr := rest_err.NewUnauthorizedError("Unauthorized")
		return nil, apiErr
	}

	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) != 2 {
		apiErr := rest_err.NewUnauthorizedError("Unauthorized")
		return nil, apiErr
	}

	token, apiErr := jwt.ValidateToken(tokenString[1])
	if apiErr != nil {
		return nil, apiErr
	}

	claims, apiErr := jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}

	switch role {
	case adminAuth:
		if !claims.IsAdmin {
			apiErr := rest_err.NewUnauthorizedError("Unauthorized, memerlukan hak akses admin")
			return nil, apiErr
		}
	case freshAuth:
		if !claims.Fresh {
			apiErr := rest_err.NewUnauthorizedError("Memerlukan token yang baru untuk mengakses halaman ini")
			return nil, apiErr
		}
	}

	return claims, nil
}
