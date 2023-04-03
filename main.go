package main

import (
	"belajar/database"
	"belajar/database/migration"
	"belajar/route"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// INITIAL DATABASE
	database.DatabaeInit()

	// RUN MIGRATION
	migration.RunMigration()

	// INITIAL ROUTE
	route.RouteInit(app)

	app.Listen(":8080")
}
