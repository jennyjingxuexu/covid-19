package model

import (
	"reflect"

	"github.com/pkg/errors"
)

var qTypes = []string{"TRUE_FALSE", "SINGLE_SELECT"}

// Question in the bank
// Currently the only implementation supported is choice based questions.
type Question struct {
	ID            string                 `xorm:"id" json:"id"`
	QuestionType  string                 `xorm:"question_type" json:"question_type" r-validate:"nonzero,question-type"`
	Question      string                 `xorm:"question" json:"question,omitempty" r-validate:"nonzero"`
	CorrectChoice string                 `xorm:"correct_choice" json:"correct_choice" r-validate:"nonzero"`
	Choices       map[string]interface{} `xorm:"choices" json:"choices" r-validate:"nonzero"`
}

// ValidateQuestionRequest validates the Question struct as the Question was constructed by the http request
// TODO: Need to better organize the code, maybe we can make the validation step more abstract.
func ValidateQuestionRequest(q Question) error {
	validator := ReuqestValidator()
	validator.SetValidationFunc("question-type", questionType)
	err := validator.Validate(q)
	if err != nil {
		return err
	}

	// cross field validation
	// TODO: maybe we need to support more than choice based questions.
	if _, ok := q.Choices[q.CorrectChoice]; !ok {
		return errors.New("correct_choice must be one of the choices in the choices object")
	}

	// Quetion Type Based validation
	// TODO: separate logic to different functions
	if q.QuestionType == "TRUE_FALSE" {
		if len(q.Choices) != 2 {
			return errors.New("TRUE_FALSE questions must have exactly two choices")
		}
		var check *bool
		for _, c := range q.Choices {
			bc, ok := c.(bool)
			if !ok {
				return errors.New("TRUE_FALSE questions must have boolean as choices")
			}
			if check == nil {
				check = &bc
			} else {
				newCheck := (*check || bc) && !(*check && bc)
				check = &newCheck
			}
		}
		if !*check {
			return errors.New("TRUE_FALSE questions must different boolean as choices")
		}
	}
	if q.QuestionType == "SINGLE_SELECT" {
		if len(q.Choices) < 2 {
			return errors.New("SINGLE_SELECT questions must have at least two choices")
		}

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
