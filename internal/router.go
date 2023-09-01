package internal

import (
	"context"
	"log"

	"git.devcara.com/hkgi/database"
	"git.devcara.com/hkgi/internal/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

func SetupRoutes(app *fiber.App) {
	app.Use(basicauth.New(basicauth.Config{
		Authorizer: func(username, password string) bool {
			db,err := database.DB.Acquire()

			if (err != nil) {
				log.Fatal(err.Error())
			}

			var user User
			err = db.QueryRow(context.Background(), "SELECT username, password FROM steads WHERE username=$1", username).Scan(&user)

			if (err != nil) {
				return false;
			}
			
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

			if (err == nil) {
				return true
			}

			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendStatus(403)
		},
	}))

	routers.SetupHkgiRoutes(app)
}
