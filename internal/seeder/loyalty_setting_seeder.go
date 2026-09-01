package seeder

import (
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

func SeedLoyaltySetting(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.LoyaltySetting{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		setting := entity.LoyaltySetting{
			EarningRate:    10000,
			RedemptionRate: 100,
			MinimumRedeem:  50,
			IsExpiryActive: false,
			ExpiryMonths:   12,
		}
		return db.Create(&setting).Error
	}

	return nil
}
