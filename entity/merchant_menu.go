package entity

import "time"

// MerchantMenu ...
type MerchantMenu struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	MerchantID  int64      `json:"merchantId" gorm:"not null"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Price       int64      `json:"price" gorm:"not null"`
	Discount    float32    `json:"discount"`
	Qty         int        `json:"qty" gorm:"not null"`
	IsActive    *bool      `json:"isActive,omitempty" gorm:"default:false"`
	IsNearby    bool       `json:"isNearby"`
	IsVisit     bool       `json:"isVisit"`
	Photo       string     `json:"photo"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"-"`
}
