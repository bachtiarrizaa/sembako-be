package seeder

import (
	"fmt"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	if err := SeedRoles(db); err != nil {
		return err
	}
	if err := SeedPermissions(db); err != nil {
		return err
	}
	if err := SeedUsers(db); err != nil {
		return err
	}
	if err := SeedStoreConfiguration(db); err != nil {
		return err
	}
	if err := SeedLoyaltySetting(db); err != nil {
		return err
	}
	if err := SeedDemoData(db); err != nil {
		return err
	}
	return nil
}

func SeedByName(db *gorm.DB, name string) error {
	switch name {
	case "roles":
		return SeedRoles(db)
	case "permissions":
		return SeedPermissions(db)
	case "users":
		return SeedUsers(db)
	case "store-config", "store_configuration":
		return SeedStoreConfiguration(db)
	case "loyalty", "loyalty_setting":
		return SeedLoyaltySetting(db)
	case "demo":
		return SeedDemoData(db)
	case "all":
		return SeedAll(db)
	default:
		return fmt.Errorf("unknown seeder target '%s'. Available targets: all, roles, permissions, users, store-config, loyalty, demo", name)
	}
}
