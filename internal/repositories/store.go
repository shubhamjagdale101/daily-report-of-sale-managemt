package repositories

import (
	"gold-management-system/internal/models"

	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) Create(store *models.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepository) Update(store *models.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepository) GetByName(storeName string) (*models.Store, error) {
	var store models.Store
	err := r.db.Where("name = ?", storeName).First(&store).Error
	return &store, err
}

func (r *StoreRepository) GetAll() ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *StoreRepository) Delete(store *models.Store) error {
	return r.db.Delete(store).Error
}
