package services

import (
	"gold-management-system/internal/config"
	"gold-management-system/internal/models"
	"gold-management-system/internal/repositories"
	"gold-management-system/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	repo *repositories.AdminRepository
	cfg  *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: repositories.NewAdminRepository(db),
		cfg:  cfg,
	}
}

func (s *AuthService) Register(name, email, password string) (*models.Admin, error) {
	admin := &models.Admin{
		Name:     name,
		Email:    email,
	}

	if err := admin.HashPassword(password); err != nil {
		return nil, err
	}

	if err := s.repo.Create(admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	admin, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if err := admin.CheckPassword(password); err != nil {
		return "", err
	}

	token, err := utils.GenerateJWTToken(admin.ID, s.cfg)
	if err != nil {
		return "", err
	}

	return token, nil
}