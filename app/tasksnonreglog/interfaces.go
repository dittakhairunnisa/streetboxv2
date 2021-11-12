package tasksnonreglog

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.TasksNonregularLog) error
	CreateWithTrx(*gorm.DB, *entity.TasksNonregularLog) error
}
