package middleware

import (
	"health-record/internal/domain"
	"path/filepath"
	"slices"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Auth interface {
	Auth() fiber.Handler
}

type auth struct {
	jwtSecret string
}

func NewAuth(jwtSecret string) Auth {
	return &auth{
		jwtSecret: jwtSecret,
	}
}

var (
	SkippedPath = []string{
		"/v1/user/it/register",
		"/v1/user/it/login",
		"/v1/user/nurse/login",
	}
	OnlyITRole = []string{
		"/v1/user/nurse/*",
		"/v1/user/nurse/register",
		"/v1/user/nurse/*/access",
		"/v1/user",
	}
)

func (a *auth) Auth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		Filter: func(ctx *fiber.Ctx) bool {
			if slices.Contains(SkippedPath, ctx.Path()) {
				return true
			}
			return false
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			user := ctx.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)

			role := claims["role"].(string)
			path := ctx.Path()
			onlyIT := false

			// Check if path is in it role paths
			for _, p := range OnlyITRole {
				match, _ := filepath.Match(p, path)
				if match {
					onlyIT = true
					break
				}
			}
			// Check if user has access to only it role paths
			if onlyIT && role != "it" {
				unauthorized := domain.NewErrUnauthorized("Only IT Role")
				return ctx.Status(unauthorized.Status()).JSON(unauthorized)
			}

			return ctx.Next()
		},
		SigningKey:   jwtware.SigningKey{Key: []byte(a.jwtSecret)},
		ErrorHandler: customAuthError,
	})
}

func customAuthError(ctx *fiber.Ctx, err error) error {
	if err.Error() == jwtware.ErrJWTMissingOrMalformed.Error() {
		missingOrMalformed := domain.NewErrUnauthorized(jwtware.ErrJWTMissingOrMalformed.Error())
		return ctx.Status(missingOrMalformed.Status()).JSON(missingOrMalformed)
	}
	invalidOrExpired := domain.NewErrUnauthorized("Invalid or expired JWT")
	return ctx.Status(invalidOrExpired.Status()).JSON(invalidOrExpired)
}
