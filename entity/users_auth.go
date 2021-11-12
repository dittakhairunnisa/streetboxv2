package entity

import "time"

// UsersAuth ...
type UsersAuth struct {
	UserName  string `gorm:"primary_key"`
	Password  string `json:"-"`
	DeviceID  string
	CreatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
}
