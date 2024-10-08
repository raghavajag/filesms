package main

import (
	"context"
	"database/sql"
	"filesms/internal/core/services/authsrv"
	"filesms/internal/core/services/cleanupservice"
	"filesms/internal/core/services/filesrv"
	"filesms/internal/handlers/authhdl"
	"filesms/internal/handlers/filehdl"
	"filesms/internal/repositories/filerepo"
	"filesms/internal/repositories/userrepo"

	response "filesms/pkg/api"
	redisStore "filesms/pkg/cache/redis"
	"filesms/pkg/jwt"
	"filesms/pkg/middleware"
	"filesms/pkg/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	/*
		// if running locally load env variables from .env file
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Error loading .env file")
			}
	*/
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run migrations
	/*
		err = database.RunMigrations(db)
		if err != nil {
			log.Fatalf("Could not run migrations: %v", err)
		}
	*/

	// Initialize Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Connected to Redis...")
	// Create Redis cache
	redisCache := redisStore.NewRedisCache(redisClient)

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
	baseURL := "http://api:8080/files"
	fileService := filesrv.NewFileService(fileRepo, localStorage, baseURL, redisCache)

	// Initialize and start cleanup service, (10 seconds for testing)
	cleanupService := cleanupservice.NewCleanupService(fileRepo, "./tmp", 10*time.Second)
	go func() {
		cleanupService.Start(context.Background())
	}()

	// Initialize handlers
	authHandler := authhdl.NewAuthHandler(authService)

	// Initialize handlers
	fileHandler := filehdl.NewFileHandler(fileService)
	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json := map[string]string{"by": "21BCT0114 Raghava Jagarwal", "postman": "https://documenter.getpostman.com/view/11141903/2sAXqp83uR"}
		response.Success(w, "Health check", json)
	})

	// Define routes
	router.HandleFunc("/register", middleware.ErrorHandler(authHandler.Register))
	router.HandleFunc("/login", middleware.ErrorHandler(authHandler.Login))

	// Define Protedted routes
	router.HandleFunc("/me", middleware.AuthMiddleware(middleware.ErrorHandler(authHandler.Me)))
	router.HandleFunc("/upload", middleware.AuthMiddleware(middleware.ErrorHandler(fileHandler.Upload)))
	router.HandleFunc("/files", middleware.AuthMiddleware(middleware.ErrorHandler(fileHandler.GetFiles)))
	router.HandleFunc("/share", middleware.AuthMiddleware(middleware.ErrorHandler(fileHandler.ShareFile)))
	router.HandleFunc("/files/search", middleware.AuthMiddleware(middleware.ErrorHandler(fileHandler.SearchFiles)))
	router.HandleFunc("/file", middleware.AuthMiddleware(middleware.ErrorHandler(fileHandler.GetFile)))

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
