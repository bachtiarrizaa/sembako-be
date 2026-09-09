package seeder

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedTransactions(db *gorm.DB) error {
	// Find cashier user
	var cashierUser entity.User
	if err := db.Where("email = ?", "cashier@sembako.com").First(&cashierUser).Error; err != nil {
		return fmt.Errorf("cashier user not found (run demo seeder first): %v", err)
	}

	// Create or find a shift
	var shift entity.Shift
	err := db.Where("cashier_id = ? AND status = ?", cashierUser.ID, entity.ShiftStatusOpen).First(&shift).Error
	if err == gorm.ErrRecordNotFound {
		shift = entity.Shift{
			ID:             uuid.NewString(),
			CashierID:      cashierUser.ID,
			OpeningBalance: 500000,
			Status:         entity.ShiftStatusOpen,
		}
		if err := db.Create(&shift).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Fetch some product units
	var productUnits []entity.ProductUnit
	if err := db.Limit(5).Find(&productUnits).Error; err != nil {
		return err
	}
	if len(productUnits) == 0 {
		return fmt.Errorf("no product units found (run demo seeder first)")
	}

	paymentMethods := []string{"cash", "qris", "transfer"}

	for i := 1; i <= 30; i++ {
		// Random date within the last 30 days
		randomDaysAgo := rand.Intn(30)
		trxDate := time.Now().AddDate(0, 0, -randomDaysAgo)
		
		receiptNum := fmt.Sprintf("TRX-%s-%04d", trxDate.Format("20060102"), rand.Intn(9999)+1000)
		paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]
		
		total := 0.0
		var items []entity.TransactionItem
		
		// generate 1-3 items per transaction
		numItems := rand.Intn(3) + 1
		for j := 0; j < numItems; j++ {
			pu := productUnits[rand.Intn(len(productUnits))]
			qty := float64(rand.Intn(5) + 1)
			sub := pu.SellingPrice * qty
			
			items = append(items, entity.TransactionItem{
				ID:              uuid.NewString(),
				ProductUnitID:   pu.ID,
				Qty:             qty,
				UnitPrice:       pu.SellingPrice,
				Subtotal:        sub,
				DiscountApplied: 0,
			})
			total += sub
		}

		trx := entity.Transaction{
			ID:            uuid.NewString(),
			ReceiptNumber: receiptNum,
			CashierID:     cashierUser.ID,
			ShiftID:       shift.ID,
			PaymentMethod: paymentMethod,
			Subtotal:      total,
			TotalDiscount: 0,
			Total:         total,
			Status:        "completed",
			Items:         items,
		}

		if err := db.Create(&trx).Error; err != nil {
			return fmt.Errorf("failed to create transaction %d: %v", i, err)
		}

		// Update created_at and updated_at directly via SQL to bypass GORM's autoCreateTime
		if err := db.Exec("UPDATE transactions SET created_at = ?, updated_at = ? WHERE id = ?", trxDate, trxDate, trx.ID).Error; err != nil {
			return fmt.Errorf("failed to update transaction date %d: %v", i, err)
		}
	}

	return nil
}
