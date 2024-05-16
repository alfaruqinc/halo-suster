package helper

import (
	"fmt"
	"health-record/internal/domain"

	"github.com/go-playground/validator/v10"
)

func msgForTag(fe validator.FieldError) string {
	field := fe.Field()
	param := fe.Param()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("minimum %s is %s", field, param)
	case "max":
		return fmt.Sprintf("maximum %s is %s", field, param)
	case "phonenumber", "nip":
		return fmt.Sprintf("wrong %s format", field)
	case "url":
		return fmt.Sprintf("%s should be url", field)
	case "oneof":
		return fmt.Sprintf("%s should be one of this value %s", field, param)
	}

	return "unhandled validation"
}

func ValidateRequest(err error) domain.ErrorMessage {
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range err {
			return domain.NewErrBadRequest(msgForTag(fe))
		}
	}
	return domain.NewErrBadRequest(err.Error())
}
