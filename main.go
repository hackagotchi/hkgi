package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hackagotchi/hkgi/database"
	"github.com/hackagotchi/hkgi/internal"
	"github.com/hackagotchi/hkgi/internal/game"
	"github.com/hackagotchi/hkgi/internal/state"
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
			var tmp *fiber.Error
			e := fiber.Error{}
			if errors.As(err, &tmp) {
				code = tmp.Code
			} else {
				e.Code = code
			}

			e.Message = err.Error()

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

	go func() {
		err := game.RunTick()
		if err != nil {
			log.Fatal(err)
		}
	}()

	app.Listen(":6000")
}
