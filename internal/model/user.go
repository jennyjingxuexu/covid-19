package model

import (
	"reflect"
	"regexp"
	"unicode"

	"github.com/pkg/errors"
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
	validator := ReuqestValidator()
	validator.SetValidationFunc("username", username)
	validator.SetValidationFunc("password", password)
	return validator.Validate(u)
}

func username(v interface{}, _ string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("username validator only validates strings")
	}

	// alphanumeric with '-' and '_', min 3, max 16
	const usernameString = "^[a-z0-9_-]{3,16}$"
	usernameRegEx := regexp.MustCompile(usernameString)

	if !usernameRegEx.MatchString(st.String()) {
		return errors.New("Invalid Username")
	}
	return nil
}

func password(v interface{}, _ string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("password validator only validates strings")
	}
	p := st.String()
	if len(p) < 6 || len(p) > 32 {
		return errors.New("Password length must be between 6 and 32")
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
		return errors.New("Password must contain 1 lowercase letter, 1 uppercase letter, 1 digit, and 1 special charactor")
	}
	return nil
}
