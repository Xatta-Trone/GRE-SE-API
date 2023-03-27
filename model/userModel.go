package model

import "time"

type UserModel struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	UserName  string     `json:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}