package entity

import "time"

// Merchant ...
type Merchant struct {
	ID         int64      `json:"id" gorm:"primary_key"`
	XenditID   string     `json:"xendit_id"`
	Name       string     `json:"name" gorm:"not null"`
	Address    string     `json:"address" gorm:"not null"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	Logo       string     `json:"logo"`
	Banner     string     `json:"banner"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	DeletedAt  *time.Time `json:"-"`
	City       string     `json:"city"`
	IGAccount  string     `json:"igAccount"`
	CategoryID int64      `json:"categoryID"`
	Terms      string     `json:"terms"`
}
