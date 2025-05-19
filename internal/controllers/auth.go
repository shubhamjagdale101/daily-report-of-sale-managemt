package controllers

import (
	"net/http"
	"gold-management-system/internal/services"
	"gold-management-system/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := ctrl.service.Login(input.Email, input.Password, c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Login successful", gin.H{"name" : admin.Name, "email": admin.Email, "createdAt" : admin.CreatedAt})
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := ctrl.service.Register(input.Name, input.Email, input.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Registration failed")
		return
	}

	if err :=ctrl.service.GenerateJWTToken(admin.ID, c); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Don't return password hash
	admin.Password = ""
	utils.SendSuccessResponse(c, http.StatusCreated, "Registration successful", admin)
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	c.SetCookie("bearer-token", "", -1, "/", "", false, true)
	utils.SendSuccessResponse(c, http.StatusOK, "Logout successful", nil)
}