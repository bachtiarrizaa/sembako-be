package entity

import "time"

type PasswordResetToken struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string     `gorm:"column:user_id;type:uuid;not null"`
	User      User       `gorm:"foreignKey:UserID"`
	TokenHash string     `gorm:"column:token_hash;not null;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
