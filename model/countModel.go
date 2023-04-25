package model

type CountModel struct {
	Count int64 `json:"count" db:"count"`
}