/*
 * Payment Registration System - Server Configuration
 * --------------------------------------------------
 * This file defines the core server logic, including:
 * - Database initialization (SQL & NoSQL)
 * - Fiber-based HTTP server setup
 * - Route configuration for API endpoints
 * - Graceful shutdown handling
 *
 * Created: Dec. 11, 2024
 * License: GNU General Public License v3.0
 */

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GabrielEValenzuela/Payment-Registration-System/src/cmd/handlers"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/config"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/services"
	nonrelational "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/non_relational"
	non_relational_repository "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/non_relational/repository"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/relational"
	relational_repository "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/storage/relational/repository"
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

/*
 * Server
 * --------------------------------------------------
 * Represents the main application server, including:
 * - Fiber HTTP server
 * - Configuration settings
 * - SQL (MySQL) and NoSQL (MongoDB) database connections
 */
type Server struct {
	app     *fiber.App
	cfg     *config.Config
	sqlDb   *gorm.DB
	noSqlDb *mongo.Database
}

/*
 * NewServer
 * --------------------------------------------------
 * Creates a new server instance with the provided configuration.
 *
 * Params:
 * - cfg (*config.Config): Configuration instance containing app settings.
 *
 * Returns:
 * - *Server: A new server instance.
 */
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

/*
 * InitDatabases
 * --------------------------------------------------
 * Initializes both SQL (MySQL) and NoSQL (MongoDB) databases.
 * Also runs data initialization if required.
 */
func (srv *Server) InitDatabases() {
	// Initial delay before retrying
	delay := 2 * time.Second
	maxDelay := 60 * time.Second

	for {
		// Attempt to connect to MySQL
		sqlDb, err := relational.NewMySQLDB(srv.cfg.SQLDb.DSN, srv.cfg.SQLDb.Clean)
		if err != nil {
			logger.Warn("Failed to initialize MySQL database: %v. Retrying in %v...", err, delay)
			time.Sleep(delay)
			if delay < maxDelay {
				delay *= 2 // Exponential backoff
			}
			continue
		}

		srv.sqlDb = sqlDb
		logger.Info("Successfully connected to MySQL database")

		// Attempt to connect to MongoDB
		mongoDb, err := nonrelational.NewMongoDB(srv.cfg.NoSQLDb.URI, srv.cfg.NoSQLDb.Database, srv.cfg.NoSQLDb.Clean)
		if err != nil {
			logger.Warn("Failed to initialize MongoDB database: %v. Retrying in %v...", err, delay)
			time.Sleep(delay)
			if delay < maxDelay {
				delay *= 2
			}
			continue
		}

		srv.noSqlDb = mongoDb
		logger.Info("Successfully connected to MongoDB database")

		// Both databases are successfully connected, break the loop
		break
	}
}

/*
 * Run
 * --------------------------------------------------
 * Starts the Fiber HTTP server, initializes databases, and handles graceful shutdown.
 *
 * Returns:
 * - error: If server shutdown fails.
 */
func (srv *Server) Run() error {
	logger.InitLogger(srv.cfg.IsProduction, srv.cfg.LogPath)
	srv.InitDatabases()
	srv.initFiber()

	// Start the server asynchronously
	go func() {
		if err := srv.app.Listen(":" + srv.cfg.App.Port); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()
	logger.Info("Server is running on port %s", srv.cfg.App.Port)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.app.Shutdown(); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("Server exited properly")
	logger.Sync() // Ensure all logs are flushed before exiting
	return nil
}

/*
 * initFiber
 * --------------------------------------------------
 * Initializes the Fiber app with default configurations and routes.
 */
func (srv *Server) initFiber() {
	srv.app = fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorHandler: srv.errorHandler, // Set global error handler
	})

	// Apply rate limiting middleware
	srv.app.Use(limiter.New(limiter.Config{
		Max:        100,             // Allow 100 requests per window
		Expiration: 1 * time.Minute, // Reset every minute
	}))

	// Enable CORS for all routes
	srv.app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // ToDo: Change to production domains
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	srv.setupRoutes()
}

/*
 * setupRoutes
 * --------------------------------------------------
 * Configures API routes for both SQL and NoSQL services.
 */
func (srv *Server) setupRoutes() {
	srv.app.Get("/swagger/*", swagger.HandlerDefault)

	srv.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Payment Registration System!")
	})

	// Initialize handlers
	bankHandlerRelational := handlers.NewBankHandler(services.NewBankService(relational_repository.NewBankRelationalRepository(srv.sqlDb)))
	bankHandlerNonRelational := handlers.NewBankHandler(services.NewBankService(non_relational_repository.NewBankNonRelationalRepository(srv.noSqlDb)))

	cardHandlerRelational := handlers.NewCardHandler(services.NewCardService(relational_repository.NewCardRelationalRepository(srv.sqlDb)))
	cardHandlerNonRelational := handlers.NewCardHandler(services.NewCardService(non_relational_repository.NewCardNonRelationalRepository(srv.noSqlDb)))

	promotionHandlerRelation := handlers.NewPromotionHandler(services.NewPromotionService(relational_repository.NewPromotionRelationRepository(srv.sqlDb)))
	promotionHandlerNonRelation := handlers.NewPromotionHandler(services.NewPromotionService(non_relational_repository.NewPromotionNonRelationalRepository(srv.noSqlDb)))

	storeHandlerRelation := handlers.NewStoreHandler(services.NewStoreService(relational_repository.NewStoreRelationalRepository(srv.sqlDb)))
	storeHandlerNonRelation := handlers.NewStoreHandler(services.NewStoreService(non_relational_repository.NewStoreNonRelationalRepository(srv.noSqlDb)))

	// API version group
	apiGroup := srv.app.Group("/v1")

	// SQL routes group
	sqlGroup := apiGroup.Group("/sql")
	// NoSQL routes group
	mongoGroup := apiGroup.Group("/no-sql")

	// -- Bank Routes --
	sqlGroup.Post("/promotions/add-promotion", bankHandlerRelational.AddFinancingPromotionToBank())
	sqlGroup.Patch("/promotions/financing/:code", bankHandlerRelational.ExtendFinancingPromotionValidity())
	sqlGroup.Delete("/promotions/financing/:code", bankHandlerRelational.DeleteFinancingPromotion())
	sqlGroup.Patch("/promotions/discount/:code", bankHandlerRelational.ExtendDiscountPromotionValidity())
	sqlGroup.Delete("/promotions/discount/:code", bankHandlerRelational.DeleteDiscountPromotion())
	sqlGroup.Get("/banks/customers/count", bankHandlerRelational.GetBankCustomerCounts())

	mongoGroup.Post("/promotions/add-promotion", bankHandlerNonRelational.AddFinancingPromotionToBank())
	mongoGroup.Patch("/promotions/financing/:code", bankHandlerNonRelational.ExtendFinancingPromotionValidity())
	mongoGroup.Delete("/promotions/financing/:code", bankHandlerNonRelational.DeleteFinancingPromotion())
	mongoGroup.Patch("/promotions/discount/:code", bankHandlerNonRelational.ExtendDiscountPromotionValidity())
	mongoGroup.Delete("/promotions/discount/:code", bankHandlerNonRelational.DeleteDiscountPromotion())
	mongoGroup.Get("/banks/customers/count", bankHandlerNonRelational.GetBankCustomerCounts())

	// -- Card Routes --
	sqlGroup.Get("/cards/summary/:cardNumber/:month/:year", cardHandlerRelational.GetPaymentSummary())
	sqlGroup.Get("/cards/expiring/:day/:month/:year", cardHandlerRelational.GetCardsExpiringInNext30Days())
	sqlGroup.Get("/cards/purchase/monthly/:cuit/:finalAmount/:paymentVoucher", cardHandlerRelational.GetPurchaseMonthly())
	sqlGroup.Get("/cards/top", cardHandlerRelational.GetTop10CardsByPurchases())

	mongoGroup.Get("/cards/summary/:cardNumber/:month/:year", cardHandlerNonRelational.GetPaymentSummary())
	mongoGroup.Get("/cards/expiring/:day/:month/:year", cardHandlerNonRelational.GetCardsExpiringInNext30Days())
	mongoGroup.Get("/cards/purchase/monthly/:cuit/:finalAmount/:paymentVoucher", cardHandlerNonRelational.GetPurchaseMonthly())
	mongoGroup.Get("/cards/top", cardHandlerNonRelational.GetTop10CardsByPurchases())

	// -- Promotion Routes --
	sqlGroup.Get("/promotions/:cuit/:startDate/:endDate", promotionHandlerRelation.GetAvailablePromotionsByStoreAndDateRange())
	sqlGroup.Get("/promotions/most-used", promotionHandlerRelation.GetMostUsedPromotion())

	mongoGroup.Get("/promotions/:cuit/:startDate/:endDate", promotionHandlerNonRelation.GetAvailablePromotionsByStoreAndDateRange())
	mongoGroup.Get("/promotions/most-used", promotionHandlerNonRelation.GetMostUsedPromotion())

	// -- Store Routes --
	sqlGroup.Get("/stores/highest-revenue/:month/:year", storeHandlerRelation.GetStoreWithHighestRevenueByMonth())
	mongoGroup.Get("/stores/highest-revenue/:month/:year", storeHandlerNonRelation.GetStoreWithHighestRevenueByMonth())
}

/*
 * errorHandler
 * --------------------------------------------------
 * Handles errors returned by the API handlers.
 */
func (srv *Server) errorHandler(ctx *fiber.Ctx, err error) error {
	// Default HTTP status
	code := fiber.StatusInternalServerError

	// Fiber error type handling
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log the error (ensure logger is configured)
	logger.Error("API Error: %v", err)

	// Return JSON error response
	return ctx.Status(code).JSON(fiber.Map{
		"error":   err.Error(),
		"code":    code,
		"message": "An error occurred while processing the request",
	})
}
