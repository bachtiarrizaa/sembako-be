package seeder

import (
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

func SeedStoreConfiguration(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.StoreConfiguration{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		address := "Jl. Raya Sembako No. 123, Jakarta"
		phone := "081234567890"
		header := "Selamat Datang di Toko Sembako Jaya!"
		footer := "Terima kasih telah berbelanja! Barang yang sudah dibeli tidak dapat dikembalikan."

		config := entity.StoreConfiguration{
			StoreName:                 "Toko Sembako Jaya",
			StoreAddress:              &address,
			StorePhone:                &phone,
			ReceiptHeaderText:         &header,
			ReceiptFooterText:         &footer,
			ReceiptShowCashierName:    true,
			ReceiptShowCustomerName:   true,
			ShiftDiscrepancyTolerance: 1000,
		}
		return db.Create(&config).Error
	}

	return nil
}
