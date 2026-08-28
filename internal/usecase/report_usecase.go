package usecase

import (
	"context"
	"time"

	"github.com/bachtiarrizaa/sembako-be/internal/model"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/errs"
	"github.com/bachtiarrizaa/sembako-be/internal/pkg/utils"
	"github.com/bachtiarrizaa/sembako-be/internal/repository"
)

type ReportUsecase interface {
	GetSalesReport(ctx context.Context, req model.GetReportRequest) (*model.SalesReportResponse, utils.Pagination, error)
	GetProfitMarginReport(ctx context.Context, req model.GetReportRequest) (*model.ProfitMarginReportResponse, utils.Pagination, error)
	GetPaymentMethodReport(ctx context.Context, req model.GetReportRequest) (*model.PaymentMethodReportResponse, error)
	GetCashierReport(ctx context.Context, req model.GetReportRequest) (*model.CashierReportResponse, utils.Pagination, error)
	GetTopSellingReport(ctx context.Context, req model.GetReportRequest) (*model.TopSellingReportResponse, utils.Pagination, error)
	GetInventoryValuationReport(ctx context.Context, req model.GetReportRequest) (*model.InventoryValuationReportResponse, utils.Pagination, error)
}

type reportUsecaseImpl struct {
	reportRepo repository.ReportRepository
}

func NewReportUsecase(reportRepo repository.ReportRepository) ReportUsecase {
	return &reportUsecaseImpl{
		reportRepo: reportRepo,
	}
}

func (u *reportUsecaseImpl) parseDateRange(req model.GetReportRequest) (*time.Time, *time.Time, error) {
	var startDateTime *time.Time
	var endDateTime *time.Time

	const dateFormat = "2006-01-02"

	// 1. Parsing Start Date jika ada
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse(dateFormat, *req.StartDate)
		if err != nil {
			return nil, nil, errs.NewBadRequest("invalid startDate format, expected YYYY-MM-DD")
		}
		startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		startDateTime = &startOfDay
	}

	// 2. Parsing End Date jika ada
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse(dateFormat, *req.EndDate)
		if err != nil {
			return nil, nil, errs.NewBadRequest("invalid endDate format, expected YYYY-MM-DD")
		}
		endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local)
		endDateTime = &endOfDay
	}

	// 3. Validasi rentang tanggal
	if startDateTime != nil && endDateTime != nil && startDateTime.After(*endDateTime) {
		return nil, nil, errs.NewBadRequest("startDate cannot be after endDate")
	}

	return startDateTime, endDateTime, nil
}

func (u *reportUsecaseImpl) GetSalesReport(ctx context.Context, req model.GetReportRequest) (*model.SalesReportResponse, utils.Pagination, error) {
	startDateTime, endDateTime, err := u.parseDateRange(req)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var cashierID *string
	if req.CashierID != nil && *req.CashierID != "" {
		cashierID = req.CashierID
	}

	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	res, total, err := u.reportRepo.GetSalesReport(ctx, startDateTime, endDateTime, cashierID, page, limit)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	pagination := utils.BuildPagination(page, limit, total)
	return res, pagination, nil
}

func (u *reportUsecaseImpl) GetProfitMarginReport(ctx context.Context, req model.GetReportRequest) (*model.ProfitMarginReportResponse, utils.Pagination, error) {
	startDateTime, endDateTime, err := u.parseDateRange(req)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var cashierID *string
	if req.CashierID != nil && *req.CashierID != "" {
		cashierID = req.CashierID
	}

	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	res, total, err := u.reportRepo.GetProfitMarginReport(ctx, startDateTime, endDateTime, cashierID, page, limit)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	pagination := utils.BuildPagination(page, limit, total)
	return res, pagination, nil
}

func (u *reportUsecaseImpl) GetPaymentMethodReport(ctx context.Context, req model.GetReportRequest) (*model.PaymentMethodReportResponse, error) {
	startDateTime, endDateTime, err := u.parseDateRange(req)
	if err != nil {
		return nil, err
	}

	var cashierID *string
	if req.CashierID != nil && *req.CashierID != "" {
		cashierID = req.CashierID
	}

	return u.reportRepo.GetPaymentMethodReport(ctx, startDateTime, endDateTime, cashierID)
}

func (u *reportUsecaseImpl) GetCashierReport(ctx context.Context, req model.GetReportRequest) (*model.CashierReportResponse, utils.Pagination, error) {
	startDateTime, endDateTime, err := u.parseDateRange(req)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var cashierID *string
	if req.CashierID != nil && *req.CashierID != "" {
		cashierID = req.CashierID
	}

	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	res, total, err := u.reportRepo.GetCashierReport(ctx, startDateTime, endDateTime, cashierID, page, limit)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	pagination := utils.BuildPagination(page, limit, total)
	return res, pagination, nil
}

func (u *reportUsecaseImpl) GetTopSellingReport(ctx context.Context, req model.GetReportRequest) (*model.TopSellingReportResponse, utils.Pagination, error) {
	startDateTime, endDateTime, err := u.parseDateRange(req)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	var cashierID *string
	if req.CashierID != nil && *req.CashierID != "" {
		cashierID = req.CashierID
	}

	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	res, total, err := u.reportRepo.GetTopSellingReport(ctx, startDateTime, endDateTime, cashierID, page, limit)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	pagination := utils.BuildPagination(page, limit, total)
	return res, pagination, nil
}

func (u *reportUsecaseImpl) GetInventoryValuationReport(ctx context.Context, req model.GetReportRequest) (*model.InventoryValuationReportResponse, utils.Pagination, error) {
	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	res, total, err := u.reportRepo.GetInventoryValuationReport(ctx, req.Search, page, limit)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	pagination := utils.BuildPagination(page, limit, total)
	return res, pagination, nil
}
