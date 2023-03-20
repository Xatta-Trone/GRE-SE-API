package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type WordGetStruct struct {
	ID   int64
	Word string
}

type ChatResp struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example"`
	Synonyms   []string `json:"synonyms"`
}

type Result struct {
	ID        uint64
	Word      string
	Google    WordStruct
	Wiki      sql.NullString
	WordsApi  sql.NullString `db:"words_api"` // because sqlx will look for column wordsapi by default
	Thesaurus sql.NullString
	Ninja     sql.NullString
}

type WordStruct struct {
	MainWord        string          `json:"word"`
	Audio           string          `json:"audio"`
	Phonetic        string          `json:"phonetic"`
	PartsOfSpeeches []PartsOfSpeech `json:"parts_of_speeches"`
}

type PartsOfSpeech struct {
	PartsOfSpeech string       `json:"parts_of_speech"`
	Phonetic      string       `json:"phonetic"`
	Audio         string       `json:"audio"`
	Definitions   []Definition `json:"definitions"`
}

type Definition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
}

func (ws *WordStruct) Scan(val interface{}) error {
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
func (ws *WordStruct) Value() (driver.Value, error) {
	return json.Marshal(ws)
}