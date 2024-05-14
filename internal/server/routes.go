package server

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	validate.RegisterValidation("nip", nipValidation)
}

func nipValidation(fl validator.FieldLevel) bool {
	nip := fl.Field().Int()
	role := fl.Param()

	// skip last three digits
	nip = nip / 1000

	// check month
	month := nip % 100
	if month < 1 || month > 12 {
		return false
	}
	nip = nip / 100

	// check year
	year := nip % 1000
	if year < 2000 || year > int64(time.Now().Year()) {
		return false
	}
	nip = nip / 1000

	// check gender
	gender := nip % 10
	if gender < 1 || gender > 2 {
		return false
	}
	nip = nip / 10

	// check role
	itRole := role == "it" && nip == 615
	nurseRole := role == "nurse" && nip == 303
	if !itRole {
		return false
	} else if !nurseRole {
		return false
	}

	return true
}
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
