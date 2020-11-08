package handler

import (
	"covid-19/internal/model"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
