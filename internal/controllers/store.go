package controllers

import (
	"gold-management-system/internal/models"
	"gold-management-system/internal/services"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminIdInterface, ok := c.Get("admin_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminId, ok := adminIdInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	admin, err := sc.adminService.GetAdminByID(adminId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	store.AdminID = adminId
	store.ManagedBy = []models.Admin{*admin}

	if err := sc.storeService.CreateStore(&store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, store)
}

func (sc *StoreController) GetAllStores(c *gin.Context) {
	stores, err := sc.storeService.GetAllStores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stores)
}

func (sc *StoreController) GetStoreByName(c *gin.Context) {
	name := c.Param("name")
	store, err := sc.storeService.GetStoreByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, store)
}

func (sc *StoreController) UpdateStore(c *gin.Context) {
	name := c.Param("name")

	store, err := sc.storeService.GetStoreByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	if err := c.ShouldBindJSON(store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sc.storeService.UpdateStore(store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, store)
}

func (sc *StoreController) DeleteStore(c *gin.Context) {
	name := c.Param("name")

	store, err := sc.storeService.GetStoreByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	if err := sc.storeService.DeleteStore(store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Store deleted"})
}

func (sc *StoreController) UpdateStoreManagedBy(c *gin.Context) {
	var input struct {
		Name      string `json:"name"`
		AdminIds   []uint   `json:"adminIds"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store, err := sc.storeService.GetStoreByName(input.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	adminIdInterface, ok := c.Get("admin_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	adminId, ok := adminIdInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	if store.AdminID != adminId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to manage this store"})
		return
	}

	for _, adminId := range input.AdminIds {
		admin, err := sc.adminService.GetAdminByID(adminId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
			return
		}
		store.ManagedBy = append(store.ManagedBy, *admin)
	}

	if err := sc.storeService.UpdateStore(store); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, store)
}