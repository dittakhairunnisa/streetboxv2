package entity

import "time"

// PaymentProvier ..
type PaymentProvier struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	Name      string     `json:"name" gorm:"not null"`
	CreatedAt time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}
