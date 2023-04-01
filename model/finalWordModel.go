package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Combined struct {
	PartsOfSpeech string   `json:"partsOfSpeech"`
	Definitions   []string `json:"definitions"`
	Examples      []string `json:"examples"`
	SynonymsG     []string `json:"synonyms_gre"`
	SynonymsN     []string `json:"synonyms_normal"`
}

func (ws *Combined) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &ws)
		return nil
	case string:
		json.Unmarshal([]byte(v), &ws)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}
func (ws *Combined) Value() (driver.Value, error) {
	return json.Marshal(ws)
}

type WordDataModel struct {
	Word            string     `json:"word"`
	PartsOfSpeeches []Combined `json:"partsOfSpeeches"`
}

func (ws *WordDataModel) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &ws)
		return nil
	case string:
		json.Unmarshal([]byte(v), &ws)
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}
func (ws *WordDataModel) Value() (driver.Value, error) {
	return json.Marshal(ws)
}
