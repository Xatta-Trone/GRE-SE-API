package model

import (
	"time"
)

type WordModel struct {
	Id         int64           `json:"id"`
	Word       string        `json:"word"`
	WordData   WordDataModel `json:"word_data" db:"word_data"`
	IsReviewed int           `json:"is_reviewed" db:"is_reviewed"`
	CreatedAt  *time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time    `json:"updated_at" db:"updated_at"`
}
