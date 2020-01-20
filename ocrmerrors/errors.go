package ocrmerrors

import (
	"errors"
	"fmt"
)

var (
	ErrParseConfig        = errors.New("error parse config for profile")
	ErrParseConfigEquals0 = errors.New("error parse config for profile, equals 0")
)

type ErrorStruct struct {
	Code      Code
	Message   string
	MessageRu string
}

func New(code Code, message, messageRu string) *ErrorStruct {
	return &ErrorStruct{
		Code:      code,
		Message:   message,
		MessageRu: messageRu,
	}
}

func (err *ErrorStruct) Error() string {
	if err.MessageRu == "" {
		return err.Message
	}
	return fmt.Sprintf("En: %s; Ru: %s", err.Message, err.MessageRu)
}
