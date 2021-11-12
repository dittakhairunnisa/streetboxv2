package entity

import "time"

// TrxOrderProductSales ..
type TrxOrderProductSales struct {
	ID                int64      `json:"id" gorm:"primary_key"`
	UniqueID          string     `json:"uniqueId" gorm:"not null"`
	TrxOrderBillID    int64      `json:"trxOrderBillId" gorm:"not null"`
	MerchantMenuID    int64      `json:"merchantMenuId" gorm:"not null"`
	Name              string     `json:"name" gorm:"not null"`
	Price             float64    `json:"price" gorm:"not null"`
	Qty               int        `json:"qty" gorm:"not null"`
	BusinessDate      time.Time  `json:"businessDate" gorm:"not null"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	DeletedAt         *time.Time `json:"-"`
	Notes             string     `json:"notes"`
	OrderBillUniqueID string     `json:"orderBillUniqueId"`
	OrderUniqueID     string     `json:"orderUniqueId"`
}
