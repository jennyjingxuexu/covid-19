package model

import (
	"reflect"

	"github.com/pkg/errors"
)

var qTypes = []string{"TRUE_FALSE", "SINGLE_SELECT"}

// Question in the bank
type Question struct {
	ID            string                 `xorm:"id" json:"id"`
	QuestionType  string                 `xorm:"question_type" json:"question_type" r-validate:"question-type"`
	Question      string                 `xorm:"question" json:"question,omitempty"`
	CorrectChoice string                 `xorm:"correct_choice" json:"correct_choce"`
	Choices       map[string]interface{} `xorm:"choices" json:"choices"`
}

// ValidateQuestionRequest validates the Question struct as the Question was constructed by the http request
// TODO: Need to better organize the code, maybe we can make the validation step more abstract.
func ValidateQuestionRequest(q Question) error {
	validator := ReuqestValidator()
	validator.SetValidationFunc("username", questionType)
	err := validator.Validate(q)
	if err != nil {
		return err
	}

	// cross field validation
	if _, ok := q.Choices[q.CorrectChoice]; !ok {
		return errors.New("correct_choice must be one of the choices in the choices object")
	}
	return nil
}

func questionType(v interface{}, _ string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("question-type validator only validates strings")
	}
	qType := st.String()

	if !contains(qType, qTypes) {
		return errors.Errorf("question_type must be one of %v", qTypes)
	}
	return nil
}

func contains(str string, strs []string) bool {
	for _, s := range strs {
		if str == s {
			return true
		}
	}
	return false
}
