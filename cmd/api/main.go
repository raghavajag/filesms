package main

import (
	"context"
	"database/sql"
	database "filesms/db"
	"filesms/internal/core/services/authsrv"
	"filesms/internal/core/services/filesrv"
	"filesms/internal/handlers/authhdl"
	"filesms/internal/handlers/filehdl"
	"filesms/internal/repositories/filerepo"
	"filesms/internal/repositories/userrepo"
	"filesms/pkg/jwt"
	"filesms/pkg/middleware"
	"filesms/pkg/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/dgrijalva/jwt-go"

	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run migrations
	err = database.RunMigrations(db)
	if err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := userrepo.NewPostgresUserRepository(db)
	fileRepo := filerepo.NewPostgresFileRepository(db)

	// Create JWT maker
	jwtMaker := jwt.NewJWTMaker(os.Getenv("JWT_SECRET"))

	// Initialize local storage
	storagePath := os.Getenv("STORAGE_PATH")
	localStorage, err := storage.NewLocalStorage(storagePath)
	if err != nil {
		log.Fatalf("Error initializing local storage: %v", err)
	}

	// Initialize services
	authService := authsrv.NewAuthService(userRepo, jwtMaker)
	baseURL := "http://localhost:8080/files"
	fileService := filesrv.NewFileService(fileRepo, localStorage, baseURL)

	// Initialize handlers
	authHandler := authhdl.NewAuthHandler(authService)

	// Initialize handlers
	fileHandler := filehdl.NewFileHandler(fileService)
	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Define routes
	router.HandleFunc("/register", middleware.ErrorHandler(authHandler.Register))
	router.HandleFunc("/login", middleware.ErrorHandler(authHandler.Login))

	// Define Protedted routes
	router.HandleFunc("/me", middleware.AuthMiddleware(middleware.ErrorHandler(authHandler.Me)))
	router.HandleFunc("/upload", middleware.AuthMiddleware(middleware.ErrorHandler(fileHandler.Upload)))
	router.HandleFunc("/files", middleware.AuthMiddleware(middleware.ErrorHandler(fileHandler.GetFiles)))

	// Define routes
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
