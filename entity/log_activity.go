package entity

import "time"

// LogActivity ...
type LogActivity struct {
	LogTime  time.Time `json:"logTime" gorm:"primary_key"`
	Activity string    `json:"activity" gorm:"not null"`
}
