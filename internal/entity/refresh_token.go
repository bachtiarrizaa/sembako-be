package entity

import "time"

type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null"`
	TokenHash string    `gorm:"column:token_hash;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
