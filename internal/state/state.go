package state

import (
	"encoding/json"
	"os"
	"time"
)

type Activity struct {
	Kind string                 `json:"kind"`
	Ts   time.Time              `json:"ts"`
	obj  map[string]interface{} `json:"obj"`
}
type State struct {
	Activity []Activity
	Manifest map[string]interface{}
}

func (s *State) ActivityPush(kind string, obj map[string]interface{}) error {

	s.Activity = append(s.Activity, Activity{
		Kind: kind,
		Ts:   time.Now(),
		obj:  obj,
	})

	return nil
}

func (s *State) ActivityPrune() {
	fiveMin, _ := time.ParseDuration("5m")
	var tmp []Activity

	for _, a := range s.Activity {
		if (time.Now().Sub(a.Ts)) < fiveMin {
			tmp = append(tmp, a)
		}
	}

	s.Activity = tmp
}

var GlobalState State

func InitState() error {
	// Read manifest
	content, err := os.ReadFile("data/manifest.json")
	if err != nil {
		return err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}

	GlobalState = State{
		Manifest: payload,
	}

	return nil
}
