package services

import (
	"gold-management-system/internal/models"
	"gold-management-system/internal/repositories"

	"gorm.io/gorm"
)

type AdminService struct {
	repo *repositories.AdminRepository
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		repo: repositories.NewAdminRepository(db),
	}
}

func (s *AdminService) GetAdminByID(id uint) (*models.Admin, error) {
	return s.repo.GetByID(id)
}

func (s *AdminService) UpdateAdmin(admin *models.Admin) error {
	return s.repo.Update(admin)
}