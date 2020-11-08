package model

import (
	"gopkg.in/validator.v2"
)

// ReuqestValidator returns a validator that validates based on struct tag r-validate.
func ReuqestValidator() *validator.Validator {
	v := validator.NewValidator()
	v.SetTag("r-validate")
	v.SetPrintJSON(true)
	return v
}
