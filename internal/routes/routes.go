// internal/routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shravanirajulu2004/go-user-api/internal/handler"
)

func SetupRoutes(app *fiber.App, userHandler handler.UserHandler) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// User routes
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Get("/users", userHandler.ListUsers)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
}