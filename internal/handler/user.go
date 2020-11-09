package handler

import (
	"covid-19/internal/model"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// UserService is an interface describing all operations related to User
type userService interface {
	CreateUser(*model.User) (*model.User, error)
	GetUserByUsername(string) (*model.User, error)
	CreateUserSession(*model.UserSession) (*model.UserSession, error)
}

// UserProvider provides handlers for handling user related http requests
type UserProvider struct {
	user userService
}

// CreateUser creates a User
func (provider UserProvider) CreateUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		var err error
		if err = json.NewDecoder(r.Body).Decode(u); err != nil {
			log.Error(err)
			err = errors.New("cannot json decode the body, invalid json or missing field")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = model.ValidateUserRequest(*u); err != nil {
			err = errors.Wrap(err, "Invalid Request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: May want a separate endpoint for checking username taken. To get faster feed back.
		if exist, err := provider.user.GetUserByUsername(u.Username); exist != nil || err != nil {
			if err != nil {
				log.Error(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				err = errors.Wrap(errors.New("Username already exist"), "Invalid Request")
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		u.Password = string(encrypted)
		u.ID = uuid.New().String()

		if u, err = provider.user.CreateUser(u); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Hide encrypted password from the response
		u.Password = ""
		json.NewEncoder(w).Encode(u)
	})
}

// CreateUserSession creates a UserSession
// Note that we are transmitting the user auth info in the body
// TODO: far in the future, need to use Oauth2
// TODO: While on this topic, we need HTTPs if this goes public. (Which will require
// 		 some more funding than I am willing to spare right now)
func (provider UserProvider) CreateUserSession() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// also due to laziness auth is stored in user for now.
		authInfo := &model.User{}
		var err error
		if err = json.NewDecoder(r.Body).Decode(authInfo); err != nil {
			log.Error(err)
			err = errors.New("cannot json decode the body, invalid json or missing field")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var u *model.User
		if u, err = provider.user.GetUserByUsername(authInfo.Username); u == nil || err != nil {
			if err != nil {
				log.Error(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				err = errors.New("username does not exist")
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(authInfo.Password)); err != nil {
			err = errors.New("incorrect password")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s := &model.UserSession{
			ID:           uuid.New().String(),
			UserID:       u.ID,
			LoginTime:    time.Now(),
			LastSeenTime: time.Now(),
		}
		log.Error(s)
		if s, err = provider.user.CreateUserSession(s); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(s)
	})
}
