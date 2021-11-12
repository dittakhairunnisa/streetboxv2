package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/tasksvisitlog"
	"streetbox.id/entity"
)

// TasksVisitLogRepo ..
type TasksVisitLogRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) tasksvisitlog.RepoInterface {
	return &TasksVisitLogRepo{db}
}

// Create ..
func (r *TasksVisitLogRepo) Create(data *entity.TasksHomevisitLog) error {
	if err := r.DB.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}
