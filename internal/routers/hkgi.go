package routers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/hackagotchi/hkgi/database"
	"github.com/hackagotchi/hkgi/internal/handlers"
	"golang.org/x/crypto/bcrypt"
)

func SetupHkgiRoutes(hkgi fiber.Router) {
	log.Println("Initializing the HKGI routes")
	hkgi.Use(basicauth.New(basicauth.Config{
		Authorizer: func(username, password string) bool {
			db := database.DB

			var user string
			var pass string
			err := db.QueryRowx("SELECT username, password FROM stead WHERE username=$1", username).Scan(&user, &pass)
			log.Println("Got the user (or an error!)")

			if err != nil {
				log.Println("Error: %s", err)
				return false
			}

			err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
			log.Println("Compared password & hash!")

			if err == nil {
				return true
			}
			log.Println("Error: %s", err)

			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendStatus(403)
		},
	}))

	hkgi.Get("/getstead", handlers.GetStead)

	hkgi.Get("/activity", handlers.Activity)

	hkgi.Get("/users", handlers.Users)

	hkgi.Get("/manifest", handlers.Manifest)

	// POST
	hkgi.Post("/useitem", handlers.UseItem)

	hkgi.Post("/craft", handlers.Craft)

	hkgi.Post("/gib", handlers.Gib)

}
