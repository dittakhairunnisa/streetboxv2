package entity

import "time"

// TrxRefundVisit ...
type TrxRefundVisit struct {
	ID                  int64      `json:"id"                     gorm:"primary_key"`
	TrxHomevisitSalesID int64      `json:"trxHomevisitSalesId"    gorm:"not null"`
	Amount              int64      `json:"amount"                 gorm:"not null"`
	TrxRefundID         int64      `json:"trxRefundId"            gorm:"not null"`
	CreatedAt           time.Time  `json:"createdAt"              gorm:"not null"`
	UpdatedAt           *time.Time `json:"updatedAt"`
	DeletedAt           *time.Time `json:"-"`
}
