package transaction

import (
	"bytes"
	"context"
	"fmt"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/xuri/excelize/v2"
)

func (u *transactionUsecaseImpl) ExportTransactionsToExcel(ctx context.Context, req model.ListTransactionsRequest, cashierID string, role string) ([]byte, error) {
	var restrictToCashierID *string
	if role == "cashier" {
		restrictToCashierID = &cashierID
	}

	transactions, err := u.transactionRepo.FindTransactionsForExport(ctx, req, restrictToCashierID)
	if err != nil {
		return nil, errs.NewInternal("failed to fetch transactions for export: " + err.Error())
	}

	f := excelize.NewFile()
	sheetName := "Laporan Transaksi"
	_ = f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"Tanggal Transaksi", "No. Struk", "Kasir", "Pelanggan", "Nama Produk",
		"Qty", "Satuan", "Harga Satuan", "Subtotal Item", "Grand Total", "Metode",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheetName, cell, h)
	}

	// Apply some styling to header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
	})
	_ = f.SetRowStyle(sheetName, 1, 1, headerStyle)

	rowIdx := 2
	for _, trx := range transactions {
		dateStr := trx.CreatedAt.Format("02-Jan-2006 15:04")
		customerName := "- (Umum)"
		if trx.Customer != nil {
			customerName = trx.Customer.Name
		}

		for _, item := range trx.Items {
			productName := ""
			unitName := ""
			if item.ProductUnit.Product.Name != "" {
				productName = item.ProductUnit.Product.Name
			}
			if item.ProductUnit.Unit.Name != "" {
				unitName = item.ProductUnit.Unit.Name
			}

			_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), dateStr)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), trx.ReceiptNumber)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), trx.Cashier.Name)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), customerName)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIdx), productName)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIdx), item.Qty)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIdx), unitName)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowIdx), item.UnitPrice)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowIdx), item.Subtotal)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowIdx), trx.Total)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowIdx), trx.PaymentMethod)

			rowIdx++
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, errs.NewInternal("failed to generate excel file")
	}

	return buf.Bytes(), nil
}
