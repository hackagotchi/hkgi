package routers

import (
	"context"
	"log"

	"git.sr.ht/~muirrum/hkgi/database"
	"git.sr.ht/~muirrum/hkgi/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"golang.org/x/crypto/bcrypt"
)

func SetupHkgiRoutes(hkgi fiber.Router) {
	log.Println("Initializing the HKGI routes")
	hkgi.Use(basicauth.New(basicauth.Config{
		Authorizer: func(username, password string) bool {
			db, err := database.DB.Acquire(context.Background())

			if err != nil {
				log.Fatal(err.Error())
			}

			var user string
			var pass string
			err = db.QueryRow(context.Background(), "SELECT username, password FROM stead WHERE username=$1", username).Scan(&user, &pass)
			log.Println("Got the user (or an error!)")

			if err != nil {
				log.Fatal(err.Error())
				return false
			}

			err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
			log.Println("Compared password & hash!")

			if err == nil {
				return true
			}

			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendStatus(403)
		},
	}))

	hkgi.Get("/getstead", handlers.GetStead)

	hkgi.Get("/activity", func(c *fiber.Ctx) error { return nil })

	hkgi.Get("/users", func(c *fiber.Ctx) error { return nil })

	hkgi.Get("/manifest", func(c *fiber.Ctx) error { return nil })

	// POST
	hkgi.Post("/useitem", func(c *fiber.Ctx) error { return nil })

	hkgi.Post("/craft", func(c *fiber.Ctx) error { return nil })

}
