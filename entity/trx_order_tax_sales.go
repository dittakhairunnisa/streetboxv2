package entity

import "time"

// TrxOrderTaxSales ...
type TrxOrderTaxSales struct {
	ID                int64      `json:"id"  gorm:"primary_key"`
	UniqueID          string     `json:"uniqueId"`
	TrxOrderBillID    int64      `json:"trxOrderBillId" gorm:"not null"`
	Name              string     `json:"name" gorm:"not null"`
	MerchantTaxID     int64      `json:"merchantTaxId" gorm:"not null"`
	Amount            float64    `json:"amount" gorm:"not null"`
	Types             int        `json:"types" gorm:"not null"` // 0: Exclusive, 1: Inclusive
	OrderBillUniqueID string     `json:"orderBillUniqueId"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	OrderUniqueID     string     `json:"orderUniqueId"`
}
