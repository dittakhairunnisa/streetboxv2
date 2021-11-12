package entity

import (
	"time"
)

// LogActivityMerchant ..
type LogActivityMerchant struct {
	ID         int64     `json:"id" gorm:"primary_key"`
	MerchantID int64     `json:"merchantId" gorm:"not null"`
	LogTime    time.Time `json:"time" gorm:"not null"`
	Activity   string    `json:"activity" gorm:"not null"`
}
