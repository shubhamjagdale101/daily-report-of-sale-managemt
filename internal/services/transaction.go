package services

import (
	"bytes"
	"gold-management-system/internal/models"
	"gold-management-system/internal/repositories"
	"time"

	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
)

type TransactionService struct {
	customerRepo    *repositories.CustomerRepository
	transactionRepo *repositories.TransactionRepository
	StoreRepo *repositories.StoreRepository
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		customerRepo:    repositories.NewCustomerRepository(db),
		transactionRepo: repositories.NewTransactionRepository(db),
		StoreRepo: repositories.NewStoreRepository(db),
	}
}

func (s *TransactionService) CreateTransaction(adminID, customerID uint, transactionType string, goldWeight, goldPrice float64, description string, storeName string, paymentMethod string) (*models.Transaction, error) {
	// Get customer
	customer, err := s.customerRepo.GetByID(customerID)
	if err != nil {
		return nil, err
	}

	// Get store
	store, err := s.StoreRepo.GetByName(storeName)
	if err != nil {
		return nil, err
	}

	// Create transaction
	transaction := &models.Transaction{
		CustomerID:  customerID,
		Type:        models.TransactionType(transactionType),
		GoldWeight:  goldWeight,
		GoldPrice:   goldPrice,
		Amount:      goldWeight * goldPrice,
		Description: description,
		PaymentMethod: models.PaymentMethod(paymentMethod),
	}

	// Update customer gold balance
	if transaction.Type == models.Buy {
		customer.TotalBought += goldWeight
		store.TotalGold -= goldWeight

		if transaction.PaymentMethod == models.BorrowedGold {
			customer.BorrowedGold += goldWeight
			store.GoldGiven += goldWeight
		} else if transaction.PaymentMethod == models.BorrowedMoney {
			customer.BorrowedAmount += transaction.Amount
			store.AmountGiven += transaction.Amount
		}
	} else {
		customer.TotalSold += goldWeight
		store.TotalGold += goldWeight

		if transaction.PaymentMethod == models.BorrowedGold {
			customer.BorrowedGold -= goldWeight
			store.GoldTaken += goldWeight
		} else if transaction.PaymentMethod == models.BorrowedMoney {
			customer.BorrowedAmount -= transaction.Amount
			store.AmountTaken += transaction.Amount
		}
	}

	// Start transaction
	tx := s.transactionRepo.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Save(customer).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Save(store).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

func (s *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	return s.transactionRepo.GetAll()
}

func (s *TransactionService) GetDailyReport(date time.Time) ([]models.Transaction, error) {
	return s.transactionRepo.GetDailyReport(date)
}

func (s *TransactionService) GetMonthlyReport(year int, month time.Month) ([]models.Transaction, error) {
	return s.transactionRepo.GetMonthlyReport(year, month)
}

func (s *TransactionService) GetReport(startDate, endDate time.Time) ([]models.Transaction, error) {
	return s.transactionRepo.GetReport(startDate, endDate)
}

func (s *TransactionService) CreateCSVFile(fileName string, transactions []models.Transaction) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	type TransactionCSV struct {
		CustomerName string  `csv:"customer_name"`
		Amount 	 float64 `csv:"amount"`
		TransactionType string  `csv:"transaction_type"`
		GoldWeight float64 `csv:"gold_weight"`
		GoldPrice float64 `csv:"gold_price"`
		PaymentMethod string  `csv:"payment_method"`
		Time string `csv:"time"`
		Description string  `csv:"description"`
	}

	// Define the desired time format
	const timeFormat = "2006-01-02 15:04:05"

	// Convert []models.Transaction to []TransactionCSV
	var csvData []TransactionCSV
	for _, t := range transactions {
		csvData = append(csvData, TransactionCSV{
			CustomerName:    t.Customer.Name,
			Amount:          t.Amount,
			TransactionType: string(t.Type),
			GoldWeight:      t.GoldWeight,
			GoldPrice:       t.GoldPrice,
			PaymentMethod:   string(t.PaymentMethod),
			Time:            t.CreatedAt.Format(timeFormat),
			Description:     t.Description,
		})
	}

	// Marshal the people array to the CSV file
	if err := gocsv.Marshal(&csvData, buf); err != nil {
		return nil, err
	}

	// Return the file for further use (or close later)
	return buf, nil
}