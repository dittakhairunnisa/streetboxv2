package entity

import "time"

// TasksNonregularLog ..
//
// activity:
// 1 -> Check In
//
// 2 -> Check Out
type TasksNonregularLog struct {
	ID                int64      `json:"id" gorm:"primary_key"`
	TasksNonregularID int64      `json:"tasksNonregularId" gorm:"not null"`
	Activity          int        `json:"activity" gorm:"not null"`
	Latitude          float64    `json:"latitude" gorm:"not null"`
	Longitude         float64    `json:"longitude" gorm:"not null"`
	Address           string     `json:"address"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	DeletedAt         *time.Time `json:"-"`
}
