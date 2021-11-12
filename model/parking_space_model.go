package model

import (
	"time"

	"github.com/lib/pq"
)

// ReqParkingSpaceCreate ... from recruit apps
type ReqParkingSpaceCreate struct {
	Name          string    `json:"name" binding:"required"`
	Address       string    `json:"address" binding:"required"`
	Latitude      float64   `json:"latitude" binding:"required"`
	Longitude     float64   `json:"longitude" binding:"required"`
	TotalSpace    int       `json:"total" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	LandlordInfo  string    `json:"landlordInfo" binding:"required"`
	Rating        float32   `json:"rating" binding:"required"`
	StartContract time.Time `json:"startContract" binding:"required" time_format:"2006-01-02"`
	EndContract   time.Time `json:"endContract" binding:"required" time_format:"2006-01-02"`
	StartTime     time.Time `json:"startTime" binding:"required"`
	EndTime       time.Time `json:"endTime" binding:"required"`
	City          string    `form:"city" binding:"required"`
}

// ReqParkingSpaceUpdate ...
type ReqParkingSpaceUpdate struct {
	Name          string    `form:"name,omitempty"`
	Address       string    `form:"address,omitempty"`
	Latitude      float64   `form:"latitude,omitempty"`
	Longitude     float64   `form:"longitude,omitempty"`
	TotalSpace    int       `form:"total,omitempty"`
	Description   string    `form:"description,omitempty"`
	LandlordInfo  string    `form:"landlordInfo,omitempty"`
	Rating        float32   `form:"rating,omitempty"`
	StartContract time.Time `form:"startContract,omitempty" time_format:"2006-01-02"` //
	EndContract   time.Time `form:"endContract,omitempty" time_format:"2006-01-02"`
	StartTime     time.Time `form:"startTime,omitempty" time_format:"2006-01-02 15:04:05"`
	EndTime       time.Time `form:"endTime,omitempty" time_format:"2006-01-02 15:04:05"`
	City          string    `form:"city,omitempty"`
}

// ReqSalesCreate ...
type ReqSalesCreate struct {
	StartDate      time.Time `form:"startDate" binding:"required" time_format:"2006-01-02 15:04:05"`
	EndDate        time.Time `form:"endDate" binding:"required" time_format:"2006-01-02 15:04:05"`
	TotalSlot      int       `form:"totalSlot" binding:"required"`
	Point          int64     `form:"point" binding:"required"`
	ParkingSpaceID int64     `form:"parkingSpaceId" binding:"required"`
}

// ReqGetSpaceBySalesDate ...
type ReqGetSpaceBySalesDate struct {
	StartDate time.Time `form:"startDate" binding:"required" time_format:"2006-01-02 15:04:05"`
	EndDate   time.Time `form:"endDate"   binding:"required" time_format:"2006-01-02 15:04:05"`
}

// ReqSalesUpdate ...
type ReqSalesUpdate struct {
	ParkingSpaceID int   `json:"parkingSpaceId,omitempty"`
	TotalSlot      int   `json:"totalSlot,omitempty"`
	AvailableSlot  int   `json:"availableSlot,omitempty"`
	Point          int64 `json:"point,omitempty"`
}

// ResMyParking ...
type ResMyParking struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	ImagesMeta  pq.StringArray `json:"images" swaggertype:"array,string"`
	Description string         `json:"description"`
	Rating      float32        `json:"rating"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	StartTime   time.Time      `json:"startTime"`
	EndTime     time.Time      `json:"endTime"`
}

// ResMyParkingList contains trx visit and parking space
type ResMyParkingList struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	Address         string         `json:"address"`
	ImagesMeta      pq.StringArray `json:"images" swaggertype:"array,string"`
	Description     string         `json:"description"`
	Rating          float32        `json:"rating"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	StartTime       time.Time      `json:"startTime"`
	EndTime         time.Time      `json:"endTime"`
	ProfilePicture  string         `json:"profilePicture"`
	TrxVisitSalesID int64          `json:"trxVisitSalesId"`
}

// ResSlotMyParking ...
type ResSlotMyParking struct {
	ID        int64     `json:"id"`
	TotalSlot int       `json:"totalSlot"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// ReqDeleteAsset req body delete image and documents
type ReqDeleteAsset struct {
	Filename string `json:"filename" binding:"required"`
}

// ReqCreateFoodtruck ..
type ReqCreateFoodtruck struct {
	UserName string `json:"userName" binding:"required"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	PlatNo   string `json:"platNo"`
}

// ResSales ..
type ResSales struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	StartDate      time.Time  `json:"startDate"`
	EndDate        time.Time  `json:"endDate"`
	TotalSlot      int        `json:"totalSlot"`
	AvailableSlot  int        `json:"availableSlot"`
	Point          int64      `json:"point"`
	ParkingSpaceID int64      `json:"parkingSpaceId"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}
