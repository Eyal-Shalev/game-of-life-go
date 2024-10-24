package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type StringError string

func (e StringError) Error() string {
	return string(e)
}

type JoinedErrors struct {
	errors []error
}

func (e JoinedErrors) Error() string {
	sb := strings.Builder{}
	sb.WriteString("game_of_life.JoinedErrors{")
	for i, err := range e.errors {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Error())
	}
	sb.WriteString("}")
	return sb.String()
}

func (e JoinedErrors) Unwrap() []error {
	return e.errors
}

func (e JoinedErrors) Is(target error) bool {
	for _, err := range e.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func JoinErrors(errors ...error) error {
	result := &JoinedErrors{errors: make([]error, 0, len(errors))}
	for _, err := range errors {
		if err != nil {
			result.errors = append(result.errors, err)
		}
	}
	if len(result.errors) == 0 {
		return nil
	}
	return result
}

type NilPointerError struct {
	Target any
}

func (e NilPointerError) TargetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(e.Target)).Type()
}

func (e NilPointerError) Error() string {
	return fmt.Sprintf("NilPointerError: %s", e.TargetType().Name())
}
