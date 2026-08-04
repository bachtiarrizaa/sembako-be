package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"column:name;not null" json:"name"`
	Email        string         `gorm:"column:email;unique;not null" json:"email"`
	Username     *string        `gorm:"column:username;unique" json:"username"`
	PasswordHash string         `gorm:"column:password_hash;not null" json:"-"`
	RoleID       string         `gorm:"column:role_id;type:uuid;not null" json:"roleId"`
	Role         Role           `gorm:"foreignKey:RoleID" json:"role"`
	IsActive     bool           `gorm:"column:is_active;not null;default:true" json:"isActive"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
