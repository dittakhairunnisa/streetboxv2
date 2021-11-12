package model

import "time"

// ReqCreateTask ...
type ReqCreateTask struct {
	TrxParkingSpaceSalesID int64 `json:"trxId" binding:"required"`
	UsersID                int64 `json:"usersId" binding:"required"` //Foodtruck
}

// ResMyTask ...
type ResMyTask struct {
	ID        int64     `json:"taskId"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Latitude  float64   `json:"latParkingSpace"`
	Longitude float64   `json:"lonParkingSpace"`
	Status    int       `json:"status"`
}

// ReqTaskLog ...
type ReqTaskLog struct {
	TaskID    int64   `json:"taskId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ResTaskBySales ...
type ResTaskBySales struct {
	ID        int64   `json:"taskId"`
	OpsName   string  `json:"foodtruckName"`
	Latitude  float64 `json:"latParkingSpace"`
	Longitude float64 `json:"lonParkingSpace"`
	Status    int     `json:"status"`
	Address   string  `json:"address"`
}

// ReqCreateTracking ...
type ReqCreateTracking struct {
	TaskID    int64   `json:"taskId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ResTracking ..
type ResTracking struct {
	LogTime   time.Time `json:"logTime"`
	TaskID    int64     `json:"taskId"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Status    int       `json:"status"`
}
