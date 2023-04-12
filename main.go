package main

import (
	"fiber_rest_api/config"
	"fiber_rest_api/routes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	fmt.Println("Go sales api course started...")
	db.Connect()

	app := fiber.New()
	app.Use(cors.New())

	routes.Setup(app)

	app.Listen(":30001")
}
