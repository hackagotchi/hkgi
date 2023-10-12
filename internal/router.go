package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hackagotchi/hkgi/internal/handlers"
	"github.com/hackagotchi/hkgi/internal/routers"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	app.Post("/signup", handlers.Signup)
	hkgi := app.Group("/hkgi", logger.New())
	routers.SetupHkgiRoutes(hkgi)
}
