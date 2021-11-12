package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/tasksreglog"
	"streetbox.id/entity"
)

// TasksRegLogRepo ..
type TasksRegLogRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) tasksreglog.RepoInterface {
	return &TasksRegLogRepo{db}
}

// Create ..
func (r *TasksRegLogRepo) Create(data *entity.TasksRegularLog) error {
	if err := r.DB.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}
