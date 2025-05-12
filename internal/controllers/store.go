package controllers

import (
	"gold-management-system/internal/models"
	"gold-management-system/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	storeService *services.StoreService
}

func NewStoreController(service *services.StoreService) *StoreController {
	return &StoreController{
		storeService: service,
	}
}

func (sc *StoreController) CreateStore(c *gin.Context) {
	var store models.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
