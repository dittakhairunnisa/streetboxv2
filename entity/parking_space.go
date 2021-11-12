package entity

import (
	"time"

	"github.com/lib/pq"
)

// ParkingSpace ... Master
type ParkingSpace struct {
	ID            int64          `json:"id" gorm:"primary_key"`
	Name          string         `json:"name" gorm:"not null"`
	Address       string         `json:"address" gorm:"not null"`
	Latitude      float64        `json:"latitude" gorm:"not null"`
	Longitude     float64        `json:"longitude" gorm:"not null"`
	Description   string         `json:"description" gorm:"not null"`
	TotalSpace    int            `json:"totalSpace" gorm:"not null"`
	StartTime     time.Time      `json:"startTime" gorm:"not null"`                                   // start operational time
	EndTime       time.Time      `json:"endTime" gorm:"not null"`                                     // end operational time
	Rating        float32        `json:"rating" gorm:"not null;default:'0.0'"`                        // superadmin
	ImagesMeta    pq.StringArray `json:"imagesMeta" gorm:"type:text[]" swaggertype:"array,string"`    // meta data binary image
	DocumentsMeta pq.StringArray `json:"documentsMeta" gorm:"type:text[]" swaggertype:"array,string"` // meta data binary file
	CreatedAt     time.Time      `json:"createdAt" gorm:"not null"`
	UpdatedAt     *time.Time     `json:"updatedAt"`
	DeletedAt     *time.Time     `json:"-"`
	LandlordInfo  string         `json:"landlordInfo" gorm:"not null"`
	StartContract time.Time      `json:"startContract" gorm:"not null" time_format:"2006-01-02"`
	EndContract   time.Time      `json:"endContract" gorm:"not null"  time_format:"2006-01-02"`
	City          string         `json:"city"`
}
