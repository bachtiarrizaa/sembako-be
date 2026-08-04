package entity

import "time"

type PasswordResetToken struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string     `gorm:"column:user_id;type:uuid;not null" json:"userId"`
	User      User       `gorm:"foreignKey:UserID" json:"user"`
	TokenHash string     `gorm:"column:token_hash;not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null" json:"expiresAt"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"usedAt"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
