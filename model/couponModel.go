package model

import "time"

type CouponModel struct {
	ID      uint64     `json:"id"`
	Coupon  string     `json:"coupon"`
	UserId  *uint64    `json:"user_id" db:"user_id"`
	Used    int        `json:"used"`
	MaxUse  int       `json:"max_use" db:"max_use"`
	Expires *time.Time `json:"expires" db:"expires"`
	Months  int        `json:"months"`
	Type    string     `json:"type"`
}
