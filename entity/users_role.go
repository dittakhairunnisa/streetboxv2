package entity

import "time"

// UsersRole ...
type UsersRole struct {
	ID        int64 `gorm:"primary_key"`
	UsersID   int64
	RoleID    int64
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt *time.Time
	DeletedAt *time.Time `json:"-"`
}
