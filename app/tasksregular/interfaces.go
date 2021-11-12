package tasksregular

import (
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(trx *gorm.DB, data *entity.TasksRegular) error
	FindByID(int64) *entity.TasksRegular
	IsAssigned(int64, time.Time, int64) bool
	MyTasks(usersID int64) *[]model.ResMyTasksReg
	MyTasksByMerchantID(int64) *[]model.ResMyTasksReg
	FindByTasksID(int64) *entity.TasksRegular
	FindByMerchantUsersID(int64) *entity.TasksRegular
	CountBySalesID(int64) int
}
