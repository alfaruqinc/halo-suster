package handler

import (
	"health-record/internal/domain"
	"health-record/internal/service"

	"github.com/gofiber/fiber/v2"
)

type User interface {
	GetAllUsers() fiber.Handler
}

type user struct {
	userService service.User
}

func NewUser(userService service.User) User {
	return &user{
		userService: userService,
	}
}

func (u *user) GetAllUsers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		queryParams := new(domain.UserQueryParams)
		ctx.QueryParser(queryParams)

		users, err := u.userService.GetAllUsers(ctx.Context(), *queryParams)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		resp := domain.NewStatusOK("success retrive all users", users)

		return ctx.Status(resp.Status()).JSON(resp)
	}
}
