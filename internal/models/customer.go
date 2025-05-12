package models

import "time"

type Customer struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Phone        string    `gorm:"unique;not null" json:"phone"`
	Address      string    `json:"address"`
	BorrowedGold  float64   `gorm:"default:0" json:"borrowed_gold"` // in grams
	TotalBought  float64   `gorm:"default:0" json:"total_bought"` // in grams
	TotalSold    float64   `gorm:"default:0" json:"total_sold"`   // in grams
	BorrowedAmount float64   `gorm:"default:0" json:"borrowed_amount"` // in currency
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Transactions []Transaction `gorm:"foreignKey:CustomerID" json:"transactions,omitempty"`
}