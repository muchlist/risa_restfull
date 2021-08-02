package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"strings"
)

var (
	jwt = mjwt.NewJwt()
)

const (
	headerKey = "Authorization"
	bearerKey = "Bearer"
)

// NormalAuth memerlukan semua role yang tertulis agar diloloskan ke proses berikutnya
// logika yang digunakan adalah AND
func NormalAuth(rolesReq ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(headerKey)
		claims, err := authMustHaveRoleValidator(authHeader, false, rolesReq)
		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{"error": err, "data": nil})
		}
		c.Locals(mjwt.CLAIMS, claims)
		return c.Next()
	}
}

// NormalAuthOr memerlukan salah satu role yang tertulis agar diloloskan ke proses berikutnya
// logika yang digunakan adalah OR
/*func NormalAuthOr(rolesReq ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(headerKey)
		claims, err := authHaveRoleValidator(authHeader, false, rolesReq)
		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{"error": err, "data": nil})
		}
		c.Locals(mjwt.CLAIMS, claims)
		return c.Next()
	}
}*/

func FreshAuth(rolesReq ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(headerKey)
		claims, err := authMustHaveRoleValidator(authHeader, true, rolesReq)
		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{"error": err, "data": nil})
		}

		c.Locals(mjwt.CLAIMS, claims)
		return c.Next()
	}
}

// Kondisi And
func authMustHaveRoleValidator(authHeader string, mustFresh bool, rolesRequired []string) (*mjwt.CustomClaim, rest_err.APIError) {
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

	if mustFresh {
		if !claims.Fresh {
			apiErr := rest_err.NewUnauthorizedError("Memerlukan token yang baru untuk mengakses halaman ini")
			return nil, apiErr
		}
	}

	if len(rolesRequired) != 0 {
		for _, roleReq := range rolesRequired {
			if !sfunc.InSlice(roleReq, claims.Roles) {
				apiErr := rest_err.NewUnauthorizedError(fmt.Sprintf("Unauthorized, memerlukan hak akses %s", roleReq))
				return nil, apiErr
			}
		}
	}
	return claims, nil
}

// Kondisi Or
/*func authHaveRoleValidator(authHeader string, mustFresh bool, rolesAllowed []string) (*mjwt.CustomClaim, rest_err.APIError) {
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

	if mustFresh {
		if !claims.Fresh {
			apiErr := rest_err.NewUnauthorizedError("Memerlukan token yang baru untuk mengakses halaman ini")
			return nil, apiErr
		}
	}

	if len(rolesAllowed) != 0 {
		for _, roleReq := range rolesAllowed {
			if sfunc.InSlice(roleReq, claims.Roles) {
				return claims, nil
			}
		}
	} else {
		return claims, nil
	}

	apiErr = rest_err.NewUnauthorizedError(fmt.Sprintf("Unauthorized, memerlukan hak akses %s", rolesAllowed))
	return nil, apiErr
}*/
