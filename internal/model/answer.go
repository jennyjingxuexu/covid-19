package model

import "github.com/pkg/errors"

// Answer by the user
type Answer struct {
	ID         string `xorm:"id" json:"id"`
	UserID     string `xorm:"user_id" json:"user_id"`
	QuestionID string `xorm:"question_id" json:"question_id" r-validate:"uuid,required"`
	Choice     string `xorm:"choice" json:"choice" r-validate:"required"`
	Point      int    `xorm:"point" json:"-"`
	// TODO: should just populate it from db
	PossiblePoint int `xorm:"-" json:"-"`
}

type questionGetter interface {
	GetQuestionByID(id string) (q *Question, err error)
}

// ValidateAndPopulateAnswerRequest validates the Answer struct as the Answer was constructed by the http request
// Also populates the Point field based on the validation.
func ValidateAndPopulateAnswerRequest(a *Answer, qg questionGetter) error {
	v := RequestValidator()
	err := TranslateError(v.Struct(a))
	if err != nil {
		return err
	}
	q, err := qg.GetQuestionByID(a.QuestionID)
	// TODO: Need to better handle this. Panic is not good, even if we have a panic handler.
	if err != nil {
		panic(err)
	}
	if q == nil {
		return errors.Errorf("question_id(%s) does not exist", a.QuestionID)
	}
	if c, ok := q.Choices[a.Choice]; !ok {
		return errors.Errorf("question_id(%s) does not have choice %s", a.QuestionID, a.Choice)
	} else {
		a.PossiblePoint = q.MaxPoint
		a.Point = c.Point
	}
	return nil
}
