package model

import "time"

// UserSession of the app in DB
type UserSession struct {
	ID           string    `xorm:"id" json:"id"`
	UserID       string    `xorm:"user_id" json:"user_id"`
	LoginTime    time.Time `xorm:"login_time" json:"login_time,omitempty" r-validate:"password"`
	LastSeenTime time.Time `xorm:"last_seen_time" json:"last_seen_time,omitempty"`
}
