package app

import (
	"fmt"
	"log"
	"net/http"

	"qubix/config"
	"qubix/internal/database"
	"qubix/internal/handler"
	"qubix/internal/repository"
	"qubix/internal/service"
)

func Run() error {
	cfg := config.LoadConfig()

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
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
	answerService := service.NewAnswerService(attemptRepository, questionRepository, optionRepository, answerRepository)
	resultService := service.NewResultService(attemptRepository, answerRepository, questionRepository, optionRepository, resultRepository)

	// Handlers
	userHandler := handler.NewUserHandler(userService, jwtService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	quizHandler := handler.NewQuizHandler(quizService)
	questionHandler := handler.NewQuestionHandler(questionService, quizService)
	optionHandler := handler.NewOptionHandler(optionService, questionService, quizService)
	attemptHandler := handler.NewAttemptHandler(attemptService, answerService, resultService)

	// Router
	router := handler.NewRouter(
		userHandler,
		categoryHandler,
		quizHandler,
		questionHandler,
		optionHandler,
		attemptHandler,
		jwtService,
	)

	log.Println("Qubix server started at http://localhost:" + cfg.Port)

	return http.ListenAndServe(":"+cfg.Port, corsMiddleware(router))
}

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
