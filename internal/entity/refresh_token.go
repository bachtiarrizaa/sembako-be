package entity

import "time"

type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null" json:"userId"`
	TokenHash string    `gorm:"column:token_hash;not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expiresAt"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
