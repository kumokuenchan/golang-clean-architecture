package main

import (
	"log"
	"net/http"
	"os"

	httpHandler "clean-arch-sample/internal/delivery/http"
	"clean-arch-sample/internal/infrastructure/database"
	"clean-arch-sample/internal/presenter"
	"clean-arch-sample/internal/usecase"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "cleanarch_sample")

	db, err := database.NewMySQLConnection(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	userRepo := database.NewMySQLUserRepository(db)
	userPresenter := presenter.NewHTTPUserPresenter()
	userUsecase := usecase.NewUserUsecase(userRepo, userPresenter)

	router := httpHandler.NewRouter(userUsecase)
	mux := router.SetupRoutes()

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}