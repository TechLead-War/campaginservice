package main

import (
	"campaign/internal/domain/models"
	"campaign/internal/infrastructure/cache"
	"campaign/internal/infrastructure/db"
	"campaign/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Setup database
	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	defer db.Close()

	// Setup cache
	memCache := setupCache(cfg)

	// Start metrics server
	metricsServer := startMetricsServer()

	// Initialize metrics
	utils.InitMetrics()

	// Setup main router
	router := setupRouter(db, memCache)

	// Create and start HTTP server
	server := createHTTPServer(cfg, router)
	go startServer(server, cfg.AppPort)

	// Wait for shutdown signal
	waitForShutdown(server, metricsServer)
}

// loadConfiguration loads and validates application configuration
func loadConfiguration() (*models.AppConfig, error) {
	cfg, err := models.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	// Set log level
	gin.SetMode(gin.ReleaseMode)
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	return cfg, nil
}

// setupDatabase establishes database connection and configures connection pool
func setupDatabase(cfg *models.AppConfig) (*sql.DB, error) {
	dbConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHOST, cfg.DBPORT, cfg.DBUSER, cfg.DBPass, cfg.DBName)

	d, err := db.Connect(dbConnString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Configure database connection pool
	d.SetMaxOpenConns(25)
	d.SetMaxIdleConns(25)
	d.SetConnMaxLifetime(5 * time.Minute)
	log.Println("Database connection pool configured successfully")

	return d, nil
}

// setupCache initializes the in-memory cache
func setupCache(cfg *models.AppConfig) *cache.MemoryCache {
	memCache := cache.NewMemoryCacheWithSize(cfg.CacheSize)
	log.Printf("Memory cache initialized with size: %d", cfg.CacheSize)
	return memCache
}

// setupRouter configures the main application router with middleware and routes
func setupRouter(db *sql.DB, memCache *cache.MemoryCache) *gin.Engine {
	router := gin.New()

	// Add middleware
	router.Use(utils.GinPrometheusMiddleware())
	router.Use(utils.RequestIDMiddleware())
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Setup routes
	baseRoute := "/api/v1"
	Delivery(router.Group(baseRoute), db, memCache)

	// Health check endpoint
	router.GET("/health", healthCheckHandler)

	return router
}

// healthCheckHandler handles health check requests
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

// createHTTPServer creates and configures the HTTP server
func createHTTPServer(cfg *models.AppConfig, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// startServer starts the HTTP server in a goroutine
func startServer(server *http.Server, port string) {
	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}

// startMetricsServer starts the Prometheus metrics server
func startMetricsServer() *http.Server {
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

	server := &http.Server{
		Addr:    ":9090",
		Handler: metricsRouter,
	}

	go func() {
		log.Println("Metrics server starting on :9090")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting metrics server: %v", err)
		}
	}()

	return server
}

// waitForShutdown waits for interrupt signal and gracefully shuts down servers
func waitForShutdown(server *http.Server, metricsServer *http.Server) {
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Shutdown metrics server
	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down metrics server: %v", err)
	}

	log.Println("Server exited")
}
