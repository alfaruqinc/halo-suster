package handler

import (
	"health-record/internal/domain"
	"health-record/internal/helper"
	"health-record/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MedicalPatient interface {
	Create() fiber.Handler
}

type medicalPatient struct {
	validator             *validator.Validate
	medicalPatientService service.MedicalPatient
}

func NewMedicalPatient(validator *validator.Validate, medicalPatientService service.MedicalPatient) MedicalPatient {
	return &medicalPatient{
		validator:             validator,
		medicalPatientService: medicalPatientService,
	}
}

func (mp *medicalPatient) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(domain.CreateMedicalPatient)
		ctx.BodyParser(body)

		if err := mp.validator.Struct(body); err != nil {
			err := helper.ValidateRequest(err)
			return ctx.Status(err.Status()).JSON(err)
		}

		patient := body.NewMedicalPatienFromDTO()
		err := mp.medicalPatientService.Create(ctx.Context(), patient)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		respCreated := domain.NewStatusCreated("success create medical patient", "")

		return ctx.Status(respCreated.Status()).JSON(respCreated)
	}
}
