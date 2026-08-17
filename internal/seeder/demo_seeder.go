package seeder

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/shopspring/decimal"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"gorm.io/gorm"
)

type unitDef struct {
	name string
	conv float64
}

type supplierDef struct {
	name    string
	contact string
	phone   string
	address string
}

type customerDef struct {
	name    string
	phone   string
	address string
	points  int
}

type purchaseDef struct {
	date     string
	unit     string
	qty      float64
	price    float64
	invoice  string
	supplier string
}

type opnameDef struct {
	date     string
	system   float64
	physical float64
	note     string
	status   string
}

type productDef struct {
	name         string
	category     string
	perKg        float64
	perLiter     float64
	perKarung25  float64
	perKarung50  float64
	perKantong5  float64
	minimumStock float64
	marginPct    float64
	purchases    []purchaseDef
	opname       *opnameDef
}

var unitDefs = []unitDef{
	{"Kg", 1},
	{"Liter", 0.8},
	{"Karung 25 Kg", 25},
	{"Karung 50 Kg", 50},
	{"Kantong 5 Kg", 5},
}

var supplierDefs = []supplierDef{
	{"PT Beras Kita Unggul", "Rudi Hartono", "021-5550011", "Jl. Raya Cikampek KM 12, Karawang"},
	{"CV Sawah Makmur", "H. Suryana", "081234567010", "Jl. Raya Indramayu KM 6, Majalengka"},
	{"Bulog Divre Jabar", "Divre Jawa Barat", "022-5201234", "Jl. Soekarno Hatta No. 1, Bandung"},
	{"Penggilingan Padi Subur Jaya", "Pak Joko", "081234567011", "Desa Sukamaju, Karawang"},
	{"CV Beras Sehat Nusantara", "Linda Wijaya", "081234567012", "Jl. Merdeka No. 88, Sukabumi"},
}

var customerDefs = []customerDef{
	{"Warung Bu Siti", "081234567801", "Jl. Melati No. 12", 150},
	{"Ibu Ratna", "081234567802", "Jl. Kenanga No. 3", 80},
	{"Rumah Makan Padang Sederhana", "081234567803", "Jl. Sudirman No. 45", 220},
	{"Toko Sembako Berkah", "081234567804", "Jl. Ahmad Yani No. 7", 95},
	{"H. Ahmad (Grosir)", "081234567805", "Jl. Diponegoro No. 21", 310},
	{"Katering Bu Sri", "081234567806", "Perum Griya Asri Blok C4", 60},
}

var catalog = []productDef{
	{
		name: "Beras Premium Rojolele", category: "Beras Premium",
		perKg: 15000, perLiter: 12000, perKarung25: 360000, perKantong5: 74000,
		minimumStock: 100, marginPct: 12,
		purchases: []purchaseDef{
			{"2026-07-22", "Karung 25 Kg", 20, 345000, "INV/BJS/2026/001", "PT Beras Kita Unggul"},
			{"2026-08-05", "Kantong 5 Kg", 20, 70000, "INV/BJS/2026/008", "PT Beras Kita Unggul"},
		},
	},
	{
		name: "Beras Premium Pandan Wangi", category: "Beras Premium",
		perKg: 16500, perLiter: 13200, perKarung25: 400000, perKantong5: 81000,
		minimumStock: 100, marginPct: 12,
		purchases: []purchaseDef{
			{"2026-07-25", "Karung 25 Kg", 16, 380000, "INV/BJS/2026/003", "PT Beras Kita Unggul"},
			{"2026-08-10", "Kg", 50, 15200, "INV/BJS/2026/011", "PT Beras Kita Unggul"},
		},
		opname: &opnameDef{
			date: "2026-08-06", system: 400, physical: 405,
			note: "Sisa kemasan & kelebihan timbangan", status: "approved",
		},
	},
	{
		name: "Beras Premium Ramos Super", category: "Beras Premium",
		perKg: 15500, perLiter: 12400, perKarung25: 375000, perKantong5: 76000,
		minimumStock: 100, marginPct: 12,
		purchases: []purchaseDef{
			{"2026-07-28", "Karung 25 Kg", 16, 360000, "INV/BJS/2026/005", "PT Beras Kita Unggul"},
			{"2026-08-12", "Kantong 5 Kg", 20, 72000, "INV/BJS/2026/013", "PT Beras Kita Unggul"},
		},
	},
	{
		name: "Beras Medium IR64", category: "Beras Medium",
		perKg: 13500, perLiter: 10800, perKarung25: 325000, perKantong5: 66000,
		minimumStock: 150, marginPct: 10,
		purchases: []purchaseDef{
			{"2026-07-21", "Karung 25 Kg", 20, 312000, "INV/BJS/2026/002", "CV Sawah Makmur"},
			{"2026-08-02", "Karung 25 Kg", 16, 312000, "INV/BJS/2026/007", "CV Sawah Makmur"},
		},
		opname: &opnameDef{
			date: "2026-08-03", system: 900, physical: 888,
			note: "Penyusutan & tumpah saat penyimpanan", status: "approved",
		},
	},
	{
		name: "Beras Medium Setra Ramos", category: "Beras Medium",
		perKg: 14000, perLiter: 11200, perKarung25: 340000, perKantong5: 69000,
		minimumStock: 150, marginPct: 10,
		purchases: []purchaseDef{
			{"2026-07-26", "Karung 25 Kg", 16, 322000, "INV/BJS/2026/004", "CV Sawah Makmur"},
			{"2026-08-07", "Kantong 5 Kg", 20, 65000, "INV/BJS/2026/009", "CV Sawah Makmur"},
		},
	},
	{
		name: "Beras SPHP Bulog", category: "Beras Ekonomis",
		perKg: 12500, perKarung25: 302500, perKarung50: 600000,
		minimumStock: 300, marginPct: 8,
		purchases: []purchaseDef{
			{"2026-07-23", "Karung 50 Kg", 12, 585000, "INV/BJS/2026/006", "Bulog Divre Jabar"},
			{"2026-08-04", "Karung 50 Kg", 10, 585000, "INV/BJS/2026/010", "Bulog Divre Jabar"},
		},
		opname: &opnameDef{
			date: "2026-08-15", system: 1100, physical: 1120,
			note: "Stok fisik lebih banyak dari sistem (sisa karung tidak tercatat)", status: "pending",
		},
	},
	{
		name: "Beras Ekonomis Cap DW", category: "Beras Ekonomis",
		perKg: 12000, perKarung25: 290000, perKarung50: 580000,
		minimumStock: 200, marginPct: 8,
		purchases: []purchaseDef{
			{"2026-07-27", "Karung 50 Kg", 10, 565000, "INV/BJS/2026/015", "CV Sawah Makmur"},
			{"2026-08-11", "Kg", 100, 11400, "INV/BJS/2026/012", "CV Sawah Makmur"},
		},
	},
	{
		name: "Beras Merah Organik", category: "Beras Khusus",
		perKg: 20000, perLiter: 16000, perKarung25: 480000,
		minimumStock: 50, marginPct: 15,
		purchases: []purchaseDef{
			{"2026-07-30", "Kg", 40, 18500, "INV/BJS/2026/016", "CV Beras Sehat Nusantara"},
			{"2026-08-06", "Karung 25 Kg", 4, 445000, "INV/BJS/2026/017", "CV Beras Sehat Nusantara"},
		},
	},
	{
		name: "Beras Hitam Organik", category: "Beras Khusus",
		perKg: 28000, perLiter: 22400, perKarung25: 680000,
		minimumStock: 50, marginPct: 15,
		purchases: []purchaseDef{
			{"2026-08-01", "Kg", 30, 26000, "INV/BJS/2026/018", "CV Beras Sehat Nusantara"},
			{"2026-08-13", "Karung 25 Kg", 2, 630000, "INV/BJS/2026/019", "CV Beras Sehat Nusantara"},
		},
	},
	{
		name: "Beras Ketan Putih", category: "Beras Khusus",
		perKg: 17000, perLiter: 13600, perKarung25: 410000,
		minimumStock: 60, marginPct: 12,
		purchases: []purchaseDef{
			{"2026-07-24", "Karung 25 Kg", 6, 380000, "INV/BJS/2026/020", "Penggilingan Padi Subur Jaya"},
			{"2026-08-08", "Kg", 30, 15200, "INV/BJS/2026/021", "Penggilingan Padi Subur Jaya"},
		},
	},
	{
		name: "Beras Ketan Hitam", category: "Beras Khusus",
		perKg: 22000, perLiter: 17600, perKarung25: 530000,
		minimumStock: 60, marginPct: 12,
		purchases: []purchaseDef{
			{"2026-07-29", "Karung 25 Kg", 4, 490000, "INV/BJS/2026/022", "Penggilingan Padi Subur Jaya"},
			{"2026-08-14", "Kg", 20, 19600, "INV/BJS/2026/023", "Penggilingan Padi Subur Jaya"},
		},
	},
	{
		name: "Beras Pera Cianjur", category: "Beras Medium",
		perKg: 16000, perLiter: 12800, perKarung25: 385000,
		minimumStock: 80, marginPct: 12,
		purchases: []purchaseDef{
			{"2026-07-31", "Karung 25 Kg", 8, 360000, "INV/BJS/2026/024", "Penggilingan Padi Subur Jaya"},
			{"2026-08-09", "Kg", 40, 14400, "INV/BJS/2026/025", "Penggilingan Padi Subur Jaya"},
		},
		opname: &opnameDef{
			date: "2026-08-11", system: 240, physical: 235,
			note: "Angka meragukan, akan dihitung ulang", status: "rejected",
		},
	},
}

func SeedDemoData(db *gorm.DB) error {
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := seedMasterData(tx); err != nil {
			return err
		}
		if err := seedProducts(tx); err != nil {
			return err
		}
		if err := seedInventory(tx); err != nil {
			return err
		}
		if err := seedDiscounts(tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	log.Println("seeding demo data done!")
	return nil
}

func seedMasterData(db *gorm.DB) error {
	for _, u := range unitDefs {
		var ent entity.Unit
		err := db.Where("name = ?", u.name).First(&ent).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&entity.Unit{Name: u.name}).Error; err != nil {
				return err
			}
		}
	}

	categories := []string{"Beras Premium", "Beras Medium", "Beras Ekonomis", "Beras Khusus"}
	for _, c := range categories {
		var ent entity.Category
		err := db.Where("name = ?", c).First(&ent).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&entity.Category{Name: c}).Error; err != nil {
				return err
			}
		}
	}

	for _, s := range supplierDefs {
		var ent entity.Supplier
		err := db.Where("name = ?", s.name).First(&ent).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			supplier := entity.Supplier{
				Name:        s.name,
				ContactName: s.contact,
				Phone:       s.phone,
				Address:     s.address,
				IsActive:    true,
			}
			if err := db.Create(&supplier).Error; err != nil {
				return err
			}
		}
	}

	for _, c := range customerDefs {
		var ent entity.Customer
		err := db.Where("phone_number = ?", c.phone).First(&ent).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			customer := entity.Customer{
				Name:        c.name,
				PhoneNumber: c.phone,
				Address:     c.address,
				TotalPoints: c.points,
				IsActive:    true,
			}
			if err := db.Create(&customer).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func seedProducts(db *gorm.DB) error {
	unitIDs := map[string]string{}
	var unitList []entity.Unit
	if err := db.Find(&unitList).Error; err != nil {
		return err
	}
	for _, u := range unitList {
		unitIDs[u.Name] = u.ID
	}

	categoryIDs := map[string]string{}
	var categoryList []entity.Category
	if err := db.Find(&categoryList).Error; err != nil {
		return err
	}
	for _, c := range categoryList {
		categoryIDs[c.Name] = c.ID
	}

	for _, pd := range catalog {
		var existing entity.Product
		err := db.Where("name = ?", pd.name).First(&existing).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			continue
		}

		minStock := pd.minimumStock
		margin := pd.marginPct
		product := entity.Product{
			CategoryID:             categoryIDs[pd.category],
			Name:                   pd.name,
			BaseUnitID:             unitIDs["Kg"],
			MinimumStock:           &minStock,
			MarginThresholdPercent: &margin,
			IsActive:               true,
		}
		if err := db.Create(&product).Error; err != nil {
			return err
		}

		productUnits := []entity.ProductUnit{
			{
				ProductID:        product.ID,
				UnitID:           unitIDs["Kg"],
				ConversionToBase: 1,
				SellingPrice:     pd.perKg,
				IsBaseUnit:       true,
				IsActive:         true,
			},
		}

		alternates := []struct {
			unit  string
			price float64
		}{
			{"Liter", pd.perLiter},
			{"Karung 25 Kg", pd.perKarung25},
			{"Karung 50 Kg", pd.perKarung50},
			{"Kantong 5 Kg", pd.perKantong5},
		}
		for _, a := range alternates {
			if a.price <= 0 {
				continue
			}
			productUnits = append(productUnits, entity.ProductUnit{
				ProductID:        product.ID,
				UnitID:           unitIDs[a.unit],
				ConversionToBase: conversionByName(a.unit),
				SellingPrice:     a.price,
				IsBaseUnit:       false,
				IsActive:         true,
			})
		}

		if err := db.Create(&productUnits).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedInventory(db *gorm.DB) error {
	adminID, cashierID, err := getActorIDs(db)
	if err != nil {
		return err
	}

	unitIDs := map[string]string{}
	var unitList []entity.Unit
	if err := db.Find(&unitList).Error; err != nil {
		return err
	}
	for _, u := range unitList {
		unitIDs[u.Name] = u.ID
	}

	supplierIDs := map[string]string{}
	var supplierList []entity.Supplier
	if err := db.Find(&supplierList).Error; err != nil {
		return err
	}
	for _, s := range supplierList {
		supplierIDs[s.Name] = s.ID
	}

	productIDs := map[string]string{}
	var productList []entity.Product
	if err := db.Find(&productList).Error; err != nil {
		return err
	}
	for _, p := range productList {
		productIDs[p.Name] = p.ID
	}

	productUnitIDs := map[string]string{}
	var productUnitList []entity.ProductUnit
	if err := db.Find(&productUnitList).Error; err != nil {
		return err
	}
	for _, pu := range productUnitList {
		productUnitIDs[pu.ProductID+"|"+pu.UnitID] = pu.ID
	}

	for _, pd := range catalog {
		productID := productIDs[pd.name]
		running := 0.0
		var batches []*entity.PurchaseBatch

		type invEvent struct {
			date     time.Time
			purchase *purchaseDef
			opname   *opnameDef
		}
		var events []invEvent
		for i := range pd.purchases {
			d, err := time.Parse("2006-01-02", pd.purchases[i].date)
			if err != nil {
				return err
			}
			events = append(events, invEvent{date: d, purchase: &pd.purchases[i]})
		}
		if pd.opname != nil {
			d, err := time.Parse("2006-01-02", pd.opname.date)
			if err != nil {
				return err
			}
			events = append(events, invEvent{date: d, opname: pd.opname})
		}
		sort.Slice(events, func(i, j int) bool { return events[i].date.Before(events[j].date) })

		for _, ev := range events {
			if ev.purchase != nil {
				pur := ev.purchase
				conv := conversionByName(pur.unit)
				qtyBase := pur.qty * conv
				pricePerBase := pur.price / conv
				unitID := unitIDs[pur.unit]
				productUnitID := productUnitIDs[productID+"|"+unitID]
				date := ev.date
				invoice := pur.invoice
				createdAt := date.Add(9 * time.Hour)

				batch := entity.PurchaseBatch{
					ProductID:     productID,
					SupplierID:    supplierIDs[pur.supplier],
					UnitID:        &productUnitID,
					UnitPrice:     &pur.price,
					InitialQty:    qtyBase,
					RemainingQty:  qtyBase,
					PurchasePrice: pricePerBase,
					InvoiceNumber: &invoice,
					PurchaseDate:  date,
					CreatedBy:     adminID,
					CreatedAt:     createdAt,
					UpdatedAt:     createdAt,
				}
				if err := db.Create(&batch).Error; err != nil {
					return err
				}
				batches = append(batches, &batch)

				running += qtyBase
				if err := upsertStock(db, productID, running); err != nil {
					return err
				}

				noteStr := "Invoice: " + pur.invoice
				mutation := entity.StockMutation{
					ProductID:   productID,
					Type:        "in",
					Qty:         qtyBase,
					QtyBefore:   running - qtyBase,
					QtyAfter:    running,
					Source:      "purchase",
					ReferenceID: &batch.ID,
					Note:        &noteStr,
					CreatedBy:   adminID,
					CreatedAt:   createdAt,
					UpdatedAt:   createdAt,
				}
				if err := db.Create(&mutation).Error; err != nil {
					return err
				}
				continue
			}

			if err := seedOpname(db, productID, *ev.opname, adminID, cashierID, batches, &running); err != nil {
				return err
			}
		}
	}

	return nil
}

func seedOpname(
	db *gorm.DB,
	productID string,
	op opnameDef,
	adminID string,
	cashierID string,
	batches []*entity.PurchaseBatch,
	running *float64,
) error {
	date, err := time.Parse("2006-01-02", op.date)
	if err != nil {
		return err
	}

	disc := op.physical - op.system
	note := op.note
	submittedAt := date.Add(9 * time.Hour)

	sc := entity.StockCount{
		ProductID:   productID,
		CountDate:   date,
		SystemQty:   op.system,
		PhysicalQty: op.physical,
		Discrepancy: disc,
		Note:        &note,
		Status:      op.status,
		SubmittedBy: cashierID,
		SubmittedAt: submittedAt,
		CreatedAt:   submittedAt,
		UpdatedAt:   submittedAt,
	}
	if op.status != "pending" {
		approvedAt := date.Add(11 * time.Hour)
		sc.ApprovedBy = &adminID
		sc.ApprovedAt = &approvedAt
	}
	if err := db.Create(&sc).Error; err != nil {
		return err
	}

	if op.status != "approved" {
		return nil
	}

	if disc < 0 {
		remaining := -disc
		for _, b := range batches {
			if remaining <= 0 {
				break
			}
			if b.RemainingQty <= 0 {
				continue
			}
			var deduct float64
			if b.RemainingQty >= remaining {
				deduct = remaining
			} else {
				deduct = b.RemainingQty
			}
			b.RemainingQty -= deduct
			if err := db.Save(b).Error; err != nil {
				return err
			}
			remaining -= deduct
		}

		before := *running
		after := before + disc
		if after < 0 {
			after = 0
		}
		*running = after
		if err := upsertStock(db, productID, after); err != nil {
			return err
		}

		noteMutation := fmt.Sprintf("Penyesuaian opname fisik (selisih kurang: %f)", disc)
		mutation := entity.StockMutation{
			ProductID:   productID,
			Type:        "out",
			Qty:         -disc,
			QtyBefore:   before,
			QtyAfter:    after,
			Source:      "stock_count",
			ReferenceID: &sc.ID,
			Note:        &noteMutation,
			CreatedBy:   adminID,
			CreatedAt:   submittedAt,
			UpdatedAt:   submittedAt,
		}
		return db.Create(&mutation).Error
	}

	if disc > 0 {
		var price float64
		var supplierID string
		if len(batches) > 0 {
			price = batches[len(batches)-1].PurchasePrice
			supplierID = batches[len(batches)-1].SupplierID
		} else {
			var supplier entity.Supplier
			if err := db.Where("is_active = ?", true).First(&supplier).Error; err != nil {
				return err
			}
			supplierID = supplier.ID
		}

		inv := "OPNAME-SURPLUS-" + sc.ID
		batch := entity.PurchaseBatch{
			ProductID:     productID,
			SupplierID:    supplierID,
			InitialQty:    disc,
			RemainingQty:  disc,
			PurchasePrice: price,
			InvoiceNumber: &inv,
			PurchaseDate:  date,
			CreatedBy:     adminID,
			CreatedAt:     submittedAt,
			UpdatedAt:     submittedAt,
		}
		if err := db.Create(&batch).Error; err != nil {
			return err
		}
		batches = append(batches, &batch)

		before := *running
		after := before + disc
		*running = after
		if err := upsertStock(db, productID, after); err != nil {
			return err
		}

		noteMutation := fmt.Sprintf("Penyesuaian opname fisik (selisih lebih: +%f)", disc)
		mutation := entity.StockMutation{
			ProductID:   productID,
			Type:        "in",
			Qty:         disc,
			QtyBefore:   before,
			QtyAfter:    after,
			Source:      "stock_count",
			ReferenceID: &sc.ID,
			Note:        &noteMutation,
			CreatedBy:   adminID,
			CreatedAt:   submittedAt,
			UpdatedAt:   submittedAt,
		}
		return db.Create(&mutation).Error
	}

	return nil
}

func seedDiscounts(db *gorm.DB) error {
	discounts := []entity.Discount{
		{Name: "Diskon Beras Khusus", Type: "percent", Value: decimal.NewFromFloat(5), IsActive: true},
		{Name: "Promo Kantong 5 Kg", Type: "fixed", Value: decimal.NewFromFloat(5000), IsActive: true},
	}
	for _, d := range discounts {
		var existing entity.Discount
		err := db.Where("name = ?", d.Name).First(&existing).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&d).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func upsertStock(db *gorm.DB, productID string, qty float64) error {
	var stock entity.Stock
	err := db.Where("product_id = ?", productID).First(&stock).Error
	if err == gorm.ErrRecordNotFound {
		stock = entity.Stock{ProductID: productID, QtyBaseUnit: qty}
		return db.Create(&stock).Error
	}
	if err != nil {
		return err
	}
	stock.QtyBaseUnit = qty
	return db.Save(&stock).Error
}

func getActorIDs(db *gorm.DB) (adminID, cashierID string, err error) {
	var adminUser entity.User
	if err = db.Where("email = ?", "scoobyd.doo89@gmail.com").First(&adminUser).Error; err != nil {
		return "", "", err
	}
	var cashierUser entity.User
	if err = db.Where("email = ?", "cashier@sembako.com").First(&cashierUser).Error; err != nil {
		return "", "", err
	}
	return adminUser.ID, cashierUser.ID, nil
}

func conversionByName(name string) float64 {
	for _, u := range unitDefs {
		if u.name == name {
			return u.conv
		}
	}
	return 1
}
