package database

import (
	"covid-19/internal/model"

	"github.com/pkg/errors"
)

// CreateQuestion inserts new question
func (db Orm) CreateQuestion(q *model.Question) (inserted *model.Question, err error) {
	if _, err := db.Table("Question").Insert(q); err != nil {
		return nil, errors.WithMessage(err, "Error Creating Question - Database Error")
	}
	return q, nil
}

// CreateQuestionSection inserts new question section
func (db Orm) CreateQuestionSection(qs *model.QuestionSection) (inserted *model.QuestionSection, err error) {
	if _, err := db.Table("QuestionSection").Insert(qs); err != nil {
		return nil, errors.WithMessage(err, "Error Creating Question Section - Database Error")
	}
	return qs, nil
}

// ListQuestions gets Questions from db
// TODO: Support Pagination
// TODO: Support Search
func (db Orm) ListQuestions() (qs []*model.Question, err error) {
	qs = []*model.Question{}
	err = db.Table("Question").Find(&qs)
	return qs, errors.WithMessage(err, "Error Getting Question - Database Error")
}

// ListQuestionSections gets Question Sections from db
// TODO: Support Pagination
// TODO: Support Search
func (db Orm) ListQuestionSections() (qs []*model.QuestionSection, err error) {
	qs = []*model.QuestionSection{}
	err = db.Table("QuestionSection").Find(&qs)
	return qs, errors.WithMessage(err, "Error Getting Question - Database Error")
}

// GetQuestionByID gets Question from db by id
func (db Orm) GetQuestionByID(id string) (q *model.Question, err error) {
	q = &model.Question{}
	var has bool
	if has, err = db.Table("Question").Where("id = ?", id).Get(q); err != nil {
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
