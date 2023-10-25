package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type User struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UseItem struct {
	Item string `json:"item" validate:"required"`
}

type GibItem struct {
	To     string `json:"username" validate:"required"`
	Item   string `json:"item" validate:"required"`
	Amount int    `json:"amount" validate:"required"`
}

type Craft struct {
	PlantId     int `json:"plant_id" validate:"required"`
	RecipeIndex int `json:"recipe_index" validate:"required"`
}

type Plant struct {
	Id           int     `json:"id"`
	Kind         string  `json:"kind"`
	Xp           int     `json:"xp"`
	XpMultiplier float32 `json:"xp_multiplier"`
	NextYield    float32 `json:"tt_yield"`
}

type Stead struct {
	Id                 int
	Username           string
	Password           string
	Inventory          JsonValue
	Ephemeral_statuses JsonValue
}

type PlantKind string

const (
	DIRT PlantKind = "dirt"
	BBC  PlantKind = "bbc"
	CYL  PlantKind = "cyl"
	HVV  PlantKind = "hvv"
)

type JsonValue map[string]interface{}

func (jval *JsonValue) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &jval)
		return nil
	case string:
		json.Unmarshal([]byte(v), &jval)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
func (jval *JsonValue) Value() (driver.Value, error) {
	return json.Marshal(jval)
}

//func (s *PlantKind) Scan(value interface{}) error {
//	asBytes, ok := value.([]byte)
//	if !ok {
//		return xerrors.New("Scan source is not []byte")
//	}
//	*s = PlantKind(string(asBytes))
//	return nil
//}
//
//func (s PlantKind) Value() (driver.Value, error) {
//	values := map[PlantKind]interface{}{
//		DIRT: nil,
//		BBC:  nil,
//		CYL:  nil,
//		HVV:  nil,
//	}
//
//	if _, ok := values[s]; !ok {
//		return nil, xerrors.New("Wrong value for PlantKind")
//	}
//
//	return string(s), nil
//}
