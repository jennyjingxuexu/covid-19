package handler

import "net/http"

// NewUserProvider returns a provider for User related operations.
func NewUserProvider(u userService) (up UserProvider) {
	return UserProvider{u}
}

// NewQuestionProvider returns a provider for Question related operations.
func NewQuestionProvider(q questionService) (qp QuestionProvider) {
	return QuestionProvider{q}
}

// DefaultMiddleware is a middleware that is used by all endpoints
func DefaultMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
