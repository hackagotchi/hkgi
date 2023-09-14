package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"git.sr.ht/~muirrum/hkgi/database"
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
	var ephemeral_statuses []string

	err := db.QueryRowx("SELECT inventory, ephemeral_statuses FROM stead WHERE username=$1", user).Scan(&inventory, &ephemeral_statuses)

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
	rows, err := db.Query("SELECT xp, xp_multiplier, next_yield FROM plant WHERE stead_owner=$1", userId)

	if err != nil {
		return nil
	}
	log.Println("Got a list of plants")

	var plants []SerializedPlant

	for rows.Next() {
		var p SerializedPlant
		var xp int
		//		var kind pgtype.Text
		var next_yield time.Time
		var xp_multiplier float32
		err = rows.Scan(xp, xp_multiplier, next_yield)
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
