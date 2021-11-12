package entity

import "time"

// TrxRefundSpace ...
type TrxRefundSpace struct {
	ID                     int64      `json:"id"                     gorm:"primary_key"`
	TrxParkingSpaceSalesID int64      `json:"trxParkingSpaceSalesId"   gorm:"not null"`
	Amount                 int64      `json:"amount"                 gorm:"not null"`
	TrxRefundID            int64      `json:"trxRefundId"            gorm:"not null"`
	CreatedAt              time.Time  `json:"createdAt"              gorm:"not null"`
	UpdatedAt              *time.Time `json:"updatedAt"`
	DeletedAt              *time.Time `json:"-"`
}
