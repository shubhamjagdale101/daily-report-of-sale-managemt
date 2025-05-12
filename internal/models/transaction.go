package models

import (
	"time"
)

type TransactionType string
type PaymentMethod string

const (
	Buy  TransactionType = "buy"
	Sell TransactionType = "sell"
)

const (
	Cash          PaymentMethod = "cash"
	BorrowedGold  PaymentMethod = "borrowed_gold"
	BorrowedMoney PaymentMethod = "borrowed_money"
	UPI           PaymentMethod = "upi"
)

type Transaction struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	CustomerID  uint            `gorm:"not null" json:"customer_id"`
	Type        TransactionType `gorm:"not null" json:"type"`
	GoldWeight  float64         `gorm:"not null" json:"gold_weight"` // in grams
	GoldPrice   float64         `gorm:"not null" json:"gold_price"`  // per gram
	Amount      float64         `gorm:"not null" json:"amount"`      // gold_weight * gold_price
	PaymentMethod PaymentMethod   `gorm:"not null" json:"payment_method"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}