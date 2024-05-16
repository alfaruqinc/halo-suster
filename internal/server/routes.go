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
	userRepository := repository.NewUser()
	medicalPatientRepository := repository.NewMedicalPatient()

	userITService := service.NewUserIT(s.db.GetDB(), jwtSecret, bcryptSalt, userITRepository)
	userNurseService := service.NewUserNurse(s.db.GetDB(), jwtSecret, bcryptSalt, userNurseRepository)
	userService := service.NewUser(s.db.GetDB(), userRepository)
	medicalPatientService := service.NewMedicalPatient(s.db.GetDB(), medicalPatientRepository)

	userITHandler := handler.NewUserIT(validate, userITService)
	userNurseHandler := handler.NewUserNurse(validate, userNurseService)
	userHandler := handler.NewUser(userService)
	medicalPatientHandler := handler.NewMedicalPatient(validate, medicalPatientService)

	s.App.Use(recover.New())

	apiV1 := s.App.Group("/v1")

	user := apiV1.Group("/user")

	it := user.Group("/it")
	it.Post("/register", userITHandler.Register())
	it.Post("/login", userITHandler.Login())

	// TODO: add auth middleware
	nurse := user.Group("/nurse")
	nurse.Post("/login", userNurseHandler.Login())
	nurse.Post("/register", userNurseHandler.Register())
	nurse.Put("/:nurseId", userNurseHandler.Update())
	nurse.Delete("/:nurseId", userNurseHandler.Delete())
	nurse.Put("/:nurseId/access", userNurseHandler.GiveAccess())

	user.Get("", userHandler.GetAllUsers())

	medical := apiV1.Group("/medical")
	medical.Post("/patient", medicalPatientHandler.Create())
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
