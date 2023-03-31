package model

import "time"

type WordGroupModel struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	NewWords  *string   `json:"new_words" db:"new_words"`
	Words     *string   `json:"words" db:"words"`
	FileName  *string   `json:"file_name" db:"file_name"`
	Status    int       `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WordGroupRelationModel struct {
	WordId      int64     `db:"word_id"`
	WordGroupId int64     `db:"word_group_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
