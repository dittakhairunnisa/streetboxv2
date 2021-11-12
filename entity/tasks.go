package entity

import "time"

// Tasks ..
//
// status :
//
// 1 -> open
//
// 2 -> ongoing
//
// 3 -> arrived
//
// 4 -> completed
type Tasks struct {
	ID              int64      `json:"id" gorm:"primary_key"`
	Types           string     `json:"types" gorm:"not null"` //REGULAR, NONREGULAR, HOMEVISIT
	MerchantUsersID int64      `json:"merchantUsersId" gorm:"not null"`
	Status          int        `json:"status" gorm:"not null"`
	CreatedAt       time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	DeletedAt       *time.Time `json:"-"`
}
