package seeder

import "gorm.io/gorm"

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
