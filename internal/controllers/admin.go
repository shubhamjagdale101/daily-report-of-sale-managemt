package controllers

import (
	"net/http"
	"gold-management-system/internal/services"
	"gold-management-system/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	service *services.AdminService
}

func NewAdminController(service *services.AdminService) *AdminController {
	return &AdminController{service: service}
}

func (ctrl *AdminController) GetProfile(c *gin.Context) {
	adminID := c.GetUint("admin_id")
	admin, err := ctrl.service.GetAdminByID(adminID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Admin not found")
		return
	}

	// Don't return password hash
	admin.Password = ""
	utils.SendSuccessResponse(c, http.StatusOK, "Admin profile retrieved", admin)
}

func (ctrl *AdminController) UpdateProfile(c *gin.Context) {
	adminID := c.GetUint("admin_id")
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := ctrl.service.GetAdminByID(adminID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Admin not found")
		return
	}

	admin.Name = input.Name
	if err := ctrl.service.UpdateAdmin(admin); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	// Don't return password hash
	admin.Password = ""
	utils.SendSuccessResponse(c, http.StatusOK, "Profile updated successfully", admin)
}