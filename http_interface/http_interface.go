package http_interface

import (
	"User-Management/database"
	"User-Management/middleware"
	"User-Management/router"
	"database/sql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Handler struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (a *Handler) Initialize() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbInstance := &database.Database{User: dbUser, Password: dbPassword, Host: dbHost, Name: dbName}
	var err error
	a.DB, err = dbInstance.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize Gin router with default middleware
	a.Router = gin.Default()

	// Apply custom logger middleware
	a.Router.Use(middleware.LoggerMiddleware())

	// Initialize routes
	router.InitializeRoutes(a.Router, a.DB)

	// Add Prometheus metrics route
	a.Router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func (a *Handler) Run() {
	// Start pprof server in a separate goroutine
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Start Gin server
	serverAddr := ":" + os.Getenv("DB_LOCALHOST")
	log.Printf("Starting server on %s", serverAddr)
	if err := a.Router.Run(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
