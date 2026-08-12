package entity

import "time"

type Permission struct {
	ID          string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string       `gorm:"column:name;unique;not null" json:"name"`
	Description string       `gorm:"column:description" json:"description"`
	ParentID    *string      `gorm:"column:parent_id;type:uuid" json:"parentId"`
	Type        string       `gorm:"column:type;not null;default:action" json:"type"`
	Path        *string      `gorm:"column:path" json:"path"`
	Parent      *Permission  `gorm:"foreignKey:ParentID" json:"-"`
	Children    []Permission `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Permission) TableName() string {
	return "permissions"
}
