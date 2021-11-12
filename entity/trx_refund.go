package entity

import "time"

// TrxRefund ...
type TrxRefund struct {
	ID        int64      `json:"id"        gorm:"primary_key"`
	Types     string     `json:"type"      gorm:"not null"`
	CreatedAt time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}
