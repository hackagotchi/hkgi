package routers

import "github.com/gofiber/fiber/v2"

func SetupHkgiRoutes(router fiber.Router) {
	hkgi := router.Group("/hkgi")

	hkgi.Get("/getstead", func(c *fiber.Ctx) error {})

	hkgi.Get("/activity", func(c *fiber.Ctx) error {})

	hkgi.Get("/users", func(c *fiber.Ctx) error {})

	hkgi.Get("/manifest", func(c *fiber.Ctx) error {})

	// POST
	hkgi.Post("/useitem", func(c *fiber.Ctx) error {})

	hkgi.Post("/craft", func(c *fiber.Ctx) error {})

}
