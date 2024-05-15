package handler

import (
	"health-record/internal/domain"
	"health-record/internal/helper"
	"health-record/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserNurse interface {
	Register() fiber.Handler
	Login() fiber.Handler
}

type userNurse struct {
	validator        *validator.Validate
	userNurseService service.UserNurse
}

func NewUserNurse(validator *validator.Validate, userNurseService service.UserNurse) UserNurse {
	return &userNurse{
		validator:        validator,
		userNurseService: userNurseService,
	}
}

func (un *userNurse) Register() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(domain.RegisterUserNurse)
		ctx.BodyParser(&body)

		nurse := body.NewUserNurseFromDTO()
		resp, err := un.userNurseService.Register(ctx.Context(), nurse)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		nurseCreated := domain.NewStatusCreated("success register nurse", resp)

		return ctx.Status(nurseCreated.Status()).JSON(nurseCreated)
	}
}

func (un *userNurse) Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(domain.LoginUserNurse)
		ctx.BodyParser(&body)

		if err := un.validator.Struct(body); err != nil {
			err := helper.ValidateRequest(err)
			return ctx.Status(err.Status()).JSON(err)
		}

		resp, err := un.userNurseService.Login(ctx.Context(), *body)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		userLogin := domain.NewStatusOK("success login", resp)

		return ctx.Status(userLogin.Status()).JSON(userLogin)
	}
}
