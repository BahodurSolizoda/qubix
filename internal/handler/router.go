package handler

import (
	"net/http"

	appmiddleware "qubix/internal/middleware"
	"qubix/internal/service"
)

func NewRouter(userHandler *UserHandler, categoryHandler *CategoryHandler, quizHandler *QuizHandler, questionHandler *QuestionHandler, optionHandler *OptionHandler, attemptHandler *AttemptHandler, jwtService *service.JWTService) http.Handler {
	mux := http.NewServeMux()

	// Authentication routes
	mux.HandleFunc("POST /sign-up", userHandler.SignUp)
	mux.HandleFunc("POST /login", userHandler.Login)

	// Private route
	mux.Handle("GET /private", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access granted!"))
	})))

	// Category routes
	mux.Handle("POST /categories", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(categoryHandler.Create)))
	mux.HandleFunc("GET /categories", categoryHandler.GetAll)
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetByID)
	mux.Handle("PUT /categories/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(categoryHandler.Update)))
	mux.Handle("DELETE /categories/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(categoryHandler.Delete)))

	// Quiz routes
	mux.Handle("POST /quizzes", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(quizHandler.Create)))
	mux.HandleFunc("GET /quizzes", quizHandler.GetAll)
	mux.HandleFunc("GET /quizzes/{id}", quizHandler.GetByID)
	mux.Handle("PUT /quizzes/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(quizHandler.Update)))
	mux.Handle("DELETE /quizzes/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(quizHandler.Delete)))

	// Question routes
	mux.Handle("POST /quizzes/{quiz_id}/questions", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(questionHandler.Create)))
	mux.HandleFunc("GET /quizzes/{quiz_id}/questions", questionHandler.GetByQuizID)
	mux.HandleFunc("GET /questions/{id}", questionHandler.GetByID)
	mux.Handle("PUT /questions/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(questionHandler.Update)))
	mux.Handle("DELETE /questions/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(questionHandler.Delete)))

	// Option routes
	mux.Handle("POST /questions/{question_id}/options", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(optionHandler.Create)))
	mux.HandleFunc("GET /questions/{question_id}/options", optionHandler.GetByQuestionID)
	mux.HandleFunc("GET /options/{id}", optionHandler.GetByID)
	mux.Handle("PUT /options/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(optionHandler.Update)))
	mux.Handle("DELETE /options/{id}", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(optionHandler.Delete)))

	// Attempt routes
	mux.Handle("POST /quizzes/{quiz_id}/attempts", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(attemptHandler.Start)))
	mux.Handle("GET /attempts/{id}/result", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(attemptHandler.GetResult)))
	mux.Handle("POST /attempts/{id}/answers", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(attemptHandler.SubmitAnswer)))
	mux.Handle("POST /attempts/{id}/finish", appmiddleware.AuthMiddleware(jwtService, http.HandlerFunc(attemptHandler.Finish)))

	return mux
}
