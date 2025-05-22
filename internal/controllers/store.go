package controllers

import (
	"gold-management-system/internal/models"
	"gold-management-system/internal/services"
	"gold-management-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	storeService *services.StoreService
	adminService *services.AdminService
}

func NewStoreController(storeService *services.StoreService, adminService *services.AdminService) *StoreController {
	return &StoreController{
		storeService: storeService,
		adminService: adminService,
	}
}

func (sc *StoreController) CreateStore(c *gin.Context) {
	var store models.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	adminIdInterface, ok := c.Get("admin_id")
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	adminId, ok := adminIdInterface.(uint)
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	admin, err := sc.adminService.GetAdminByID(adminId)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Admin not found")
		return
	}

	store.AdminID = adminId
	store.ManagedBy = []models.Admin{*admin}

	if err := sc.storeService.CreateStore(&store); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "stores created successfully", store);
}

func (sc *StoreController) GetAllStores(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	// Convert to int
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid page value")
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 10 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid size value")
		return
	}

	res, err := sc.storeService.GetAllStores(page, size)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "stores retrieved successfully", res)
}

func (sc *StoreController) GetStoreByName(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	store, err := sc.storeService.GetStoreByID(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Store not found")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Store retrieved successfully", store)
}

func (sc *StoreController) UpdateStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	store, err := sc.storeService.GetStoreByID(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Store not found")
		return
	}

	if err := c.ShouldBindJSON(store); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := sc.storeService.UpdateStore(store); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SendSuccessResponse(c, http.StatusOK, "Store updated successfully", store)
}

func (sc *StoreController) DeleteStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	store, err := sc.storeService.GetStoreByID(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Store not found")
		return
	}

	if err := sc.storeService.DeleteStore(store); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SendSuccessResponse(c, http.StatusOK, "Store deleted successfully", gin.H{"message": "Store deleted"})
}

func (sc *StoreController) UpdateStoreManagedBy(c *gin.Context) {
	var input struct {
		Name      string `json:"name"`
		AdminIds   []uint   `json:"adminIds"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	store, err := sc.storeService.GetStoreByName(input.Name)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Store not found")
		return
	}

	adminIdInterface, ok := c.Get("admin_id")
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	adminId, ok := adminIdInterface.(uint)
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	if store.AdminID != adminId {
		utils.SendErrorResponse(c, http.StatusForbidden, "You are not authorized to manage this store")
		return
	}

	for _, adminId := range input.AdminIds {
		admin, err := sc.adminService.GetAdminByID(adminId)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusNotFound, "Admin not found")
			return
		}
		store.ManagedBy = append(store.ManagedBy, *admin)
	}

	if err := sc.storeService.UpdateStore(store); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Store updated successfully", store)
}

func (sc *StoreController) GetAllStoresByAdmin(c *gin.Context) {
	adminIdInterface, ok := c.Get("admin_id")
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	adminId, ok := adminIdInterface.(uint64)
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	stores, err := sc.storeService.GetStoreByAdminID(adminId)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SendSuccessResponse(c, http.StatusOK, "Stores retrieved successfully", stores)
}
