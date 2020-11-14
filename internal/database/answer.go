package database

import (
	"covid-19/internal/model"

	"github.com/pkg/errors"
)

// BulkUpsertUserAnswer upserts answers for a single user
// TODO: wrap all db actions in this function into a single transaction. REALLY IMPORTANT.
// 		 I regret not doing it in the first place, but deadline is a b**ch. This have serious
// 		 implications. Business logic will definitely break.
func (db Orm) BulkUpsertUserAnswer(uid string, as []*model.Answer) (err error) {
	if len(as) == 0 {
		return nil
	}
	answerMap := map[string]*model.Answer{}
	for _, a := range as {
		answerMap[a.QuestionID] = a
	}
	existingAnswers := []*model.Answer{}
	if err := db.Table("Answer").Where("user_id = ?", uid).Find(&existingAnswers); err != nil {
		return errors.WithMessage(err, "Error Getting Answer for user - Database Error")
	}
	for _, ea := range existingAnswers {
		if _, err := db.Table("Answer").Cols("choice", "point").Update(answerMap[ea.QuestionID], ea); err != nil {
			return errors.WithMessage(err, "Error Updating Answer - Database Error")
		}
		delete(answerMap, ea.QuestionID)
	}
	toInsert := []*model.Answer{}
	for _, a := range answerMap {
		toInsert = append(toInsert, a)
	}
	if len(toInsert) != 0 {
		if _, err := db.Table("Answer").Insert(toInsert); err != nil {
			return errors.WithMessage(err, "Error Creating Answer - Database Error")
		}
	}
	return nil
}

// ListAnswerByUser gets the user's answers to all questions he answered.
func (db Orm) ListAnswerByUser(uid string) (results []*model.Answer, err error) {
	err = db.Table("Answer").Where("user_id = ?", uid).Find(&results)
	return results, errors.WithMessage(err, "Error Getting Answers - Database Error")
}
