package handler

import (
	"health-record/internal/domain"
	"health-record/internal/helper"
	"health-record/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserIT interface {
	Register() fiber.Handler
}

type userIT struct {
	validator   *validator.Validate
	userService service.UserIT
}

func NewUserIT(validator *validator.Validate, userSerivce service.UserIT) UserIT {
	return &userIT{
		validator:   validator,
		userService: userSerivce,
	}
}

func (uit *userIT) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(domain.RegisterUserIT)
		ctx.BodyParser(&body)

		if err := uit.validator.Struct(body); err != nil {
			err := helper.ValidateRequest(err)
			return ctx.Status(err.Status()).JSON(err)
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
