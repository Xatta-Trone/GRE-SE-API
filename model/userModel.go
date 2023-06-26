package model

import "time"

type UserModel struct {
	ID        uint64     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	UserName  string     `json:"username,omitempty"`
	ExpiresOn *time.Time `json:"expires_on" db:"expires_on"`
	CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" db:"updated_at"`
}
