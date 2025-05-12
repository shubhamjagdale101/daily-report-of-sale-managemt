package models

import (
	"time"
)

type Store struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	TotalGold float64   `gorm:"default:0" json:"total_gold"` 
	TotalAmount float64   `gorm:"default:0" json:"total_amount"` 
	GoldTaken float64 `gorm:"default:0" json:"gold_taken"`
	AmountTaken float64 `gorm:"default:0" json:"amount_taken"`
	GoldGiven float64 `gorm:"default:0" json:"gold_given"`
	AmountGiven float64 `gorm:"default:0" json:"amount_given"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}