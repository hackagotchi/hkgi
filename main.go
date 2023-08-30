package main

import (
	"log"

	"git.devcara.com/hkgi/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	database.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber!"})
	})

	app.Listen(":6000")
}
