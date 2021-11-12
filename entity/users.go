package entity

import (
	"time"
)

// Users ...
type Users struct {
	ID                int64      `json:"id" gorm:"primary_key"`
	UserName          string     `json:"userName" gorm:"unique;not null"`
	Name              string     `json:"name"`
	Phone             string     `json:"phone"`
	Address           string     `json:"address"`
	PlatNo            string     `json:"platNo"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	DeletedAt         *time.Time `json:"-"`
	ProfilePicture    string     `json:"profilePicture"`
	RegistrationToken string     `json:"-"`
}
