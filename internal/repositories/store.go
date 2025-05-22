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

func (r *StoreRepository) GetByID(id uint64) (*models.Store, error) {
	var store models.Store
	err := r.db.Preload("CreatedBy").Preload("ManagedBy").Where("id = ?", id).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) GetByName(storeName string) (*models.Store, error) {
	var store models.Store
	err := r.db.Preload("CreatedBy").Preload("ManagedBy").Where("name = ?", storeName).First(&store).Error
	return &store, err
}

func (r *StoreRepository) GetAll(page int, size int) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Limit(size).Offset(page*size).Order("created_at desc").Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *StoreRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Store{}).Count(&count).Error
	return count, err
}

func (r *StoreRepository) Delete(store *models.Store) error {
	return r.db.Delete(store).Error
}

func (r *StoreRepository) GetByAdminID(adminID uint64) ([]models.Store, error) {
	var stores []models.Store
	err := r.db.Where("admin_id = ?", adminID).Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}
