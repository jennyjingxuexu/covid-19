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
	ID           string             `xorm:"id" json:"id"`
	QuestionType string             `xorm:"question_type" json:"question_type" r-validate:"required,question-type"`
	Question     string             `xorm:"question" json:"question" r-validate:"required"`
	Choices      map[string]*Choice `xorm:"choices" json:"choices" r-validate:"required"`
	MaxPoint     int                `xorm:"-> max_point" json:"-"`
	// TODO: should just use QuestionSection Model, but xorm right now has a bug with "extends" tag.
	QuestionSectionID   string `xorm:"-> question_section_id" json:"question_section_id,omitempty" r-validate:"uuid,required"`
	QuestionSectionName string `xorm:"<- question_section_name" json:"question_section_name,omitempty" r-validate:"-"`
}

// Choice is what a question's choices can be, it includes a choice text and then a point assigned to the choice
type Choice struct {
	Value interface{} `json:"value,omitempty" r-validate:"required"`
	Point int         `json:"point,omitempty" r-validate:"required"`
}

// ValidateQuestionRequest validates the Question struct as the Question was constructed by the http request
// TODO: Need to better organize the code, maybe we can make the validation step more abstract.
func ValidateQuestionRequest(q Question) error {
	v := RequestValidator()
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

	// Quetion Type Based validation
	// TODO: separate logic to different functions
	if q.QuestionType == "TRUE_FALSE" {
		if len(q.Choices) != 2 {
			return errors.New("TRUE_FALSE questions must have exactly two choices")
		}
		var check *bool
		for _, c := range q.Choices {
			bc, ok := c.Value.(bool)
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
