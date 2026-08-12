package seeder

import (
	"log"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roleNames := []string{"admin", "cashier"}

	for _, name := range roleNames {
		role := entity.Role{Name: name}
		result := db.FirstOrCreate(&role, entity.Role{Name: name})
		if result.Error != nil {
			return result.Error
		}
	}
	log.Println("seeding roles done!")
	return nil
}
