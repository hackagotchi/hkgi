package handlers

import (
	"context"
	"log"

	"git.sr.ht/~muirrum/hkgi/database"
	"git.sr.ht/~muirrum/hkgi/internal/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	u := new(models.User)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	db, err := database.DB.Acquire(context.Background())

	if err != nil {
		return err
	}

	// See if they already exist
	var u2 models.User
	err = db.QueryRow(context.Background(), "SELECT username FROM stead WHERE username=$1", u.Username).Scan(&u2)
	log.Println("Queried database for existing user")

	if err != nil {
		log.Printf(err.Error())
	}

	if err == nil {
		return fiber.NewError(400, "User already exists!")
	} else {
		// Let's make a new hackstead!
		starting_inventory := fiber.Map{
			"nest_egg": 1,
			"bbc_seed": 1,
			"hvv_seed": 1,
			"cyl_seed": 1,
		}
		// Create a stead!
		pw_hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
			return err
		}
		var steadId int
		log.Println("Creating new stead...")
		err = db.QueryRow(context.Background(), "INSERT INTO stead (username, password, inventory) VALUES ($1, $2, $3) RETURNING id", u.Username, pw_hash, starting_inventory).Scan(&steadId)
		log.Println("Created a new stead for user" + u.Username)

		// Create a plant first
		_, err = db.Exec(context.Background(), "INSERT INTO plant (stead_owner, kind, xp, xp_multiplier) VALUES ($1, 'dirt', 0, 0)", steadId)
		log.Printf("Created a new patch of dirt\n")

		return c.JSON(fiber.Map{
			"ok": true,
		})
	}

	return nil
}
