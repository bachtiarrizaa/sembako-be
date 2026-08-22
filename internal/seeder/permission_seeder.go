package seeder

import (
	"log"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) error {
	ptrString := func(s string) *string {
		return &s
	}

	permissions := []entity.Permission{
		// Level 1: Root / Group Menus
		{ID: "10000000-0000-0000-0000-000000000000", Name: "dashboard", Description: "Dashboard", Type: "menu", Path: ptrString("/dashboard")},
		{ID: "20000000-0000-0000-0000-000000000000", Name: "pos:create", Description: "POS / Kasir", Type: "menu", Path: ptrString("/pos")},
		{ID: "30000000-0000-0000-0000-000000000000", Name: "transactions:read", Description: "Riwayat Transaksi", Type: "menu", Path: ptrString("/transactions")},
		{ID: "40000000-0000-0000-0000-000000000000", Name: "shifts:read", Description: "Shift Kasir", Type: "menu", Path: ptrString("/shifts")},
		{ID: "50000000-0000-0000-0000-000000000000", Name: "products", Description: "Produk & Satuan", Type: "menu", Path: ptrString("/products")},
		{ID: "60000000-0000-0000-0000-000000000000", Name: "discounts", Description: "Promosi & Diskon", Type: "menu", Path: ptrString("/discounts")},
		{ID: "70000000-0000-0000-0000-000000000000", Name: "inventory", Description: "Stok & Opname", Type: "menu", Path: ptrString("/inventory")},
		{ID: "80000000-0000-0000-0000-000000000000", Name: "suppliers-purchases", Description: "Pengadaan", Type: "menu", Path: ptrString("/suppliers-purchases")},
		{ID: "90000000-0000-0000-0000-000000000000", Name: "customers-loyalty", Description: "Pelanggan & Poin", Type: "menu", Path: ptrString("/customers-loyalty")},
		{ID: "a0000000-0000-0000-0000-000000000000", Name: "reports", Description: "Laporan", Type: "menu", Path: ptrString("/reports")},
		{ID: "b0000000-0000-0000-0000-000000000000", Name: "employees", Description: "Pegawai & Akun", Type: "menu", Path: ptrString("/employees")},
		{ID: "c0000000-0000-0000-0000-000000000000", Name: "settings:read", Description: "Pengaturan", Type: "menu", Path: ptrString("/settings")},

		// Level 2: Sub-menus (Parent ID points to Level 1)
		// Products
		{ID: "51000000-0000-0000-0000-000000000000", Name: "products:read", Description: "Daftar Produk", ParentID: ptrString("50000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/products/list")},
		{ID: "52000000-0000-0000-0000-000000000000", Name: "units:read", Description: "Satuan", ParentID: ptrString("50000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/products/units")},
		{ID: "53000000-0000-0000-0000-000000000000", Name: "categories:read", Description: "Kategori", ParentID: ptrString("50000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/products/categories")},
		// Discounts
		{ID: "61000000-0000-0000-0000-000000000000", Name: "discounts:read", Description: "Daftar Diskon", ParentID: ptrString("60000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/discounts")},
		{ID: "62000000-0000-0000-0000-000000000000", Name: "product-discounts:read", Description: "Diskon per Produk", ParentID: ptrString("60000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/discounts/products")},
		// Inventory
		{ID: "71000000-0000-0000-0000-000000000000", Name: "stocks:read", Description: "Mutasi Stok", ParentID: ptrString("70000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/inventory/stock")},
		{ID: "72000000-0000-0000-0000-000000000000", Name: "opname:approve", Description: "Persetujuan Opname", ParentID: ptrString("70000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/inventory/opname")},
		// Suppliers & Purchases
		{ID: "81000000-0000-0000-0000-000000000000", Name: "suppliers:read", Description: "Supplier", ParentID: ptrString("80000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/suppliers")},
		{ID: "82000000-0000-0000-0000-000000000000", Name: "purchases:read", Description: "Pembelian", ParentID: ptrString("80000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/purchases")},
		// Customers & Loyalty
		{ID: "91000000-0000-0000-0000-000000000000", Name: "customers:read", Description: "Pelanggan", ParentID: ptrString("90000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/customers")},
		{ID: "92000000-0000-0000-0000-000000000000", Name: "loyalty:write", Description: "Pengaturan Poin", ParentID: ptrString("90000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/loyalty-settings")},
		// Reports
		{ID: "a1000000-0000-0000-0000-000000000000", Name: "reports:read", Description: "Semua Laporan", ParentID: ptrString("a0000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/reports")},
		// Employees
		{ID: "b1000000-0000-0000-0000-000000000000", Name: "users:read", Description: "Pegawai", ParentID: ptrString("b0000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/users")},
		{ID: "b2000000-0000-0000-0000-000000000000", Name: "roles:read", Description: "Role & Hak Akses", ParentID: ptrString("b0000000-0000-0000-0000-000000000000"), Type: "menu", Path: ptrString("/roles")},

		// Level 3: Actions (Parent ID points to Level 2 sub-menus or Level 1 menus directly where applicable)
		// Transactions (from Level 1 directly)
		{Name: "transactions:void", Description: "Melakukan void transaksi", ParentID: ptrString("30000000-0000-0000-0000-000000000000"), Type: "action"},
		// Shifts (from Level 1 directly)
		{Name: "shifts:create", Description: "Membuka shift kasir", ParentID: ptrString("40000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "shifts:close", Description: "Menutup shift sendiri", ParentID: ptrString("40000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "shifts:force-close", Description: "Menutup shift kasir lain secara paksa", ParentID: ptrString("40000000-0000-0000-0000-000000000000"), Type: "action"},
		// Settings (from Level 1 directly)
		{Name: "settings:update", Description: "Mengubah pengaturan global", ParentID: ptrString("c0000000-0000-0000-0000-000000000000"), Type: "action"},

		// Under Products sub-menu (Level 2: 51000000-...)
		{Name: "products:create", Description: "Membuat produk baru", ParentID: ptrString("51000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "products:update", Description: "Mengubah data produk", ParentID: ptrString("51000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "products:delete", Description: "Menghapus produk", ParentID: ptrString("51000000-0000-0000-0000-000000000000"), Type: "action"},
		// Under Units sub-menu (Level 2: 52000000-...)
		{Name: "units:create", Description: "Membuat satuan baru", ParentID: ptrString("52000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "units:update", Description: "Mengubah data satuan", ParentID: ptrString("52000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "units:delete", Description: "Menghapus satuan", ParentID: ptrString("52000000-0000-0000-0000-000000000000"), Type: "action"},
		// Under Categories sub-menu (Level 2: 53000000-...)
		{Name: "categories:create", Description: "Membuat kategori baru", ParentID: ptrString("53000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "categories:update", Description: "Mengubah data kategori", ParentID: ptrString("53000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "categories:delete", Description: "Menghapus kategori", ParentID: ptrString("53000000-0000-0000-0000-000000000000"), Type: "action"},

		// Under Discounts sub-menu (Level 2: 61000000-...)
		{Name: "discounts:create", Description: "Membuat template diskon baru", ParentID: ptrString("61000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "discounts:update", Description: "Mengubah data diskon", ParentID: ptrString("61000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "discounts:delete", Description: "Menghapus diskon", ParentID: ptrString("61000000-0000-0000-0000-000000000000"), Type: "action"},
		// Under Product Discounts sub-menu (Level 2: 62000000-...)
		{Name: "product-discounts:create", Description: "Menghubungkan diskon ke produk", ParentID: ptrString("62000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "product-discounts:update", Description: "Mengubah pengaturan diskon produk", ParentID: ptrString("62000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "product-discounts:delete", Description: "Mencabut diskon dari produk", ParentID: ptrString("62000000-0000-0000-0000-000000000000"), Type: "action"},

		// Under Inventory Stok sub-menu (Level 2: 71000000-...)
		{Name: "stocks:create", Description: "Menambah penyesuaian stok manual", ParentID: ptrString("71000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "stocks:update", Description: "Mengubah penyesuaian stok manual", ParentID: ptrString("71000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "stocks:delete", Description: "Menghapus penyesuaian stok manual", ParentID: ptrString("71000000-0000-0000-0000-000000000000"), Type: "action"},
		// Under Approval Opname sub-menu (Level 2: 72000000-...)
		{Name: "opname:create", Description: "Melakukan/mengajukan stok opname fisik", ParentID: ptrString("72000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "opname:read", Description: "Melihat riwayat opname", ParentID: ptrString("72000000-0000-0000-0000-000000000000"), Type: "action"},

		// Under Suppliers sub-menu (Level 2: 81000000-...)
		{Name: "suppliers:create", Description: "Membuat supplier baru", ParentID: ptrString("81000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "suppliers:update", Description: "Mengubah data supplier", ParentID: ptrString("81000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "suppliers:delete", Description: "Menghapus supplier", ParentID: ptrString("81000000-0000-0000-0000-000000000000"), Type: "action"},
		// Under Purchases sub-menu (Level 2: 82000000-...)
		{Name: "purchases:create", Description: "Mencatat transaksi pembelian baru", ParentID: ptrString("82000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "purchases:update", Description: "Mengubah catatan pembelian", ParentID: ptrString("82000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "purchases:delete", Description: "Menghapus catatan pembelian", ParentID: ptrString("82000000-0000-0000-0000-000000000000"), Type: "action"},

		// Under Customers sub-menu (Level 2: 91000000-...)
		{Name: "customers:create", Description: "Membuat customer baru", ParentID: ptrString("91000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "customers:update", Description: "Mengubah data customer", ParentID: ptrString("91000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "customers:delete", Description: "Menghapus customer", ParentID: ptrString("91000000-0000-0000-0000-000000000000"), Type: "action"},

		// Under Employees sub-menu (Level 2: b1000000-...)
		{Name: "users:create", Description: "Mengundang pegawai baru", ParentID: ptrString("b1000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "users:update", Description: "Mengubah data pegawai", ParentID: ptrString("b1000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "users:delete", Description: "Menghapus pegawai", ParentID: ptrString("b1000000-0000-0000-0000-000000000000"), Type: "action"},
		// Under Roles sub-menu (Level 2: b2000000-...)
		{Name: "roles:create", Description: "Membuat role baru", ParentID: ptrString("b2000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "roles:update", Description: "Mengubah permissions role", ParentID: ptrString("b2000000-0000-0000-0000-000000000000"), Type: "action"},
		{Name: "roles:delete", Description: "Menghapus role", ParentID: ptrString("b2000000-0000-0000-0000-000000000000"), Type: "action"},
	}

	for _, perm := range permissions {
		var existing entity.Permission
		var result *gorm.DB
		if perm.ID != "" {
			result = db.Where("id = ? OR name = ?", perm.ID, perm.Name).First(&existing)
		} else {
			result = db.Where("name = ?", perm.Name).First(&existing)
		}

		if result.Error != nil {
			if err := db.Create(&perm).Error; err != nil {
				return err
			}
		} else {
			existing.Description = perm.Description
			existing.ParentID = perm.ParentID
			existing.Type = perm.Type
			existing.Path = perm.Path
			if err := db.Save(&existing).Error; err != nil {
				return err
			}
		}
	}
	log.Println("seeding permissions done!")

	var allPerms []entity.Permission
	if err := db.Find(&allPerms).Error; err != nil {
		return err
	}

	var adminRole entity.Role
	if err := db.First(&adminRole, "name = 'admin'").Error; err != nil {
		return err
	}
	if err := db.Model(&adminRole).Association("Permissions").Replace(allPerms); err != nil {
		return err
	}

	var cashierPerms []entity.Permission
	cashierPermNames := []string{
		"dashboard",
		"pos:create",
		"transactions:read",
		"shifts:read", "shifts:create", "shifts:close",
		"products", "products:read",
		"discounts", "discounts:read",
		"inventory", "stocks:read", "opname:create", "opname:read",
		"customers-loyalty", "customers:read", "customers:create", "customers:update", "customers:delete",
	}
	if err := db.Where("name IN ?", cashierPermNames).Find(&cashierPerms).Error; err != nil {
		return err
	}

	var cashierRole entity.Role
	if err := db.First(&cashierRole, "name = 'cashier'").Error; err != nil {
		return err
	}
	if err := db.Model(&cashierRole).Association("Permissions").Replace(cashierPerms); err != nil {
		return err
	}

	log.Println("assigning role permissions done!")
	return nil
}
