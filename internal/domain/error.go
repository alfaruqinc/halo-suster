package domain

import "github.com/gofiber/fiber/v2"

type ErrorMessage interface {
	Message() string
	Status() int
	Error() string
}

type errorData struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
}

func (e *errorData) Message() string {
	return e.ErrMessage
}

func (e *errorData) Status() int {
	return e.ErrStatus
}

func (e *errorData) Error() string {
	return e.ErrError
}

func NewErrUnauthorized(message string) ErrorMessage {
	return &errorData{
		ErrMessage: message,
		ErrStatus:  fiber.StatusUnauthorized,
		ErrError:   "UNAUTHORIZED",
	}
}

func NewErrNotFound(message string) ErrorMessage {
	return &errorData{
		ErrMessage: message,
		ErrStatus:  fiber.StatusNotFound,
		ErrError:   "NOT_FOUND",
	}
}

func NewErrBadRequest(message string) ErrorMessage {
	return &errorData{
		ErrMessage: message,
		ErrStatus:  fiber.StatusBadRequest,
		ErrError:   "BAD_REQUEST",
	}
}

func NewErrConflict(message string) ErrorMessage {
	return &errorData{
		ErrMessage: message,
		ErrStatus:  fiber.StatusConflict,
		ErrError:   "CONFLICT",
	}
}

func NewErrInternalServerError(message string) ErrorMessage {
	return &errorData{
		ErrMessage: message,
		ErrStatus:  fiber.StatusInternalServerError,
		ErrError:   "INTERNAL_SERVER_ERROR",
	}
}
