package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/bachtiarrizaa/sembako-be/internal/config"
	"github.com/bachtiarrizaa/sembako-be/internal/entity"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}

	rolesMap := make(map[string]entity.Role)
	roleNames := []string{"admin", "cashier"}

	for _, name := range roleNames {
		role := entity.Role{Name: name}
		result := db.FirstOrCreate(&role, entity.Role{Name: name})
		if result.Error != nil {
			log.Fatalf("failed to seed role %s: %v", name, result.Error)
		}
		rolesMap[name] = role
	}
	log.Println("seeding roles done!")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash default password: ", err)
	}
	passStr := string(hashedPassword)

	// 3. Seed Users
	ownerUsername := "admin"
	kasirUsername := "kasir"

	users := []entity.User{
		{
			Name:         "Admin Store",
			Email:        "scoobyd.doo89@gmail.com",
			Username:     &ownerUsername,
			PasswordHash: passStr,
			RoleID:       rolesMap["admin"].ID,
			IsActive:     true,
		},
		{
			Name:         "Cashier",
			Email:        "cashier@sembako.com",
			Username:     &kasirUsername,
			PasswordHash: passStr,
			RoleID:       rolesMap["cashier"].ID,
			IsActive:     true,
		},
	}

	for _, user := range users {
		var existing entity.User
		result := db.Where("email = ?", user.Email).First(&existing)
		if result.Error != nil {
			// Not found, create
			if err := db.Create(&user).Error; err != nil {
				log.Fatalf("failed to seed user %s: %v", user.Email, err)
			}
			log.Printf("seeded user: %s (%s)", user.Name, user.Email)
		} else {
			log.Printf("user %s already exists, skipping", user.Email)
		}
	}

	log.Println("seeding users done!")
}
