package main

// This package is the entry point for the server application.

//

// We need to setup the server to listen on a port and handle requests. We also need to configure the server to use the customer handler.
// The server will be responsible for handling the requests and responses from the client.

import (
	"github.com/GabrielEValenzuela/Payment-Registration-System/src/cmd/server/handlers"
	_ "github.com/GabrielEValenzuela/Payment-Registration-System/src/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Server is the main struct for the server application.
type Server struct {
	app             *fiber.App
	customerHandler *handlers.CustomerHandler
}

func NewServer(customerHandler *handlers.CustomerHandler) *Server {
	return &Server{
		app:             fiber.New(),
		customerHandler: customerHandler,
	}
}

// ConfigureServer configures the server to listen on a port and handle requests.
func (srv *Server) ConfigureServer() {
	// Initialize the server
	srv.app = fiber.New()

	// Define the routes for the server
	srv.ConfigureRoutes()
}

func (srv *Server) ConfigureRoutes() {
	// Define the route for the documentation
	srv.app.Use("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/doc.json",
	}))

	// Define the routes for the server
	sqlRouting := srv.app.Group("/sql")
	mongoRouting := srv.app.Group("/mongo")

	sqlRouting.Get("/customer/:id", srv.customerHandler.GetCustomerById())
	mongoRouting.Get("/customer/:id", srv.customerHandler.GetCustomerById())
}

func ConfigureServer() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, client 👋!")
	})

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
