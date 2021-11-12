package taskshomevisit

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(db *gorm.DB, data *entity.TasksHomevisit) error
	MyTasks(usersID int64) *[]model.ResMyTasksReg
	IsAssigned(salesID int64) bool
	MyTasksByMerchantID(int64) *[]model.ResMyTasksReg
	FindOne(*entity.TasksHomevisit) *entity.TasksHomevisit
}
