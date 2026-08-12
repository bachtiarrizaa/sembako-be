package seeder

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

func SeedUsers(db *gorm.DB) error {
	var adminRole entity.Role
	if err := db.First(&adminRole, "name = 'admin'").Error; err != nil {
		return err
	}

	var cashierRole entity.Role
	if err := db.First(&cashierRole, "name = 'cashier'").Error; err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passStr := string(hashedPassword)

	ownerUsername := "admin"
	kasirUsername := "kasir"

	users := []entity.User{
		{
			Name:         "Admin Store",
			Email:        "scoobyd.doo89@gmail.com",
			Username:     &ownerUsername,
			PasswordHash: passStr,
			RoleID:       adminRole.ID,
			IsActive:     true,
		},
		{
			Name:         "Cashier",
			Email:        "cashier@sembako.com",
			Username:     &kasirUsername,
			PasswordHash: passStr,
			RoleID:       cashierRole.ID,
			IsActive:     true,
		},
	}

	for _, user := range users {
		var existing entity.User
		result := db.Where("email = ?", user.Email).First(&existing)
		if result.Error != nil {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
			log.Printf("seeded user: %s (%s)", user.Name, user.Email)
		} else {
			log.Printf("user %s already exists, skipping", user.Email)
		}
	}

	log.Println("seeding users done!")
	return nil
}
