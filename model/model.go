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
	Google    Google
	Wiki      Wiki     `db:"wiki"`
	WordsApi  WordsApi `db:"words_api"` // because sqlx will look for column words-api by default
	Thesaurus Thesaurus
	Ninja     sql.NullString
}

type Google struct {
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

func (ws *Google) Scan(val interface{}) error {
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
func (ws *Google) Value() (driver.Value, error) {
	return json.Marshal(ws)
}

type Wiki struct {
	MainWord        string              `json:"word"`
	Phonetic        string              `json:"phonetic"`
	PartsOfSpeeches []WikiPartsOfSpeech `json:"meanings"`
}

type WikiPartsOfSpeech struct {
	PartsOfSpeech string                        `json:"partOfSpeech"`
	Synonyms      []string                      `json:"synonyms"`
	Antonyms      []string                      `json:"antonyms"`
	Definition    []WikiPartsOfSpeechDefinition `json:"definitions"`
}

type WikiPartsOfSpeechDefinition struct {
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
	Definition string   `json:"definition"`
}

func (ws *Wiki) Scan(val interface{}) error {
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
func (ws *Wiki) Value() (driver.Value, error) {
	return json.Marshal(ws)
}

type WordsApi struct {
	Word          string        `json:"word"`
	Results       []Results     `json:"results"`
	Frequency     float64       `json:"frequency"`
	Syllables     Syllables     `json:"syllables"`
	Pronunciation Pronunciation `json:"pronunciation"`
}
type Results struct {
	TypeOf       []string `json:"typeOf,omitempty"`
	HasTypes     []string `json:"hasTypes,omitempty"`
	Synonyms     []string `json:"synonyms,omitempty"`
	Definition   string   `json:"definition"`
	Derivation   []string `json:"derivation,omitempty"`
	PartOfSpeech string   `json:"partOfSpeech"`
	HasMembers   []string `json:"hasMembers,omitempty"`
	Examples     []string `json:"examples,omitempty"`
	SimilarTo    []string `json:"similarTo,omitempty"`
	InCategory   []string `json:"inCategory,omitempty"`
}
type Syllables struct {
	List  []string `json:"list"`
	Count int      `json:"count"`
}
type Pronunciation struct {
	All string `json:"all"`
}

func (ws *WordsApi) Scan(val interface{}) error {
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
func (ws *WordsApi) Value() (driver.Value, error) {
	return json.Marshal(ws)
}

type Thesaurus struct {
	Data Data `json:"data"`
}
type Synonyms struct {
	Synonym       []string `json:"synonym"`
	Definition    string   `json:"definition"`
	PartsOfSpeech string   `json:"parts_of_speech"`
}
type Data struct {
	Antonyms []string   `json:"antonyms"`
	Synonyms []Synonyms `json:"synonyms"`
}

func (ws *Thesaurus) Scan(val interface{}) error {
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
func (ws *Thesaurus) Value() (driver.Value, error) {
	return json.Marshal(ws)
}
