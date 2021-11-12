package entity

import "time"

//TrxVisit ...
type TrxVisit struct {
	ID              int64      `json:"id" gorm:"primary_key"`
	TrxID           string     `json:"trxId" gorm:"not null"`
	GrandTotal      float64    `json:"grandTotal" gorm:"not null"`
	PaymentMethodID int64      `json:"paymentMethodId" gorm:"not null"`
	CustomerName    string     `json:"customerName" gorm:"not null"`
	Address         string     `json:"address" gorm:"not null"`
	Phone           string     `json:"phone"`
	Notes           string     `json:"notes"`
	Longitude       float64    `json:"longitude" gorm:"not null"`
	Latitude        float64    `json:"latitude" gorm:"not null"`
	CreatedAt       time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	DeletedAt       *time.Time `json:"-"`
}
