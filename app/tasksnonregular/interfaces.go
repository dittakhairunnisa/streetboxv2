package tasksnonregular

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(trx *gorm.DB, data *entity.TasksNonregular) error
	GetTasksByUsersID(int64) *model.ResMyTasksNonReg
	GetByTasksID(int64) *entity.TasksNonregular
	MyTasksByMerchantID(int64) *[]model.ResMyTasksReg
}
