package handler

import (
	"covid-19/internal/model"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// UserService is an interface describing all operations related to User
type userService interface {
	CreateUser(*model.User) (*model.User, error)
	GetUserByUsername(string) (*model.User, error)
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
				err = errors.New("Username already exist")
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
