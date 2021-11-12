package entity

import "time"

// HomevisitSales ..
type HomevisitSales struct {
	ID         int64      `json:"id" gorm:"primary_key"`
	MerchantID int64      `json:"merchantId" gorm:"not null"`
	StartDate  time.Time  `json:"startTime" gorm:"not null"`
	EndDate    time.Time  `json:"endTime" gorm:"not null"`
	Deposit    int64      `json:"deposit" gorm:"not null"`
	Total      int        `json:"total" gorm:"not null"`
	Available  int        `json:"available" gorm:"not null"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	DeletedAt  *time.Time `json:"-"`
}
