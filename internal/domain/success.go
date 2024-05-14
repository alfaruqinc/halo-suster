package domain

import "github.com/gofiber/fiber/v2"

type SuccessMessage interface {
	Message() string
	Status() int
	Data() any
}

type successData struct {
	SuccessMessage string `json:"message"`
	SuccessStatus  int    `json:"status"`
	SuccessData    any    `json:"data,omitempty"`
}

func (s *successData) Message() string {
	return s.SuccessMessage
}

func (s *successData) Status() int {
	return s.SuccessStatus
}

func (s *successData) Data() any {
	return s.SuccessData
}

func NewStatusOK(message string, data any) SuccessMessage {
	return &successData{
		SuccessMessage: message,
		SuccessStatus:  fiber.StatusOK,
		SuccessData:    data,
	}
}

func NewStatusCreated(message string, data any) SuccessMessage {
	return &successData{
		SuccessMessage: message,
		SuccessStatus:  fiber.StatusCreated,
		SuccessData:    data,
	}
}
