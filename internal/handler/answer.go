package handler

import (
	"covid-19/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// answerService is an interface describing all operations related to User
type answerService interface {
	// UpsertAnswer(*model.Answer) (*model.Answer, error)
	BulkUpsertUserAnswer(string, []*model.Answer) error
	ListAnswerByUser(string) ([]*model.Answer, error)
}

// AnswerProvider provides handlers for handling answer related http requests
type AnswerProvider struct {
	answer   answerService
	question questionService
}

// BulkUpsertUserAnswer upserts answers in bulk
func (provider AnswerProvider) BulkUpsertUserAnswer() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		score := 0
		total := 0
		as := []*model.Answer{}
		if err := json.NewDecoder(r.Body).Decode(&as); err != nil {
			log.Error(err)
			err = errors.New("cannot json decode the body, invalid json or missing field")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		validationErr := map[string]string{}
		toBeUpserted := []*model.Answer{}
		for i, a := range as {
			if err := model.ValidateAndPopulateAnswerRequest(a, provider.question); err != nil {
				validationErr["entry "+strconv.Itoa(i+1)] = errors.Wrap(err, "Invalid Answers").Error()
			} else {
				a.UserID = r.Context().Value("user_id").(string)
				a.ID = uuid.New().String()
				toBeUpserted = append(toBeUpserted, a)
			}
			score += a.Point
			total += a.PossiblePoint
		}
		if err := provider.answer.BulkUpsertUserAnswer(r.Context().Value("user_id").(string), toBeUpserted); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(validationErr) != 0 {
			json.NewEncoder(w).Encode(validationErr)
		}
		// TODO: Move this to model package
		respBody := struct {
			Score int
			Total int
		}{score, total}
		json.NewEncoder(w).Encode(respBody)
	})
}

// GetUserAnswers list all answers given by user
// TODO: Support pagination
func (provider AnswerProvider) GetUserAnswers() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if qs, err := provider.answer.ListAnswerByUser(r.Context().Value("user_id").(string)); err != nil {
			log.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			if qs == nil {
				qs = []*model.Answer{}
			}
			json.NewEncoder(w).Encode(qs)
		}
	})
}
