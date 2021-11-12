package entity

import "time"

// TasksHomevisitLog ..
//
// activity:
// 1 -> Check In
//
// 2 -> Check Out
type TasksHomevisitLog struct {
	ID               int64      `json:"id" gorm:"primary_key"`
	TasksHomevisitID int64      `json:"tasksHomevisitId" gorm:"not null"`
	Activity         int        `json:"activity" gorm:"not null"`
	Latitude         float64    `json:"latitude" gorm:"not null"`
	Longitude        float64    `json:"longitude" gorm:"not null"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	DeletedAt        *time.Time `json:"-"`
}
