package entity

import "time"

// ParkingSpaceSales ... Transaction
type ParkingSpaceSales struct {
	ID             int64      `json:"id" gorm:"primary_key"`
	StartDate      time.Time  `json:"startDate" gorm:"not null;type:date"` // date
	EndDate        time.Time  `json:"endDate" gorm:"not null;type:date"`   // date
	TotalSlot      int        `json:"totalSlot" gorm:"not null"`
	AvailableSlot  int        `json:"availableSlot" gorm:"not null"`
	Point          int64      `json:"point" gorm:"not null"`
	ParkingSpaceID int64      `json:"parkingSpaceId" gorm:"not null"` // foreign key
	CreatedAt      time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	DeletedAt      *time.Time `json:"-"`
}
