package handler

import (
	"health-record/internal/domain"
	"health-record/internal/helper"
	"health-record/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type MedicalRecord interface {
	Create() fiber.Handler
}

type medicalRecord struct {
	validator            *validator.Validate
	medicalRecordService service.MedicalRecord
}

func NewMedicalRecord(validator *validator.Validate, medicalRecordService service.MedicalRecord) MedicalRecord {
	return &medicalRecord{
		validator:            validator,
		medicalRecordService: medicalRecordService,
	}
}

func (mr *medicalRecord) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		body := new(domain.CreateMedicalRecord)
		ctx.BodyParser(body)

		if err := mr.validator.Struct(body); err != nil {
			err := helper.ValidateRequest(err)
			return ctx.Status(err.Status()).JSON(err)
		}

		record := body.NewMedicalRecordFromDTO()
		record.CreatedByID = claims["id"].(string)
		err := mr.medicalRecordService.Create(ctx.Context(), record)
		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		respCreated := domain.NewStatusCreated("success create medical record", "")

		return ctx.Status(respCreated.Status()).JSON(respCreated)
	}
}
