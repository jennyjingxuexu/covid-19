package database

import (
	"covid-19/internal/model"
	"time"

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
	if _, err := db.Table("UserSession").Insert(u); err != nil {
		return nil, errors.WithMessage(err, "Error Creating UserSession - Database Error")
	}
	return u, nil
}

// RefreshUserSession refreshes the user session.
func (db Orm) RefreshUserSession(sessionID string) (session *model.UserSession, err error) {
	session = &model.UserSession{}
	if has, err := db.Table("UserSession").Where("id = ?", sessionID).Get(session); !has || err != nil {
		if !has {
			return nil, nil
		} else {
			return nil, errors.WithMessage(err, "Error Getting UserSession - Database Error")
		}
	}
	if time.Since(session.LastSeenTime).Hours() > 2 {
		return nil, nil
	}
	session.LastSeenTime = time.Now()
	if _, err := db.Table("UserSession").Update(session, &model.UserSession{ID: session.ID}); err != nil {
		return nil, errors.WithMessage(err, "Error Updating UserSession - Database Error")
	}
	return session, nil
}
