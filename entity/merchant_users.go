package entity

import "time"

// MerchantUsers ...
type MerchantUsers struct {
	ID                int64      `json:"id" gorm:"primary_key"`
	MerchantID        int64      `json:"merchantId" gorm:"not null"`
	UsersID           int64      `json:"usersId" gorm:"not null"`
	RegistrationToken string     `json:"-"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	DeletedAt         *time.Time `json:"-"`
}
