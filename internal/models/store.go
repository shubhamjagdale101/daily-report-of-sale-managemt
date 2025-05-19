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

	AdminID uint `gorm:"not null" json:"admin_id"`
	CreatedBy Admin `gorm:"foreignKey:AdminID;constraint:onUpdate:CASCADE,OnDelete:SET NULL" json:"created_by,omitempty"`
	ManagedBy []Admin `gorm:"many2many:store_admins" json:"managed_by,omitempty"`
}

func (s *Store) HaveAccessToManage(adminId uint) bool {
	for _, admin := range s.ManagedBy {
		if admin.ID == adminId {
			return true
		}
	}
	return false
}