package model

// QuestionSection for Quetions
type QuestionSection struct {
	ID   string `xorm:"id" json:"id"`
	Name string `xorm:"name" json:"name,omitempty" r-validate:"required"`
}

// ValidateQuestionSectionRequest validates the QuestionSection struct as the Question was constructed by the http request
// TODO: Need to better organize the code, maybe we can make the validation step more abstract.
func ValidateQuestionSectionRequest(qs QuestionSection) error {
	v := RequestValidator()
	return TranslateError(v.Struct(qs))
}
