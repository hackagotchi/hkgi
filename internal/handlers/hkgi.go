package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hackagotchi/hkgi/database"
	"github.com/hackagotchi/hkgi/internal/game"
	"github.com/hackagotchi/hkgi/internal/models"
	"github.com/hackagotchi/hkgi/internal/state"
)

type SerializedPlant struct {
	Kind      string
	Lvl       int
	TtYield   float32
	YieldProg float32
	TtLevelUp float32
	LvlupProg float32
}

func getLevelCost(lvl int) int {
	return int(math.Floor(10 * math.Pow(1.3, float64(lvl+1))))
}

func getXpRemaining(xp int) int {
	lvl := game.LvlFromXp(xp)

	cost := getLevelCost(lvl + 1)

	return int(math.Round(float64(cost) - float64(xp)))
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

			xppy := game.XpPerYield(xp)
			xp_to_go := getXpRemaining(xp)

			// MAGIC: 10 XP/sec
			p.TtYield = (float32(xppy) - float32(xp%xppy)) / 10 * 1000 / xp_multiplier
			p.YieldProg = float32(xp%xppy) / float32(xppy)
			p.TtLevelUp = float32(xp_to_go) / 10 * 1000 / xp_multiplier
			p.LvlupProg = float32(xp_to_go) / float32(getLevelCost(game.LvlFromXp(xp)+1))

			p.Lvl = game.LvlFromXp(xp)
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
	val := &XValidator{
		validator: validate,
	}

	var cReq models.Craft
	if err := c.BodyParser(cReq); err != nil {
		return err
	}

	if errs := val.Validate(cReq); len(errs) > 0 && errs[0].Error {
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
	manifest := state.GlobalState.Manifest
	u := c.Locals("username").(string)

	var p models.Plant
	err := db.Get(&p, "SELECT * FROM plant WHERE stead_owner=(SELECT id FROM stead WHERE username=$1) AND id=$2", u, cReq.PlantId)
	if err != nil {
		return err
	}

	recipe := manifest["plant_recipes"].(map[string]interface{})[p.Kind].([]map[string]interface{})[cReq.RecipeIndex]

	if !(game.LvlFromXp(p.Xp) == recipe["xp"].(int)) {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "come back when you're older, plant!",
		}
	}

	if err = game.TakeItem(u, recipe["needs"].(map[string]interface{})); err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "you can't afford that!",
		}
	}

	state.GlobalState.ActivityPush("craft", map[string]interface{}{
		"who": u,
		"what": map[string]interface{}{
			"recipe_index": cReq.RecipeIndex,
			"plant_id":     cReq.PlantId,
		},
	})

	if recipe["change_plant_to"] != nil {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE plant SET kind=$2, xp=0 WHERE id=$1", cReq.PlantId, recipe["change_plant_to"].(string))
		if err != nil {
			tx.Rollback()
			game.GiveItem(u, recipe["needs"].(map[string]interface{}))
			return err
		}
		tx.Commit()
	}

	if recipe["make_item"] != nil {
		switch recipe["make_item"].(type) {
		case string:
			game.GiveItem(u, map[string]interface{}{
				recipe["make_item"].(string): 1,
			})
		default:
			one_of := recipe["make_item"].(map[string]interface{})["one_of"].([][]any)
			r := rand.Float64()
			i := 0
			var item string
			for next := true; next; next = r > 0 && i < len(one_of) {
				pair := one_of[i]
				chance := pair[0].(float64)
				item = pair[1].(string)
				r -= chance
			}

			game.GiveItem(u, map[string]interface{}{
				item: 1,
			})
		}
	}

	return nil
}
