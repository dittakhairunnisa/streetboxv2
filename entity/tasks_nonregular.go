package entity

import "time"

// TasksNonregular ..
type TasksNonregular struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	TasksID   int64      `json:"tasksId" gorm:"not null"`
	CreatedAt time.Time  `json:"createdAt" gorm:"not null"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}
