package handler

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// NewUserProvider returns a provider for User related operations.
func NewUserProvider(u userService) (up UserProvider) {
	return UserProvider{u}
}

// NewQuestionProvider returns a provider for Question related operations.
func NewQuestionProvider(q questionService) (qp QuestionProvider) {
	return QuestionProvider{q}
}

// NewAnswerProvider returns a provider for Answer related operations.
func NewAnswerProvider(a answerService, q questionService) (ap AnswerProvider) {
	return AnswerProvider{a, q}
}

// DefaultMiddleware is a middleware that is used by all endpoints
func DefaultMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		// TODO: May want to tighten this up
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		if r.Method == "OPTIONS" {
			//handle preflight, need to do better, should handle each case separately
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

// UserContextMiddleware adds the user-id to request
func UserContextMiddleware(u userService, isProd bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := ""
			cookie, _ := r.Cookie("user_session_id")
			if cookie != nil {
				sessionID = cookie.Value
			} else {
				// When developing there is this problem of CORS and we HTTPS, so get session from query instead
				if !isProd {
					sessionIDs := r.URL.Query()["user_session_id"]
					if len(sessionIDs) == 0 {
						http.Error(w, "User not authenticated", http.StatusUnauthorized)
						return
					}
					sessionID = sessionIDs[0]
				}
			}
			session, err := u.RefreshUserSession(sessionID)
			if err != nil {
				log.Error(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if session == nil {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}
			// TODO: Create a const type for context key
			ctx := context.WithValue(r.Context(), "user_id", session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NotFoundHandler is the catch all handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 - Not Found")
}
