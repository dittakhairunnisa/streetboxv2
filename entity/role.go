package entity

import "time"

// Role ...
type Role struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	Name        string     `json:"name" gorm:"unique;not null"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"-"`
}
