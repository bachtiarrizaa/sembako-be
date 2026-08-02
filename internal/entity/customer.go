package entity

import "time"

type Customer struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(50);unique;not null" json:"name"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(20);unique;not null" json:"phoneNumber"`
	Address     string    `gorm:"column:address;type:text" json:"address"`
	TotalPoints int       `gorm:"column:total_points;not null;default:0" json:"totalPoints"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Customer) TableName() string {
	return "customers"
}
