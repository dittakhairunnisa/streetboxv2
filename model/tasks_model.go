package model

import "time"

// ReqCreateTasksRegular ..
type ReqCreateTasksRegular struct {
	TrxSalesID   int64     `form:"trxSalesId" binding:"required"`
	UsersID      int64     `form:"usersId" binding:"required"` //foodtruck
	ScheduleDate time.Time `form:"scheduleDate" binding:"required"  time_format:"2006-01-02"`
}

// ReqCreateTasksHomevisit ..
type ReqCreateTasksHomevisit struct {
	TrxHomevisitSalesID int64 `json:"trxVisitSalesId" binding:"required"`
	UsersID             int64 `json:"usersId" binding:"required"` //foodtruck
}

// ResTasksStatus ..
type ResTasksStatus struct {
	ID        int64     `json:"id"`
	Types     string    `json:"types"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ResMyTasksReg ...
type ResMyTasksReg struct {
	ID           int64     `json:"tasksId"`
	TypesID      int64     `json:"typesId"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	ScheduleDate time.Time `json:"scheduleDate"`
	Latitude     float64   `json:"latParkingSpace"`
	Longitude    float64   `json:"lonParkingSpace"`
	Status       int       `json:"status"`
	Types        string    `json:"types"`
}

// ResMyTasksNonReg ...
type ResMyTasksNonReg struct {
	ID        int64   `json:"tasksId"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latParkingSpace"`
	Longitude float64 `json:"lonParkingSpace"`
	Status    int     `json:"status"`
}

// ReqChangeStatusTasks ..
type ReqChangeStatusTasks struct {
	ID     int64 `json:"tasksId" binding:"required"`
	Status int64 `json:"status" binding:"required"`
}

// ReqTasksRegLog ...
type ReqTasksRegLog struct {
	TasksID          int64   `json:"tasksId" binding:"required"`
	TasksRegularID   int64   `json:"typesId" binding:"required"`
	Latitude         float64 `json:"latitude" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	ParkingSpaceName string  `json:"parkingSpaceName" binding:"required"`
}

// ReqTasksVisitLog ...
type ReqTasksVisitLog struct {
	TasksID          int64   `json:"tasksId" binding:"required"`
	TasksHomevisitID int64   `json:"typesId" binding:"required"`
	Latitude         float64 `json:"latitude" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	CustomerName     string  `json:"customerName" binding:"required"`
}

// ReqTasksNonRegLog ...
type ReqTasksNonRegLog struct {
	TasksID   int64   `json:"tasksId" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Address   string  `json:"address" binding:"required"`
}

// ResTasksTracking ..
type ResTasksTracking struct {
	LogTime   time.Time `json:"logTime"`
	TasksID   int64     `json:"tasksId"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Status    int       `json:"status"`
}

// ReqCreateTasksTracking ...
type ReqCreateTasksTracking struct {
	TasksID   int64   `json:"tasksId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TasksVisit dto tasks + trxHomeVisitSalesId
type TasksVisit struct {
	ID                  int64
	Types               string
	MerchantUsersID     int64
	Status              int
	CreatedAt           time.Time
	UpdatedAt           *time.Time
	DeletedAt           *time.Time
	TrxHomevisitSalesID int64
}
