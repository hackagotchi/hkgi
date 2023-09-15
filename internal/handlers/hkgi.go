package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"git.sr.ht/~muirrum/hkgi/database"
	"git.sr.ht/~muirrum/hkgi/internal/game"
	"git.sr.ht/~muirrum/hkgi/internal/models"
	"git.sr.ht/~muirrum/hkgi/internal/state"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type SerializedPlant struct {
	Kind      pgtype.Text
	Lvl       int
	TtYield   float32
	YieldProg float32
	TtLevelUp float32
	LvlupProg float32
}

func lvlfromxp(xp int) int {
	// get log_1.3...
	log_13 := math.Log(1.3)
	lvl := int(math.Log(float64(xp)/10) / log_13)
	if lvl < 0 {
		lvl = 0
	}
	return lvl
}

func getLevelCost(lvl int) int {
	return int(math.Floor(10 * math.Pow(1.3, float64(lvl+1))))
}

func getXpRemaining(xp int) int {
	lvl := lvlfromxp(xp)

	cost := getLevelCost(lvl + 1)

	return int(math.Round(float64(cost) - float64(xp)))
}

func xpPerYield(xp int) int {
	level := lvlfromxp(xp)

	// magic numbers!!

	return int(math.Floor(900 * (1 - float64(level)/27)))
}

func GetStead(c *fiber.Ctx) error {
	db := database.DB

	user := c.Locals("username")

	var inventory json.RawMessage
	var ephemeral_statuses []uint8

	err := db.QueryRowx("SELECT inventory, (CASE WHEN ephemeral_statuses IS NULL THEN '{\"\"}' ELSE ephemeral_statuses END) as ephemeral_statuses FROM stead WHERE username=$1", user).Scan(&inventory, &ephemeral_statuses)

	if err != nil {
		return err
	}
	log.Println("Retrieved inventory and statuses!")

	// Now we have the inventory and any statuses, let's get the plants!

	var userId int
	err = db.QueryRowx("SELECT id FROM stead WHERE username=$1", user).Scan(&userId)
	if err != nil {
		return err
	}
	log.Println("Got the user id")
	rows, err := db.Query("SELECT kind, xp, xp_multiplier, (CASE WHEN next_yield IS NULL THEN '1970-01-01 00:00:00' ELSE next_yield END) as next_yield FROM plant WHERE stead_owner=$1", userId)

	if err != nil {
		return nil
	}
	log.Println("Got a list of plants")

	var plants []SerializedPlant

	for rows.Next() {
		var p SerializedPlant
		var xp int
		var kind string
		var next_yield time.Time
		var xp_multiplier float32
		err = rows.Scan(&kind, &xp, &xp_multiplier, &next_yield)
		if err != nil {
			return err
		}

		if xp > 10 {

			xppy := xpPerYield(xp)
			xp_to_go := getXpRemaining(xp)

			// MAGIC: 10 XP/sec
			p.TtYield = (float32(xppy) - float32(xp%xppy)) / 10 * 1000 / xp_multiplier
			p.YieldProg = float32(xp%xppy) / float32(xppy)
			p.TtLevelUp = float32(xp_to_go) / 10 * 1000 / xp_multiplier
			p.LvlupProg = float32(xp_to_go) / float32(getLevelCost(lvlfromxp(xp)+1))

			p.Lvl = lvlfromxp(xp)
		} else {
			p.TtYield = math.MaxInt
			p.YieldProg = 0
			p.TtLevelUp = math.MaxInt
			p.LvlupProg = 0
		}

		plants = append(plants, p)
		fmt.Printf("%+v", p)
	}

	fmt.Printf("%v", plants)

	return c.JSON(fiber.Map{
		"inv":    inventory,
		"plants": plants,
	})

}

func Users(c *fiber.Ctx) error {
	db := database.DB

	rows, err := db.Queryx("SELECT username FROM stead")
	if err != nil {
		return err
	}

	var res []string

	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			return err
		}

		res = append(res, username)
	}
	return c.JSON(res)
}

func Manifest(c *fiber.Ctx) error {
	return c.JSON(state.GlobalState.Manifest)
}

func Activity(c *fiber.Ctx) error {
	return c.JSON(state.GlobalState.Activity)
}

// POST
func UseItem(c *fiber.Ctx) error {
	val := &XValidator{
		validator: validate,
	}

	var item models.UseItem
	if err := c.BodyParser(item); err != nil {
		return err
	}
	if errs := val.Validate(item); len(errs) > 0 && errs[0].Error {
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

	manifest := state.GlobalState.Manifest

	if manifest[item.Item] == nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "no such item found",
		}
	}

	if !manifest[item.Item].(map[string]interface{})["usable"].(bool) {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "that's not an item you can use!",
		}
	}

	err := game.TakeItem(c.Locals("username").(string), map[string]interface{}{
		item.Item: 1,
	})

	if err != nil {
		return err
	}

	log.Printf("%s is using item %s", c.Locals("username"), item.Item)
	state.GlobalState.ActivityPush("useitem", map[string]interface{}{
		"who":  c.Locals("username"),
		"what": item.Item,
	})

	u := c.Locals("username").(string)
	i := item.Item

	// Item recipes!
	if i == "bag_egg_t1" {
		game.GiveItem(u, game.ScaleDrop(game.MegaboxDrop(), 0.2))
	}
	if i == "bag_egg_t2" {
		game.GiveItem(u, game.ScaleDrop(game.MegaboxDrop(), 0.6))
	}
	if i == "bag_egg_t3" {
		game.GiveItem(u, game.MegaboxDrop())
	}

	if i == "bbc_egg" || i == "hvv_egg" || i == "cyl_egg" {
		inverted_types := map[string]interface{}{
			"bbc_egg": []string{"hvv_item", "cyl_item"},
			"hvv_egg": []string{"bbc_item", "cyl_item"},
			"cyl_egg": []string{"bbc_item", "hvv_item"},
		}

		var reward map[string]interface{}
		r := rand.Float64()

		if r < 0.1 {
			reward = map[string]interface{}{"land_deed": 1}
		} else if r < 0.5 {
			reward = map[string]interface{}{
				game.Choose[string](inverted_types[i].([]string)): 1,
			}
		} else if r < 0.55 {
			reward = map[string]interface{}{"powder_t1": 1}
		} else {
			reward = map[string]interface{}{"powder_t2": 1}
		}

		game.GiveItem(u, reward)
	}

	if i == "land_deed" {
		game.NewPlant(u, "dirt")
	}

	return c.JSON(fiber.Map{"ok": true})
}

func Craft(c *fiber.Ctx) error {

	return nil
}
