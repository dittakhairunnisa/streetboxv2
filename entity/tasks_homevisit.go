package entity

import "time"

// TasksHomevisit ..
type TasksHomevisit struct {
	ID                  int64                `json:"id" gorm:"primary_key"`
	TasksID             int64                `json:"tasksId" gorm:"not null"`
	TrxHomevisitSalesID int64                `json:"salesId" gorm:"not null"`
	CreatedAt           time.Time            `json:"createdAt" gorm:"not null"`
	UpdatedAt           *time.Time           `json:"updatedAt"`
	DeletedAt           *time.Time           `json:"-"`
	Log                 []TasksNonregularLog `json:"log"`
}
