package model

import (
	"time"

	"github.com/xatta-trone/words-combinator/requests"
)

type WordModel struct {
	Id         int           `json:"id"`
	Word       string        `json:"word"`
	WordData   WordDataModel `json:"word_data" db:"word_data"`
	IsReviewed int           `json:"is_reviewed" db:"is_reviewed"`
	CreatedAt  *time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time    `json:"updated_at" db:"updated_at"`
}

// repository
type WordRepository interface {
	// FindByID(ID int) (*WordModel, error)
	// Save(user *WordModel) error
	FindAll(req requests.WordIndexReqStruct) ([]WordModel, error)
	FindOne(id int) (WordModel, error)
}
