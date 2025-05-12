package repositories

import (
	"gold-management-system/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) GetByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Preload("Transactions").First(&customer, id).Error
	return &customer, err
}

func (r *CustomerRepository) GetAll() ([]models.Customer, error) {
	var customers []models.Customer
	err := r.db.Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Customer{}, id).Error
}

func (r *CustomerRepository) GetTransactions(customerID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("customer_id = ?", customerID).Find(&transactions).Error
	return transactions, err
}

func (r *CustomerRepository) GetByName(name string) ([]*models.Customer, error) {
	var customers []*models.Customer
	err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&customers).Error
	return customers, err
}