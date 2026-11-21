package main

import (
	"log"
	"net/http"

	"qubix/config"
	"qubix/internal/database"
	"qubix/internal/handler"
	"qubix/internal/repository"
	"qubix/internal/service"

	appmiddleware "qubix/middleware"
)

func main() {
	cfg := config.LoadConfig()

	if err := cfg.Validate(); err != nil {
		log.Fatal("Configuration error:", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Repositories
	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	quizRepository := repository.NewQuizRepository(db)
	questionRepository := repository.NewQuestionRepository(db)
	optionRepository := repository.NewOptionRepository(db)
	attemptRepository := repository.NewAttemptRepository(db)
	answerRepository := repository.NewAnswerRepository(db)
	resultRepository := repository.NewResultRepository(db)

	// Services
	jwtService := service.NewJWTService(cfg.JWTSecret)
	userService := service.NewUserService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	quizService := service.NewQuizService(quizRepository)
	questionService := service.NewQuestionService(questionRepository)
	optionService := service.NewOptionService(optionRepository)
	attemptService := service.NewAttemptService(attemptRepository, quizRepository)

	answerService := service.NewAnswerService(
		attemptRepository,
		questionRepository,
		optionRepository,
		answerRepository,
	)

	resultService := service.NewResultService(
		attemptRepository,
		answerRepository,
		questionRepository,
		optionRepository,
		resultRepository,
	)

	// Handlers
	userHandler := handler.NewUserHandler(userService, jwtService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	quizHandler := handler.NewQuizHandler(quizService)
	questionHandler := handler.NewQuestionHandler(questionService, quizService)
	optionHandler := handler.NewOptionHandler(
		optionService,
		questionService,
		quizService,
	)
	attemptHandler := handler.NewAttemptHandler(
		attemptService,
		answerService,
		resultService,
	)

	// Router
	mux := http.NewServeMux()

	// Authentication routes
	mux.HandleFunc("POST /sign-up", userHandler.SignUp)
	mux.HandleFunc("POST /login", userHandler.Login)

	// Private route
	mux.Handle(
		"GET /private",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Access granted!"))
			}),
		),
	)

	// Category routes
	mux.Handle(
		"POST /categories",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(categoryHandler.Create),
		),
	)
	mux.HandleFunc("GET /categories", categoryHandler.GetAll)
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetByID)
	mux.Handle(
		"PUT /categories/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(categoryHandler.Update),
		),
	)
	mux.Handle(
		"DELETE /categories/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(categoryHandler.Delete),
		),
	)

	// Quiz routes
	mux.Handle(
		"POST /quizzes",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(quizHandler.Create),
		),
	)
	mux.HandleFunc("GET /quizzes", quizHandler.GetAll)
	mux.HandleFunc("GET /quizzes/{id}", quizHandler.GetByID)
	mux.Handle(
		"PUT /quizzes/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(quizHandler.Update),
		),
	)
	mux.Handle(
		"DELETE /quizzes/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(quizHandler.Delete),
		),
	)

	// Question routes
	mux.Handle(
		"POST /quizzes/{quiz_id}/questions",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(questionHandler.Create),
		),
	)
	mux.HandleFunc("GET /quizzes/{quiz_id}/questions", questionHandler.GetByQuizID)
	mux.HandleFunc("GET /questions/{id}", questionHandler.GetByID)
	mux.Handle(
		"PUT /questions/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(questionHandler.Update),
		),
	)
	mux.Handle(
		"DELETE /questions/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(questionHandler.Delete),
		),
	)

	// Option routes
	mux.Handle(
		"POST /questions/{question_id}/options",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(optionHandler.Create),
		),
	)
	mux.HandleFunc(
		"GET /questions/{question_id}/options",
		optionHandler.GetByQuestionID,
	)
	mux.HandleFunc(
		"GET /options/{id}",
		optionHandler.GetByID,
	)
	mux.Handle(
		"PUT /options/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(optionHandler.Update),
		),
	)
	mux.Handle(
		"DELETE /options/{id}",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(optionHandler.Delete),
		),
	)

	// Attempt routes
	mux.Handle(
		"POST /quizzes/{quiz_id}/attempts",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(attemptHandler.Start),
		),
	)
	mux.Handle(
		"GET /attempts/{id}/result",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(attemptHandler.GetResult),
		),
	)
	mux.Handle(
		"POST /attempts/{id}/answers",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(attemptHandler.SubmitAnswer),
		),
	)
	mux.Handle(
		"POST /attempts/{id}/finish",
		appmiddleware.AuthMiddleware(
			jwtService,
			http.HandlerFunc(attemptHandler.Finish),
		),
	)

	// Start server
	log.Println("Qubix server started at http://localhost:" + cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, corsMiddleware(mux)); err != nil {
		log.Fatal("Server error:", err)
	}
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "https://editor.swagger.io" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			w.Header().Set("Access-Control-Max-Age", "600")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
