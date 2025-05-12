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

	token, err := ctrl.service.Login(input.Email, input.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	c.SetCookie("bearer-token", token, 3600, "/", "", true, true)
	utils.SendSuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
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

	// Don't return password hash
	admin.Password = ""
	utils.SendSuccessResponse(c, http.StatusCreated, "Registration successful", admin)
}