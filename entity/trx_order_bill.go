package entity

import "time"

// TrxOrderBill ..
type TrxOrderBill struct {
	ID            int64      `json:"id" gorm:"primary_key"`
	BillNo        string     `json:"billNo"`
	IsClose       bool       `json:"isClose" gorm:"not null"`
	TotalDiscount float32    `json:"totalDiscount" gorm:"not null"`
	SubTotal      float64    `json:"subTotal" gorm:"not null"`
	TotalTax      float64    `json:"totalTax" gorm:"not null"`
	GrandTotal    float64    `json:"grandTotal" gorm:"not null"`
	BusinessDate  time.Time  `json:"businessDate" gorm:"not null"`
	TrxOrderID    int64      `json:"trxOrderId" gorm:"not null"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	DeletedAt     *time.Time `json:"-"`
	UniqueID      string     `json:"uniqueId"`
	OrderUniqueID string     `json:"orderUniqueId"`
}
