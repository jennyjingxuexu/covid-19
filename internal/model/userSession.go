package model

// UserSession of the app in DB
type UserSession struct {
	ID           string `xorm:"id" json:"id"`
	UserID       string `xorm:"user_id" json:"user_id"`
	LoginTime    string `xorm:"login_time" json:"password,omitempty" r-validate:"password"`
	LastSeenTime string `xorm:"last_seen_time" json:"last_seen_time,omitempty"`
}
