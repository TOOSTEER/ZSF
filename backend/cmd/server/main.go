package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"quest-backend/internal/auth"
	"quest-backend/internal/crossword"
	"quest-backend/internal/db"
	"quest-backend/internal/progress"
	"quest-backend/internal/tasks"
	"quest-backend/internal/teams"
	"quest-backend/internal/users"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	port := getEnv("PORT", "8080")
	dsn := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/questdb?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "change_me_in_production")
	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		jwtExpiryHours = 24
	}
	corsOrigin := getEnv("CORS_ORIGIN", "http://localhost:5500")

	pool, err := db.NewPool(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	authCfg := auth.Config{Secret: jwtSecret, ExpiryHours: jwtExpiryHours}

	usersRepo := users.NewRepository(pool)
	teamsRepo := teams.NewRepository(pool)
	tasksRepo := tasks.NewRepository(pool)
	progressRepo := progress.NewRepository(pool)
	crosswordRepo := crossword.NewRepository(pool)

	authHandler := auth.NewHandler(usersRepo, teamsRepo, authCfg)
	usersHandler := users.NewHandler(usersRepo, teamsRepo)
	teamsHandler := teams.NewHandler(teamsRepo)
	tasksHandler := tasks.NewHandler(tasksRepo, crosswordRepo)
	progressHandler := progress.NewHandler(progressRepo)
	crosswordHandler := crossword.NewHandler(crosswordRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{corsOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		// Публичные эндпоинты
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Защищённые эндпоинты (JWT)
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware(authCfg))

			r.Get("/me", usersHandler.Me)
			r.Get("/teams/{id}/progress", teamsHandler.GetProgress)
			r.Get("/tasks/{id}", tasksHandler.GetByID)
			r.Get("/tasks/{id}/crossword", tasksHandler.GetCrossword)

			// Только для host
			r.With(auth.RequireHost).Post("/progress/award", progressHandler.Award)

			// Только для member
			r.With(auth.RequireMember).Post("/crossword/{task_id}/check", crosswordHandler.Check)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":" + port
	log.Printf("quest-backend listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
