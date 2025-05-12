package repositories

import (
	"gold-management-system/internal/models"
	"gorm.io/gorm"
	"time"
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

func (r *TransactionRepository) GetAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Find(&transactions).Error
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
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error
	return transactions, err
}