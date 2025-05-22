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

func (s *StoreService) GetStoreByID(id uint64) (*models.Store, error) {
	return s.storeRepo.GetByID(id)
}

func (s *StoreService) GetAllStores(page int, size int) (interface{}, error) {
	count, err := s.storeRepo.Count()
	if err != nil {
		return nil, err
	}

	stores, err := s.storeRepo.GetAll(page, size)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stores": stores,
		"count":  count,
	}, nil
}

func (s *StoreService) UpdateStore(store *models.Store) error {
	return s.storeRepo.Update(store)
}

func (s *StoreService) DeleteStore(store *models.Store) error {
	return s.storeRepo.Delete(store)
}

func (s *StoreService) GetStoreByAdminID(adminID uint64) ([]models.Store, error) {
	return s.storeRepo.GetByAdminID(adminID)
}