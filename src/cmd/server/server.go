package server

import (
	"context"
	"fmt"
	"log"
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
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type Server struct {
	app     *fiber.App
	cfg     *config.Config
	sqlDb   *gorm.DB
	noSqlDb *mongo.Database
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (srv *Server) InitDatabases() {
	// Initialize the SQL database
	sqlDb, err := relational.NewMySQLDB(srv.cfg.SQLDb.DSN, srv.cfg.SQLDb.Clean)
	if err != nil {
		logger.Fatal("Failed to initialize MySQL database: %v", err)
	}
	srv.sqlDb = sqlDb

	// Check if data initialization is required
	shouldInitialize, err := relational.ShouldInitializeData(sqlDb)
	if err != nil {
		logger.Fatal("Failed to check data initialization status: %v", err)
	}

	if shouldInitialize {
		logger.Info("Running data initialization: adding default data to the database.")
		// For testing purposes, this SQL script is executed to populate the database with sample data.
		err = relational.ExecuteSQLFile(sqlDb, "./src/internal/storage/relational/insert.sql")
		if err != nil {
			log.Fatalf("Failed to execute SQL file: %v", err)
		}
	} else {
		logger.Info("Data initialization script not required.")
	}

	// Initialize the MongoDB database
	mongoDb, err := nonrelational.NewMongoDB(srv.cfg.NoSQLDb.URI, srv.cfg.NoSQLDb.Database, srv.cfg.NoSQLDb.Clean)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB database: %v", err)
	}

	srv.noSqlDb = mongoDb
}

func (srv *Server) Run() error {

	// Init Logger
	logger.InitLogger(srv.cfg.IsProduction, srv.cfg.LogPath)

	// Init databases
	srv.InitDatabases()

	// Initialize the Fiber app
	srv.initFiber()

	// Start the server in a goroutine
	go func() {
		if err := srv.app.Listen(":" + srv.cfg.App.Port); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()
	logger.Info("Server is running on port %s", srv.cfg.App.Port)

	// Wait for interrupt signal for graceful shutdown
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
	logger.Sync() // Flush the logs, ensuring that all logs are written to the file before exit application
	return nil
}

func (srv *Server) initFiber() {
	// Create a new Fiber app
	srv.app = fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	// Routes
	srv.setupRoutes()
}

func (srv *Server) setupRoutes() {
	// Swagger documentation route
	srv.app.Get("/swagger/*", swagger.HandlerDefault)

	srv.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Payment Registration System!")
	})

	// Initialize services, and handlers
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
