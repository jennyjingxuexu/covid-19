package handler

import (
	"covid-19/internal/model"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// questionService is an interface describing all operations related to User
type questionService interface {
	CreateQuestion(*model.Question) (*model.Question, error)
	CreateQuestionSection(*model.QuestionSection) (*model.QuestionSection, error)
	ListQuestions() ([]*model.Question, error)
	ListQuestionSections() ([]*model.QuestionSection, error)
	GetQuestionByID(id string) (*model.Question, error)
	GetQuestionSectionByID(string) (*model.QuestionSection, error)
	GetQuestionSectionByName(string) (*model.QuestionSection, error)
}

// QuestionProvider provides handlers for handling question related http requests
type QuestionProvider struct {
	question questionService
}

// ListQuestions list all questions
// TODO: Support pagination
func (provider QuestionProvider) ListQuestions() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if qs, err := provider.question.ListQuestions(); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			json.NewEncoder(w).Encode(qs)
		}
	})
}

// ListQuestionSections list all question sections
// TODO: Support pagination
func (provider QuestionProvider) ListQuestionSections() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if qs, err := provider.question.ListQuestionSections(); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			json.NewEncoder(w).Encode(qs)
		}
	})
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

		if has, _ := provider.question.GetQuestionSectionByID(q.QuestionSectionID); has == nil {
			err = errors.Wrap(errors.Errorf("Question Section with id(%v) does not exist", q.QuestionSectionID), "Invalid Request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q.ID = uuid.New().String()
		q.MaxPoint = 0
		for _, choice := range q.Choices {
			if choice.Point > q.MaxPoint {
				q.MaxPoint = choice.Point
			}
		}

		if q, err = provider.question.CreateQuestion(q); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(q)
	})
}

// CreateQuestionSection creates a QuestionSection
func (provider QuestionProvider) CreateQuestionSection() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qs := &model.QuestionSection{}
		var err error
		if err = json.NewDecoder(r.Body).Decode(qs); err != nil {
			log.Error(err)
			err = errors.New("cannot json decode the body, invalid json or missing field")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = model.ValidateQuestionSectionRequest(*qs); err != nil {
			err = errors.Wrap(err, "Invalid Request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if has, _ := provider.question.GetQuestionSectionByName(qs.Name); has != nil {
			err = errors.Wrap(errors.Errorf("Question Section with name(%v) already exist", qs.Name), "Invalid Request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		qs.ID = uuid.New().String()

		if qs, err = provider.question.CreateQuestionSection(qs); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(qs)
	})
}
