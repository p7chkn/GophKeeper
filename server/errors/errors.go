package customerrors

import (
	"fmt"
)

func NewCustomError(err error, statusCode string) error {
	return &CustomError{
		Err:        err,
		StatusCode: statusCode,
	}
}

func ParseError(err error) string {
	switch e := err.(type) {
	case *CustomError:
		return e.StatusCode
	default:
		return "internal server error"
	}
}

type CustomError struct {
	Err        error
	StatusCode string
}

func (err *CustomError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *CustomError) Unwrap() error {
	return err.Err
}
