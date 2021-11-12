package entity

import "time"

// TrxSync ..
type TrxSync struct {
	ID           int64     `json:"id"            gorm:"primary_key"`
	UniqueID     string    `json:"uniqueId"     gorm:"not null"`
	MerchantID   int64     `json:"merchantId"   gorm:"not null"`
	Status       int       `json:"status"        gorm:"not null"`
	BusinessDate time.Time `json:"businessDate" gorm:"not null"`
	SyncDate     time.Time `json:"syncDate"     gorm:"not null"`
	Data         string    `json:"data"          gorm:"not null"`
}
