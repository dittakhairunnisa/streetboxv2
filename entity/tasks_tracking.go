package entity

import "time"

// TasksTracking temporary table for live tracking using REST API
// This table will be Hard Delete after task status completed
type TasksTracking struct {
	LogTime   time.Time `json:"logTime" gorm:"not null"`
	TasksID   int64     `json:"tasksId" gorm:"not null"`
	Latitude  float64   `json:"latitude" gorm:"not null"`
	Longitude float64   `json:"longitude" gorm:"not null"`
}
