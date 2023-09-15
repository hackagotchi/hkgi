package game

import (
	"errors"
	"math"
	"math/rand"

	"git.sr.ht/~muirrum/hkgi/database"
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
