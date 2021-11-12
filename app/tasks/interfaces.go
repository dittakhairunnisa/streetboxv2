package tasks

import (
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.Tasks) (*gorm.DB, error)
	FindOne(id int64) *entity.Tasks
	FindNonRegByUsersID(int64) *entity.Tasks
	FindByUsersID(int64) *entity.Tasks
	UpdateByID(*entity.Tasks, int64) error
	UpdateByMerchantUsersID(*entity.Tasks, int64) error
	FindUncompletedRegByMerchantUsersID(int64) *[]entity.Tasks
	FindUncompletedVisitByMerchantUsersID(int64) *[]model.TasksVisit
	FindUncompletedNonRegByMerchantUsersID(int64) *[]entity.Tasks
	FindByStatusUsersID(usersID int64, status int) *entity.Tasks
}

// ServiceInterface ..
type ServiceInterface interface {
	CreateTasksRegular(*model.ReqCreateTasksRegular) (*entity.Tasks, error)
	CreateTasksHomevisit(*model.ReqCreateTasksHomevisit) (*entity.Tasks, error)
	CreateTasksNonRegular(usersID int64) (*entity.Tasks, error)
	IsTasksRegularAssigned(salesID int64, schedudleDate time.Time, foodtruckID int64) bool
	IsTasksHomevisitAssigned(trxVisitSalesID, foodtruckID int64) bool
	MyTasksRegByUsersID(usersID int64) *[]model.ResMyTasksReg
	MyTasksNonRegByUsersID(int64) *model.ResMyTasksNonReg
	CreateRegLog(*model.ReqTasksRegLog, int) (*entity.TasksRegularLog, error)
	CreateVisitLog(*model.ReqTasksVisitLog, int) (*entity.TasksHomevisitLog, error)
	CreateNonRegLog(*model.ReqTasksNonRegLog, int) (*entity.TasksNonregularLog, error)
	UpdateTasksStatus(tasksID int64, status int) error
	CreateTracking(*model.ReqCreateTasksTracking) (*model.ResTasksTracking, error)
	GetTasksByID(int64) *entity.Tasks
	GetTrackingByTasksID(int64) *model.ResTasksTracking
	GetTasksNonRegByUsersID(int64) *entity.Tasks
	GetAllByMerchantID(int64) *[]model.ResMyTasksReg
	RegularStatusCompletedByUsersID(int64)
	VisitStatusCompletedByUsersID(int64)
	NonRegStatusCompletedByUsersID(int64)
	ClosedTrxHomevisitSales(tasksHomevisitID int64) error
}
