package handler

import (
	"health-record/internal/domain"
	"health-record/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserIT interface {
	Register() fiber.Handler
}

type userIT struct {
	userService service.UserIT
}

func NewUserIT(userSerivce service.UserIT) UserIT {
	return &userIT{
		userService: userSerivce,
	}
}

func (uit *userIT) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(domain.RegisterUserIT)

		if err := ctx.BodyParser(&body); err != nil {
			return ctx.JSON(err)
		}

		user := body.NewUserITFromDTO()
		resp, err := uit.userService.Register(ctx.Context(), user)
		if err != nil {
			return ctx.JSON(err)
		}

		userCreated := domain.NewStatusCreated("success register", resp)

		return ctx.Status(userCreated.Status()).JSON(userCreated)
	}
}
