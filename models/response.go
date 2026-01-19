package models

import "encoding/json"

type Response struct {
	State  int             `json:"state"`
	Result json.RawMessage `json:"result"` // Delayed parsing
}
