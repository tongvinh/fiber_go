package routes

import (
	"fiber_rest_api/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/cashiers/:cashierId/login", controller.Login)
	app.Get("/cashiers/:cashierId/logout", controller.Logout)
	app.Post("/cashiers/cashierId/passcode", controller.Passcode)

	//cashier routes
	app.Post("/cashiers", controller.CreateCashier)
	app.Get("/cashiers", controller.CashiersList)
	app.Get("/cashiers/:cashierId", controller.GetCashierDetail)
	app.Delete("/cashiers/:cashierId", controller.DeleteCashier)
	app.Put("/cashiers/:cashierId", controller.UpdateCashier)
}
