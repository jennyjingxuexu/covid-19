package model

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var qTypes = []string{"TRUE_FALSE", "SINGLE_SELECT"}

// Question in the bank
// Currently the only implementation supported is choice based questions.
type Question struct {
	ID                string                 `xorm:"id" json:"id"`
	QuestionType      string                 `xorm:"question_type" json:"question_type" r-validate:"required,question-type"`
	Question          string                 `xorm:"question" json:"question" r-validate:"required"`
	QuestionSectionID string                 `xorm:"question_section_id ->" json:"question_section_id,omitempty" r-validate:"uuid,required"`
	QuestionSection   QuestionSection        `xorm:"question_section <- extends" json:"question_section,omitempty" r-validate:"-"`
	CorrectChoice     string                 `xorm:"correct_choice" json:"correct_choice" r-validate:"required"`
	Choices           map[string]interface{} `xorm:"choices" json:"choices" r-validate:"required"`
}

// ValidateQuestionRequest validates the Question struct as the Question was constructed by the http request
// TODO: Need to better organize the code, maybe we can make the validation step more abstract.
func ValidateQuestionRequest(q Question) error {
	v := ReuqestValidator()
	v.RegisterValidation("question-type", questionType)
	v.RegisterTranslation("question-type", *ValidatorTranslator(), func(ut ut.Translator) error {
		return ut.Add("question-type", fmt.Sprintf("{0} must be one of %v", qTypes), true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("question-type", fe.Field())
		return t
	})

	err := TranslateError(v.Struct(q))
	if err != nil {
		return err
	}

	// TODO: ALL following can be refactored to use validator package, should do it.

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

func questionType(fl validator.FieldLevel) bool {
	st := fl.Field()
	qType := st.String()
	return contains(qType, qTypes)
}

func contains(str string, strs []string) bool {
	for _, s := range strs {
		if str == s {
			return true
		}
	}
	return false
}
