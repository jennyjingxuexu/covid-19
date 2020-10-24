package database

import (
	"covid-19/internal/model"

	"github.com/pkg/errors"
)

// CreateUser Inserts new user
func (db Orm) CreateUser(u *model.User) (inserted *model.User, err error) {
	if _, err := db.Table("User").Insert(u); err != nil {
		return nil, errors.WithMessage(err, "Error Creating User - Database Error")
	}
	return u, nil
}

// GetUserByUsername gets user from db by username
func (db Orm) GetUserByUsername(username string) (u *model.User, err error) {
	u = &model.User{}
	var has bool
	if has, err = db.Table("User").Where("username = ?", username).Get(u); err != nil {
		return nil, errors.WithMessage(err, "Error Getting User - Database Error")
	}
	if !has {
		return nil, nil
	}
	return u, nil
}

// CreateUserSession creates a new user session
func (db Orm) CreateUserSession(u *model.UserSession) (inserted *model.UserSession, err error) {
	u = &model.UserSession{}
	if _, err := db.Table("UserSession").Insert(u); err != nil {
		return nil, errors.WithMessage(err, "Error Creating UserSession - Database Error")
	}
	return u, nil
}
