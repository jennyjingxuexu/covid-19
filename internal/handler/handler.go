package handler

// NewUserProvider returns a provider for User related operations.
func NewUserProvider(u userService) (up UserProvider) {
	return UserProvider{u}
}

// NewQuestionProvider returns a provider for Question related operations.
func NewQuestionProvider(q questionService) (qp QuestionProvider) {
	return QuestionProvider{q}
}
