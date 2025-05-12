package services

import (
	"gold-management-system/internal/models"
	"gold-management-system/internal/repositories"

	"gorm.io/gorm"
)

type StoreService struct {
	storeRepo *repositories.StoreRepository
}

func NewStoreService(db *gorm.DB) *StoreService {
	return &StoreService{
		storeRepo: repositories.NewStoreRepository(db),
	}
}

func (s *StoreService) CreateStore(store *models.Store) error {
	return s.storeRepo.Create(store)
}	

func (s *StoreService) GetStoreByName(name string) (*models.Store, error) {
	return s.storeRepo.GetByName(name)
}

func (s *StoreService) GetAllStores() ([]models.Store, error) {
	return s.storeRepo.GetAll()
}

func (s *StoreService) UpdateStore(store *models.Store) error {
	return s.storeRepo.Update(store)
}

func (s *StoreService) DeleteStore(store *models.Store) error {
	return s.storeRepo.Delete(store)
}