package model

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: put these in their own package
// global varialbe is not desired.
var (
	validate  *validator.Validate
	translate *ut.Translator
)

// ReuqestValidator returns a validator that validates based on struct tag r-validate.
func ReuqestValidator() *validator.Validate {
	if validate != nil {
		return validate
	}
	validate = validator.New()
	validate.SetTagName("r-validate")
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	vt := *ValidatorTranslator()
	en_translations.RegisterDefaultTranslations(validate, vt)
	translate = &vt

	return validate
}

// ValidatorTranslator returns the cached translator
func ValidatorTranslator() *ut.Translator {
	if translate != nil {
		return translate
	}
	en := en.New()
	uni := ut.New(en, en)

	t, _ := uni.GetTranslator("en")
	translate = &t
	return translate
}

// TranslateError translates the validation error.
func TranslateError(err error) error {
	if err == nil {
		return nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		log.Error("cannot cast the following error")
		log.Error(err)
		panic(err)
	}
	msg := []string{}
	for _, e := range errs {
		msg = append(msg, e.Translate(*translate))
	}
	return errors.New(strings.Join(msg, ", "))
}
