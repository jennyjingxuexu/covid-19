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
type questionService interface {
	CreateQuestion(*model.Question) (*model.Question, error)
	GetQuestionByID(string) (*model.Question, error)
}

// QuestionProvider provides handlers for handling question related http requests
type QuestionProvider struct {
	question questionService
}

// CreateQuestion creates a Question
func (provider QuestionProvider) CreateQuestion() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := &model.Question{}
		var err error
		if err = json.NewDecoder(r.Body).Decode(q); err != nil {
			log.Error(err)
			err = errors.New("cannot json decode the body, invalid json or missing field")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = model.ValidateQuestionRequest(*q); err != nil {
			err = errors.Wrap(err, "Invalid Request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q.ID = uuid.New().String()

		if q, err = provider.question.CreateQuestion(q); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Hide encrypted password from the response
		json.NewEncoder(w).Encode(q)
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
