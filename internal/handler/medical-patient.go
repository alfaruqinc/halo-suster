package handler

import (
	"health-record/internal/domain"
	"health-record/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MedicalPatient interface {
	Create() fiber.Handler
}

type medicalPatient struct {
	medicalPatientService service.MedicalPatient
}

func NewMedicalPatient(medicalPatientService service.MedicalPatient) MedicalPatient {
	return &medicalPatient{
		medicalPatientService: medicalPatientService,
	}
}

func (mp *medicalPatient) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := new(domain.CreateMedicalPatient)
		ctx.BodyParser(body)

		patient := body.NewMedicalPatienFromDTO()
		err := mp.medicalPatientService.Create(ctx.Context(), patient)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		respCreated := domain.NewStatusCreated("success create medical patient", "")

		return ctx.Status(respCreated.Status()).JSON(respCreated)
	}
}
