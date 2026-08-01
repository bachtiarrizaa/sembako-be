package entity

import "time"

type BlacklistedToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TokenHash string    `gorm:"column:token_hash;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time
}

func (BlacklistedToken) TableName() string {
	return "blacklisted_tokens"
}
