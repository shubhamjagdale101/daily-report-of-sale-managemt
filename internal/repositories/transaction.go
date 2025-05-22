package repositories

import (
	"errors"
	"gold-management-system/internal/models"
	"log"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) DB() *gorm.DB {
	return r.db
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	return &transaction, err
}

func (r *TransactionRepository) GetAll(page int, size int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Limit(size).Offset(page*size).Order("created_at desc").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetDailyReport(date time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)
	err := r.db.Where("created_at BETWEEN ? AND ?", start, end).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetMonthlyReport(year int, month time.Month) ([]models.Transaction, error) {
	var transactions []models.Transaction
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	err := r.db.Where("created_at BETWEEN ? AND ?", start, end).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetReport(startDate time.Time, endDate time.Time) ([]models.Transaction, error) {	
	var transactions []models.Transaction
	err := r.db.Preload("Customer").Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].PaymentMethod < transactions[j].PaymentMethod
	})
	return transactions, err
}

func (r *TransactionRepository) GetByType(page int, size int, transactionType string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("type = ?", transactionType).Limit(size).Offset((page-1)*size).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetByPaymentMethod(page int, size int, paymentMethod string) ([]models.Transaction, error) {
	log.Print(paymentMethod)
	var transactions []models.Transaction
	err := r.db.Where("payment_method = ?", paymentMethod).Limit(size).Offset((page-1)*size).Find(&transactions).Error
	log.Print(transactions)
	return transactions, err
}

func (r *TransactionRepository) GetByCreatedAt(page int, size int, dates string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	// Split the comma-separated date string
	dateRange := strings.Split(dates, ",")
	if len(dateRange) != 2 {
		return nil, errors.New("invalid date range format; expected 'start,end'")
	}

	startDateStr := strings.TrimSpace(dateRange[0])
	endDateStr := strings.TrimSpace(dateRange[1])

	// Optional: parse and validate the date format (assuming "2006-01-02")
	_, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}
	_, err = time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	// Perform the DB query with pagination
	err = r.db.
		Where("DATE(created_at) BETWEEN ? AND ?", startDateStr, endDateStr).
		Limit(size).
		Offset((page - 1) * size).
		Find(&transactions).Error

	return transactions, err
}