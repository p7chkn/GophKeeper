// Package customerrors пакет для работы с своими ошибками
package customerrors

import (
	"fmt"
)

// NewCustomError функция получения новой ошибки
func NewCustomError(err error, statusCode string) error {
	return &CustomError{
		Err:        err,
		StatusCode: statusCode,
	}
}

// ParseError функция по распаковки ошибки в строку
func ParseError(err error) string {
	switch e := err.(type) {
	case *CustomError:
		return e.StatusCode
	default:
		return "internal server error"
	}
}

// CustomError структура для хранения собвтенных ошибок
type CustomError struct {
	Err        error
	StatusCode string
}

// Error функция для возвранеия ошибки
func (err *CustomError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

// Unwrap функция для распоковки ошибки
func (err *CustomError) Unwrap() error {
	return err.Err
}
