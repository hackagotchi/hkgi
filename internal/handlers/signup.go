package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hackagotchi/hkgi/database"
	"github.com/hackagotchi/hkgi/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type (
	ValidationError struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}
)

func (v XValidator) Validate(data interface{}) []ValidationError {
	validationErrors := []ValidationError{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ValidationError

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

// This is the validator instance
// for more information see: https://github.com/go-playground/validator
var validate = validator.New()

func Signup(c *fiber.Ctx) error {
	val := &XValidator{
		validator: validate,
	}
	u := new(models.User)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	if errs := val.Validate(u); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	db := database.DB

	// See if they already exist
	var u2 models.User
	err := db.QueryRowx("SELECT username FROM stead WHERE username=$1", u.Username).Scan(&u2)
	log.Println("Queried database for existing user")

	if err != nil {
		log.Printf(err.Error())
	}

	if err == nil {
		return fiber.NewError(400, "User already exists!")
	} else {
		// Let's make a new hackstead!
		starting_inventory := map[string]interface{}{
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
		var steadId int64
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Creating new stead...")
		inv, _ := json.Marshal(starting_inventory)
		statuses, _ := json.Marshal(map[string]interface{}{})
		log.Printf("Username: %s", u.Username)
		_, err = tx.Exec("INSERT INTO stead (username, password, inventory, ephemeral_statuses) VALUES ($1, $2, $3, $4)", u.Username, pw_hash, inv, statuses)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Println("Created a new stead for user" + u.Username)
		tx.Commit()

		tx, err = db.Begin()

		err = db.QueryRowx("SELECT id FROM stead WHERE username=$1", u.Username).Scan(&steadId)

		// Create a plant first
		_, err = tx.Exec("INSERT INTO plant (stead_owner, kind, xp, xp_multiplier) VALUES ($1, 'dirt', 0, 0)", steadId)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Created a new patch of dirt\n")
		tx.Commit()

		return c.JSON(fiber.Map{
			"ok": true,
		})
	}

	return nil
}
