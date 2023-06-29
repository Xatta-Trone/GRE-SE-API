package model

import "time"

type NotificationModel struct {
	Id        string     `json:"id"`
	Content   string     `json:"content"`
	URL       string     `json:"url"`
	UserId    uint64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty" db:"read_at"`
}
