package internal

import (
	"git.sr.ht/~muirrum/hkgi/internal/handlers"
	"git.sr.ht/~muirrum/hkgi/internal/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	app.Post("/signup", handlers.Signup)
	hkgi := app.Group("/hkgi", logger.New())
	routers.SetupHkgiRoutes(hkgi)
}
