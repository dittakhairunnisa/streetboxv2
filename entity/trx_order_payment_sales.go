package entity

import "time"

// TrxOrderPaymentSales ..
type TrxOrderPaymentSales struct {
	ID                int64      `json:"id" gorm:"primary_key"`
	TrxOrderBillID    int64      `json:"trxOrderBillId" gorm:"not null"`
	Name              string     `json:"name" gorm:"not null"`
	Amount            float64    `json:"amount" gorm:"not null"`
	UniqueID          string     `json:"uniqueId" gorm:"not null"`
	OrderBillUniqueID string     `json:"orderBillUniqueId"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	PaymentMethodID   int64      `json:"paymentMethodId" gorm:"not null"`
	OrderUniqueID     string     `json:"orderUniqueId"`
}
