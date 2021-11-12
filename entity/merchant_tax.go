package entity

import "time"

// MerchantTax ..
type MerchantTax struct {
	ID         int64      `json:"id" gorm:"primary_key"`
	MerchantID int64      `json:"merchantId" gorm:"not null"`
	Name       string     `json:"name" gorm:"not null"`
	Amount     *float32   `json:"amount,omitempty" gorm:"default:0"`
	IsActive   *bool      `json:"isActive,omitempty" gorm:"default:false"`
	Type       *int       `json:"type,omitempty" gorm:"default:0"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
