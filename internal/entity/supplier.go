package entity

import "time"

type Supplier struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(150);not null" json:"name"`
	ContactName string    `gorm:"column:contact_name;type:varchar(100)" json:"contactName"`
	Phone       string    `gorm:"column:phone;type:varchar(20)" json:"phone"`
	Address     string    `gorm:"column:address;type:text" json:"address"`
	IsActive    bool      `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
