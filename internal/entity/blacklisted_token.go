package entity

import "time"

type BlacklistedToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TokenHash string    `gorm:"column:token_hash;not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expiresAt"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (BlacklistedToken) TableName() string {
	return "blacklisted_tokens"
}
