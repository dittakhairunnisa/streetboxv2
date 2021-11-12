package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/tasksnonreglog"
	"streetbox.id/entity"
)

// TasksNonRegLogRepo ..
type TasksNonRegLogRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) tasksnonreglog.RepoInterface {
	return &TasksNonRegLogRepo{db}
}

// Create ..
func (r *TasksNonRegLogRepo) Create(data *entity.TasksNonregularLog) error {
	if err := r.DB.Create(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created TasksNonRegLog: %+v", data)
	return nil
}

// CreateWithTrx ..
func (r *TasksNonRegLogRepo) CreateWithTrx(
	db *gorm.DB, data *entity.TasksNonregularLog) error {
	if err := db.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created TasksNonRegLog: %+v", data)
	return nil
}
