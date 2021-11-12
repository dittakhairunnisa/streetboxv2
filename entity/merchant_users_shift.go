package entity

import "time"

// MerchantUsersShift ...
type MerchantUsersShift struct {
	ID              int64      `json:"id" gorm:"primary_key"`
	MerchantUsersID int64      `json:"merchantUsersId" gorm:"not null"`
	Shift           string     `json:"shift" gorm:"not null"`
	CreatedAt       time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	DeletedAt       *time.Time `json:"-"`
}
