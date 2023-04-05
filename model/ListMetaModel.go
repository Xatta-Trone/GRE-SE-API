package model

import "time"

type ListMetaModel struct {
	Id         uint64    `json:"id"`
	UserId     uint64    `json:"user_id" db:"user_id"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	Words      string    `json:"words"`
	Visibility int       `json:"visibility"`
	Status     int       `json:"status"`
	CratedAt   time.Time `json:"crated_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
