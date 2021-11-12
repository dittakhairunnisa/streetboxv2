package entity

import "time"

// LogXenditRequest ..
type LogXenditRequest struct {
	ID           int64     `json:"id" gorm:"primary_key"`
	TrxID        string    `json:"trxId" gorm:"not null"`
	RequestData  string    `json:"requestData" gorm:"not null"`
	ResponseData string    `json:"responseDate" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"not null"`
}
