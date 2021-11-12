package entity

import "time"

// TrxHomevisitSales ..
type TrxHomevisitSales struct {
	ID               int64      `json:"id" gorm:"primary_key"`
	HomevisitSalesID int64      `json:"homeVisitSalesId" gorm:"not null"`
	Total            int64      `json:"total" gorm:"not null"`
	TrxVisitID       int64      `json:"visitId" gorm:"not null"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	DeletedAt        *time.Time `json:"-"`
	Status           string     `json:"status" gorm:"not null"` // CLOSED, OPEN
}
