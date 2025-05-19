package services

import (
	"gold-management-system/internal/config"
	"gold-management-system/internal/models"
	"gold-management-system/internal/repositories"
	"gold-management-system/internal/utils"

	"github.com/gin-gonic/gin"
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

func (s *AuthService) GenerateJWTToken(adminID uint, c *gin.Context)  error {
	token, err := utils.GenerateJWTToken(adminID, s.cfg)
	if err != nil {
		return err
	}

	// Set the token in the cookie
	c.SetCookie("bearer-token", token, 3600, "/", "", false, true)
	return nil
}


func (s *AuthService) Login(email, password string, c *gin.Context) (*models.Admin, error) {
	admin, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := admin.CheckPassword(password); err != nil {
		return nil, err
	}

	if err := s.GenerateJWTToken(admin.ID, c); err != nil {
    	return nil, err
	}

	return admin, nil
}