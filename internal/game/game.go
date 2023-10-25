package game

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hackagotchi/hkgi/internal/state"

	"github.com/hackagotchi/hkgi/database"
	"github.com/hackagotchi/hkgi/internal/models"
)

func NewPlant(username string, plant_kind string) error {
	db := database.DB

	var steadId int
	err := db.Get(&steadId, "SELECT id FROM stead WHERE username=$1", username)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO plant (stead_owner, kind, xp, xp_multiplier) VALUES ($1, $2, 0, 0)", steadId, plant_kind)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GiveItem(username string, items map[string]interface{}) error {
	db := database.DB

	var stead map[string]interface{}
	err := db.Get(&stead, "SELECT inventory FROM Stead WHERE username=$1", username)
	if err != nil {
		return err
	}

	for k, v := range items {
		if stead[k] != nil {
			stead[k] = stead[k].(int) + v.(int)
		} else {
			stead[k] = v
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE stead SET inventory=$2 WHERE username=$1", username, stead)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func TakeItem(username string, items map[string]interface{}) error {
	db := database.DB

	var stead map[string]interface{}
	err := db.Get(&stead, "SELECT inventory FROM Stead WHERE username=$1", username)
	if err != nil {
		return err
	}

	for k, v := range items {
		if stead[k] == nil || stead[k].(int) < v.(int) {
			return errors.New("you don't have enough!")
		} else {
			stead[k] = stead[k].(int) - v.(int)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE stead SET inventory=$1 WHERE username=$1", username, stead)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func Choose[T any](arr []T) T {
	rand := rand.Intn(len(arr))
	return arr[rand]
}

func ScaleDrop(drop map[string]interface{}, scale float64) map[string]interface{} {
	for k, v := range drop {
		drop[k] = math.Floor(float64(v.(int)) * scale)
	}
	return drop
}

func MegaboxDrop() map[string]interface{} {
	return map[string]interface{}{
		Choose[string]([]string{"bbc_seed", "cyl_seed", "hvv_seed"}): Choose([]int{
			33, 33, 33, 33, 33,
			38, 38, 38, 38, 38,
			45, 45, 45, 45, 45,
			87}),
		"powder_t2": Choose[int]([]int{3, 3, 3, 3, 3, 3, 3, 5, 5, 5, 5, 8}),
		"powder_t3": Choose[int]([]int{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 3}),
	}
}

func LvlFromXp(xp int) int {
	// get log_1.3...
	log_13 := math.Log(1.3)
	lvl := int(math.Log(float64(xp)/10) / log_13)
	if lvl < 0 {
		lvl = 0
	}
	return lvl
}

func XpPerYield(xp int) int {
	level := LvlFromXp(xp)

	// magic numbers!!

	return int(math.Floor(900 * (1 - float64(level)/27)))
}

const SECONDS_PER_TICK = 0.5

func RunTick() error {
	db := database.DB

	for {
		log.Info("Running game tick...")
		var users []models.Stead
		err := db.Select(&users, "SELECT * FROM stead")
		if err != nil {
			log.Error(err)
			continue
		}

		for _, u := range users {
			for k, v := range u.Ephemeral_statuses {
				status := v.(map[string]interface{})
				tt_expire := status["tt_expire"].(string)
				time_expire, _ := time.Parse(time.UnixDate, tt_expire)
				if time_expire.Before(time.Now()) {
					u.Ephemeral_statuses[k] = nil
				}
			}

			var plants []models.Plant
			err = db.Select(&plants, "SELECT id,kind,xp,xp_multiplier,next_yield FROM plant WHERE stead_owner=$1", u.Id)
			if err != nil {
				log.Error(err)
				continue
			}

			tx, err := db.Begin()
			if err != nil {
				log.Error(err)
				continue
			}
			for _, p := range plants {
				if p.Kind == "dirt" {
					continue
				}
				mult := plant_multiplier(p, u)

				for mult > 0 {
					in_mult := math.Min(mult, 1.0)
					xp_per_tick := SECONDS_PER_TICK * 10 // 10 = XP per tick
					xppy := XpPerYield(p.Xp)

					p.Xp += int(xp_per_tick * in_mult)
					xp_since_yield := p.Xp % xppy
					if float64(xp_since_yield) <= xp_per_tick {
						err := GiveItem(u.Username, map[string]interface{}{fmt.Sprintf("%s_essence", p.Kind): 1})
						if err != nil {
							log.Error(err)
						}

						if rand.Float64() < 0.004 {
							lvl := LvlFromXp(p.Xp)
							if lvl > 5 {
								err := GiveItem(u.Username, map[string]interface{}{fmt.Sprintf("%s_bag_t1", p.Kind): 1})
								if err != nil {
									log.Error(err)
								}
							}
						}
					}

				}
				_, err := tx.Exec("UPDATE plant SET xp=$2 WHERE id=$1", p.Xp, p.Id)
				if err != nil {
					log.Error(err)
					tx.Rollback()
					break
				}
			}
			tx.Commit()

		}
		state.GlobalState.ActivityPrune()
		log.Info("Done")

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func plant_multiplier(p models.Plant, s models.Stead) float64 {
	base_mult := 1.0
	for _, status := range s.Ephemeral_statuses["statuses"].([]map[string]interface{}) {
		if status["kind"].(string) == p.Kind {
			base_mult += status["xp_multiplier"].(float64)
		}
	}
	return base_mult
}
