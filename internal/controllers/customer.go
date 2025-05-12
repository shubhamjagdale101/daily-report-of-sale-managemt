package controllers

import (
	"net/http"
	"strconv"
	"gold-management-system/internal/models"
	"gold-management-system/internal/services"
	"gold-management-system/internal/utils"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	service *services.CustomerService
}

func NewCustomerController(service *services.CustomerService) *CustomerController {
	return &CustomerController{service: service}
}

func (ctrl *CustomerController) CreateCustomer(c *gin.Context) {
	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctrl.service.CreateCustomer(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Customer created successfully", input)
}

func (ctrl *CustomerController) GetAllCustomers(c *gin.Context) {
	name := c.Query("name")
	if name != "" {
		customer, err := ctrl.service.GetCustomerByName(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusOK, customer)
		return
	}

	customers, err := ctrl.service.GetAllCustomers()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Customers retrieved successfully", customers)
}

func (ctrl *CustomerController) GetCustomerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := ctrl.service.GetCustomerByID(uint(id))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Customer not found")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Customer retrieved successfully", customer)
}

func (ctrl *CustomerController) UpdateCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Ensure ID in URL matches ID in body
	input.ID = uint(id)

	if err := ctrl.service.UpdateCustomer(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Customer updated successfully", input)
}

func (ctrl *CustomerController) DeleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := ctrl.service.DeleteCustomer(uint(id)); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Customer deleted successfully", nil)
}

func (ctrl *CustomerController) GetCustomerTransactions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	transactions, err := ctrl.service.GetCustomerTransactions(uint(id))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Customer transactions retrieved successfully", transactions)
}

func (ctrl *CustomerController) GetCustomerGoldBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	balance, err := ctrl.service.GetCustomerGoldBalance(uint(id))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Customer gold balance retrieved successfully", gin.H{"balance": balance})
}