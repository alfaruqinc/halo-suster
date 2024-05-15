package server

import (
	"health-record/internal/handler"
	"health-record/internal/repository"
	"health-record/internal/service"
	"health-record/internal/validation"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	validate   = validator.New()
	jwtSecret  = os.Getenv("JWT_SECRET")
	bcryptSalt = os.Getenv("BCRYPT_SALT")
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	validate.RegisterValidation("nip", validation.NIP)
	validate.RegisterValidation("url", validation.URL)
	validate.RegisterValidation("intlen", validation.IntLen)
	validate.RegisterValidation("iso8601", validation.ISO8601)

	userITRepository := repository.NewUserIT()
	userNurseRepository := repository.NewUserNurse()

	userITService := service.NewUserIT(s.db.GetDB(), jwtSecret, bcryptSalt, userITRepository)
	userNurseService := service.NewUserNurse(s.db.GetDB(), jwtSecret, bcryptSalt, userNurseRepository)

	userITHandler := handler.NewUserIT(validate, userITService)
	userNurseHandler := handler.NewUserNurse(validate, userNurseService)

	s.App.Use(recover.New())

	apiV1 := s.App.Group("/v1")

	user := apiV1.Group("/user")

	it := user.Group("/it")
	it.Post("/register", userITHandler.Register())
	it.Post("/login", userITHandler.Login())

	nurse := user.Group("/nurse")
	nurse.Post("/login", userNurseHandler.Login())
	nurse.Post("/register", userNurseHandler.Register())
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
