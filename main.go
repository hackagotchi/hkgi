package main

import (
	"errors"
	"log"

	"git.sr.ht/~muirrum/hkgi/database"
	"git.sr.ht/~muirrum/hkgi/internal"
	"git.sr.ht/~muirrum/hkgi/internal/game"
	"git.sr.ht/~muirrum/hkgi/internal/state"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	state.InitState()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			return c.Status(code).SendString(e.Message)

			// Return from handler
			return nil
		},
	})

	internal.SetupRoutes(app)

	database.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber!"})
	})

	go game.RunTick()

	app.Listen(":6000")
}
