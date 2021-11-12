package entity

import "time"

// TrxParkingSpaceSales ...
type TrxParkingSpaceSales struct {
	ID                  int64      `json:"id" gorm:"primary_key"`
	ParkingSpaceSalesID int64      `json:"parkingSpaceSalesId" gorm:"not null"`
	MerchantID          int64      `json:"merchantId" gorm:"not null"`
	TotalSlot           int        `json:"totalSlot" gorm:"not null"`
	CreatedAt           time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt           *time.Time `json:"updatedAt"`
	DeletedAt           *time.Time `json:"-"`
}
