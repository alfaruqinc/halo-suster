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
	Update() fiber.Handler
	Delete() fiber.Handler
	GiveAccess() fiber.Handler
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

		if err := un.validator.Struct(body); err != nil {
			err := helper.ValidateRequest(err)
			return ctx.Status(err.Status()).JSON(err)
		}

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

func (un *userNurse) Update() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		nurseId := ctx.Params("nurseId")

		body := new(domain.UpdateUserNurse)
		ctx.BodyParser(body)

		if err := un.validator.Struct(body); err != nil {
			err := helper.ValidateRequest(err)
			return ctx.Status(err.Status()).JSON(err)
		}

		body.ID = nurseId
		err := un.userNurseService.Update(ctx.Context(), *body)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		resp := domain.NewStatusOK("success update user nurse", "")

		return ctx.Status(resp.Status()).JSON(resp)
	}
}

func (un *userNurse) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		nurseId := ctx.Params("nurseId")

		err := un.userNurseService.Delete(ctx.Context(), nurseId)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		respDelete := domain.NewStatusOK("success delete nurse", "")

		return ctx.Status(respDelete.Status()).JSON(respDelete)
	}
}

func (un *userNurse) GiveAccess() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		nurseId := ctx.Params("nurseId")

		body := new(domain.AccessSystemUserNurse)
		ctx.BodyParser(body)

		if err := un.validator.Struct(body); err != nil {
			err := helper.ValidateRequest(err)
			return ctx.Status(err.Status()).JSON(err)
		}

		body.ID = nurseId
		err := un.userNurseService.GiveAccess(ctx.Context(), *body)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		resp := domain.NewStatusOK("success give nurse access", "")

		return ctx.Status(resp.Status()).JSON(resp)
	}
}
