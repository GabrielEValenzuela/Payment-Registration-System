package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/cmd/handlers"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/config"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/customer"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

type Server struct {
	app   *fiber.App
	cfg   *config.Config
	sqlDb *sql.DB
	// mongoDb *mongo.Client // Uncomment this line when we made the MongoDB connection
	// Add other fields as needed (e.g., MongoDB client)
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (srv *Server) Run() error {
	// Initialize the database connections
	if err := srv.initSqlDatabase(); err != nil {
		return fmt.Errorf("failed to initialize databases: %w", err)
	}
	defer srv.sqlDb.Close() // Close the MySQL connection when the server exits, defer ensures this runs even if an error occurs

	// Initialize the Fiber app
	srv.initFiber()

	// Start the server in a goroutine
	go func() {
		if err := srv.app.Listen(":" + srv.cfg.App.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Infof("Server is running on port %s", srv.cfg.App.Port)

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	_, cancel := context.WithTimeout(context.Background(), srv.cfg.App.GracefulShutdown*time.Second)
	defer cancel()

	if err := srv.app.Shutdown(); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info("Server exited properly")
	return nil
}

func (srv *Server) initSqlDatabase() error {
	// Initialize MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		srv.cfg.MySQL.User,
		srv.cfg.MySQL.Password,
		srv.cfg.MySQL.Host,
		srv.cfg.MySQL.Port,
		srv.cfg.MySQL.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	db.SetMaxOpenConns(srv.cfg.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(srv.cfg.MySQL.MaxIdleConns)
	db.SetConnMaxLifetime(srv.cfg.MySQL.ConnMaxLifetime * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	srv.sqlDb = db
	return nil
}

func (srv *Server) initFiber() {
	// Create a new Fiber app
	srv.app = fiber.New(fiber.Config{
		ReadTimeout:  srv.cfg.App.ReadTimeout * time.Second,
		WriteTimeout: srv.cfg.App.WriteTimeout * time.Second,
	})

	// Middleware
	srv.app.Use(logger.New())

	// Routes
	srv.setupRoutes()
}

func (srv *Server) setupRoutes() {
	// Swagger documentation route
	srv.app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize repositories, services, and handlers
	sqlStorage := storage.NewSqlStorage(srv.sqlDb)
	customerHandlerRelational := handlers.NewCustomerHandler(customer.NewCustomerService(sqlStorage))
	customerHandlerNonRelational := handlers.NewCustomerHandler(customer.NewCustomerService(sqlStorage))

	// API version group
	apiGroup := srv.app.Group("/v1")

	// SQL routes
	sqlGroup := apiGroup.Group("/sql")
	sqlGroup.Get("/customer/:id", customerHandlerRelational.GetCustomerById())

	mongoGroup := apiGroup.Group("/no-sql")
	mongoGroup.Get("/customer/:id", customerHandlerNonRelational.GetCustomerById())
}
