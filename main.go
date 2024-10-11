package main

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sistemakreditasi/backend-akreditasi/routes"
)

func main() {
	app := fiber.New()
	// Define a fiber handler for all requests
	app.All("/*", adaptor.HTTPHandlerFunc(routes.URL))

	port := ":8080"
	app.Listen(port)
}
