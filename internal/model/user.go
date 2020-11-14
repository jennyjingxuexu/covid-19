package model

import (
	"regexp"
	"unicode"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// User of the app
type User struct {
	ID       string `xorm:"id" json:"id"`
	Username string `xorm:"username" json:"username" r-validate:"username"`
	Password string `xorm:"password" json:"password,omitempty" r-validate:"password"`
}

// ValidateUserRequest validates the User struct as the User was constructed by the http request
// TODO: Need to better organize the code, maybe we can make the validation step more abstract.
func ValidateUserRequest(u User) error {
	v := RequestValidator()
	v.RegisterValidation("username", username)
	v.RegisterValidation("password", password)
	v.RegisterTranslation("username", *ValidatorTranslator(), func(ut ut.Translator) error {
		return ut.Add("username", "{0} must have a minimum of 3 maximum of 16 letters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("username", fe.Field())
		return t
	})
	v.RegisterTranslation("password", *ValidatorTranslator(), func(ut ut.Translator) error {
		return ut.Add("password", "{0} must contain 1 special character, 1 uppercase, 1 lower case letter and have a minimum of 3 maximum of 16 letters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})
	return TranslateError(v.Struct(u))
}

func username(fl validator.FieldLevel) bool {
	st := fl.Field()

	// alphanumeric with '-' and '_', min 3, max 16
	const usernameString = "^[a-z0-9_-]{3,16}$"
	usernameRegEx := regexp.MustCompile(usernameString)

	return usernameRegEx.MatchString(st.String())
}

func password(fl validator.FieldLevel) bool {
	st := fl.Field()
	p := st.String()
	if len(p) < 6 || len(p) > 32 {
		return false
	}
	var lower, upper, digit, special bool
	for _, r := range p {
		// TODO: Consider a swich
		if !lower && unicode.IsLower(r) {
			lower = true
		}
		if !upper && unicode.IsUpper(r) {
			upper = true
		}
		if !digit && unicode.IsDigit(r) {
			digit = true
		}
		if !special && (unicode.IsPunct(r) || unicode.IsSymbol(r)) {
			special = true
		}
	}
	if !(lower && upper && digit && special) {
		return false
	}
	return true
}
