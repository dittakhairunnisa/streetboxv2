package entity

import (
	"time"
)

// TrxOrder relation from Trx parent ...
type TrxOrder struct {
	ID                int64      `json:"id" gorm:"primary_key"`
	OrderNo           string     `json:"orderNo"`
	BillNo            string     `json:"billNo"`
	IsClose           bool       `json:"isClose" gorm:"not null"`
	Note              string     `json:"note"`
	Types             int        `json:"types" gorm:"not null"` // 0: Offline 1: Online
	MerchantUsersID   int64      `json:"merchantUsersId" gorm:"not null"`
	BusinessDate      time.Time  `json:"businessDate" gorm:"not null"`
	TotalDiscount     float64    `json:"totalDiscount" gorm:"not null"`
	GrandTotal        float64    `json:"grandTotal" gorm:"not null"`
	TrxID             string     `json:"trxId" gorm:"not null"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	DeletedAt         *time.Time `json:"-"`
	UniqueID          string     `json:"uniqueId"`
	PaymentMethodId   string     `json:"paymentMethodId"`
	PaymentMethodName string     `json:"paymentMethodName"`
	TypeOrder         string     `json:"typeOrder"`
	TypePayment       string     `json:"typePayment"`
}
