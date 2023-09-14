package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Plant struct {
	Kind         string  `json:"kind"`
	Xp           int     `json:"xp"`
	XpMultiplier float32 `json:"xp_multiplier"`
	NextYield    float32 `json:"tt_yield"`
}

type PlantKind string

const (
	DIRT PlantKind = "dirt"
	BBC  PlantKind = "bbc"
	CYL  PlantKind = "cyl"
	HVV  PlantKind = "hvv"
)

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
