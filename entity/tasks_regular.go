package entity

import "time"

// TasksRegular ..
type TasksRegular struct {
	ID                     int64             `json:"id" gorm:"primary_key"`
	TasksID                int64             `json:"tasksId" gorm:"not null"`
	TrxParkingSpaceSalesID int64             `json:"salesId" gorm:"not null"`
	ScheduleDate           time.Time         `json:"scheduleDate" gorm:"not null"` //date
	CreatedAt              time.Time         `json:"createdAt" gorm:"not null"`
	UpdatedAt              *time.Time        `json:"updatedAt"`
	DeletedAt              *time.Time        `json:"-"`
	Log                    []TasksRegularLog `json:"log"`
}
