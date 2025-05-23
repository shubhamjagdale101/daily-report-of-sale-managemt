package controllers

import (
	"net/http"
	"strconv"
	"time"
	"gold-management-system/internal/services"
	"gold-management-system/internal/utils"
	"gold-management-system/dto"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service *services.TransactionService
}

func NewTransactionController(service *services.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

func (ctrl *TransactionController) CreateTransaction(c *gin.Context) {
	var input struct {
		CustomerID  uint    `json:"customer_id" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=buy sell"`
		GoldWeight  float64 `json:"gold_weight" binding:"required,gt=0"`
		GoldPrice   float64 `json:"gold_price" binding:"required,gt=0"`
		Description string  `json:"description"`
		StoreName  string  `json:"store_name" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"oneof=cash borrowed_gold borrowed_money upi"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	adminIdInterface, ok := c.Get("admin_id")
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Admin ID not found in context")
		return
	}
	adminID, ok := adminIdInterface.(uint)
	if !ok {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid admin ID")
		return
	}

	transaction, err := ctrl.service.CreateTransaction(
		adminID,
		input.CustomerID,
		input.Type,
		input.GoldWeight,
		input.GoldPrice,
		input.Description,
		input.StoreName,
		input.PaymentMethod,
		adminID,
	)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Transaction created successfully", transaction)
}

func (ctrl *TransactionController) GetTransactionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	transaction, err := ctrl.service.GetTransactionByID(uint(id))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

func (ctrl *TransactionController) GetAllTransactions(c *gin.Context) {
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

	transactions, err := ctrl.service.GetAllTransactions(page, size)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}

func (ctrl *TransactionController) GetDailyReport(c *gin.Context) {
	dateStr := c.Query("date")
	var date time.Time
	var err error

	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
	}

	transactions, err := ctrl.service.GetDailyReport(date)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Daily report retrieved successfully", transactions)
}

func (ctrl *TransactionController) GetReport(c *gin.Context) {	
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" || endDateStr == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	startDate, err = time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD")
		return
	}

	endDate, err = time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD")
		return
	}

	if startDate.After(endDate) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Start date cannot be after end date")
		return
	}

	transactions, err := ctrl.service.GetReport(startDate, endDate)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Create CSV file
	buf, err := ctrl.service.CreateCSVFile("transactions_report.csv", transactions)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())		
		return
	}

	utils.SendCSVResponse(c, buf, "transactions_report.csv")
}

func (ctrl *TransactionController) GetMonthlyReport(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	var year int
	var month time.Month
	var err error

	if yearStr == "" || monthStr == "" {
		now := time.Now()
		year = now.Year()
		month = now.Month()
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid year")
			return
		}

		monthInt, err := strconv.Atoi(monthStr)
		if err != nil || monthInt < 1 || monthInt > 12 {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid month")
			return
		}
		month = time.Month(monthInt)
	}

	transactions, err := ctrl.service.GetMonthlyReport(year, month)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Monthly report retrieved successfully", transactions)
}

func (ctrl *TransactionController) GetTransactionsByFilters(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "0")
	sizeStr := c.DefaultQuery("size", "10")
	

	var filters []dto.Filter
	if err := c.ShouldBindJSON(&filters); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

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

	res, err := ctrl.service.GetTransactionsByFilters(page, size, filters)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", res)
}

func (ctrl *TransactionController) GetDashboard(c *gin.Context) {
	dateStr := c.Query("date")
	durationStr := c.Query("duration")

	var date time.Time
	var err error
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
	}

	var duration int
	if durationStr == "" {
		duration = 1
	} else {
		duration, err = strconv.Atoi(durationStr)
		if err != nil || duration < 1 {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid duration value")
			return
		}
	}

	dashboard, err := ctrl.service.GetDashboard(date, duration)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Daily dashboard retrieved successfully", dashboard)
}
