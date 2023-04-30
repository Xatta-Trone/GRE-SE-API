package model

import (
	"time"
)

type WordModel struct {
	Id         int64      `json:"id"`
	Word       string     `json:"word"`
	WordData   WordDataModel `json:"word_data,omitempty" db:"word_data"`
	IsReviewed int        `json:"is_reviewed" db:"is_reviewed"`
	CreatedAt  *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at" db:"updated_at"`
}

// func (ws *WordModel.WordData) Scan(val interface{}) error {
// 	switch v := val.(type) {
// 	case []byte:
// 		json.Unmarshal(v, &ws)
// 		return nil
// 	case string:
// 		json.Unmarshal([]byte(v), &ws)
// 		return nil
// 	case nil:
// 		return nil
// 	default:
// 		return fmt.Errorf("unsupported type: %T", v)
// 	}
// }
// func (ws *WordModel) Value() (driver.Value, error) {
// 	return json.Marshal(ws)
// }
