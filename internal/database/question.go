package database

import (
	"covid-19/internal/model"

	"github.com/pkg/errors"
)

// CreateQuestion inserts new question
func (db Orm) CreateQuestion(q *model.Question) (inserted *model.Question, err error) {
	if _, err := db.Table("Question").Insert(q); err != nil {
		return nil, errors.WithMessage(err, "Error Creating User - Database Error")
	}
	return q, nil
}

// GetQuestionByID gets Question from db by id
func (db Orm) GetQuestionByID(id string) (q *model.Question, err error) {
	q = &model.Question{}
	var has bool
	if has, err = db.Table("User").Where("id = ?", id).Get(q); err != nil {
		return nil, errors.WithMessage(err, "Error Getting Question - Database Error")
	}
	if !has {
		return nil, nil
	}
	return q, nil
}

// // CreateUserSession creates a new user session
// func (db Orm) CreateUserSession(u *model.UserSession) (inserted *model.UserSession, err error) {
// 	if _, err := db.Table("UserSession").Insert(u); err != nil {
// 		return nil, errors.WithMessage(err, "Error Creating UserSession - Database Error")
// 	}
// 	return u, nil
// }
