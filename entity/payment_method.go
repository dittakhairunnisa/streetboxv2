package entity

import "time"

// PaymentMethod ...
type PaymentMethod struct {
	ID                int64      `json:"id" gorm:"primary_key"`
	Name              string     `json:"name" gorm:"not null"`
	Types             string     `json:"types" gorm:"not null"`
	PaymentProviderID int64      `json:"paymentProviderId" gorm:"not null"`
	IsActive          bool       `json:"isActive" gorm:"not null"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	DeletedAt         *time.Time `json:"-"`
}
