package xvalidator

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validator) Validate(r *http.Request, s any) error {
	err := v.validate.StructCtx(r.Context(), s)

	switch data := s.(type) {
	case Validatable:
		return customError(data, err)
	default:
		return err
	}
}

func customError(s Validatable, err error) error {
	if err == nil {
		return nil
	}
	validationErrs, ok := err.(validator.ValidationErrors)
	// return err if err is not a list of validation errors
	if !ok {
		return err
	}

	// err is a list of validation errors
	errMsg := ""
	for _, v := range validationErrs {
		validationErrKey := fmt.Sprintf("%s.%s", v.Field(), v.Tag())
		if msg, ok := s.GetErrors()[validationErrKey]; ok {
			// custom err message found
			errMsg += fmt.Sprintf("%s; ", msg)
		} else {
			// use default err message if custom err not found
			errMsg += fmt.Sprintf("%s; ", v.Error())
		}
	}
	errMsg, _ = strings.CutSuffix(errMsg, "; ")

	return errors.New(errMsg)
}
