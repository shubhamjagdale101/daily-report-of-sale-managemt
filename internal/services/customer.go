package services

import (
	"gold-management-system/internal/models"
	"gold-management-system/internal/repositories"

	"gorm.io/gorm"
)

type CustomerService struct {
	repo *repositories.CustomerRepository
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{
		repo: repositories.NewCustomerRepository(db),
	}
}

func (s *CustomerService) CreateCustomer(customer *models.Customer) error {
	return s.repo.Create(customer)
}

func (s *CustomerService) GetCustomerByID(id uint) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *CustomerService) GetAllCustomers() ([]models.Customer, error) {
	return s.repo.GetAll()
}

func (s *CustomerService) UpdateCustomer(customer *models.Customer) error {
	return s.repo.Update(customer)
}

func (s *CustomerService) DeleteCustomer(id uint) error {
	return s.repo.Delete(id)
}

func (s *CustomerService) GetCustomerByName(name string) ([]*models.Customer, error) {	
	return s.repo.GetByName(name)
}

func (s *CustomerService) GetCustomerTransactions(customerID uint) ([]models.Transaction, error) {
	return s.repo.GetTransactions(customerID)
}

func (s *CustomerService) GetCustomerGoldBalance(customerID uint) (float64, error) {
	customer, err := s.repo.GetByID(customerID)
	if err != nil {
		return 0, err
	}
	return customer.BorrowedGold, nil
}