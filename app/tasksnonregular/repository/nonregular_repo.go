package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/tasksnonregular"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// TasksNonRegRepo ..
type TasksNonRegRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) tasksnonregular.RepoInterface {
	return &TasksNonRegRepo{db}
}

// Create ..
func (r *TasksNonRegRepo) Create(
	trx *gorm.DB, data *entity.TasksNonregular) error {
	if err := trx.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		trx.Rollback()
		return err
	}
	return nil
}

// GetTasksByUsersID ..
func (r *TasksNonRegRepo) GetTasksByUsersID(usersID int64) *model.ResMyTasksNonReg {
	data := new(model.ResMyTasksNonReg)
	r.DB.Select("t.id, tnl.address, tnl.latitude, tnl.longitude, "+
		"t.status").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").Joins("JOIN "+
		"tasks_nonregular tn on t.id = tn.tasks_id").Joins("JOIN "+
		"tasks_nonregular_log tnl on tn.id = tnl.tasks_nonregular_id").
		Where("mu.users_id = ? and t.deleted_at is null", usersID).
		Order("tnl.id desc").Table("tasks t").Limit(1).Scan(&data)
	if data.ID == 0 || data.Status == 4 {
		return nil
	}
	return data
}

// GetByTasksID ..
func (r *TasksNonRegRepo) GetByTasksID(id int64) *entity.TasksNonregular {
	data := new(entity.TasksNonregular)
	r.DB.Find(&data, "tasks_id = ?", id)
	if data.ID == 0 {
		return nil
	}
	return data
}

// MyTasksByMerchantID ..
func (r *TasksNonRegRepo) MyTasksByMerchantID(id int64) *[]model.ResMyTasksReg {
	data := new([]model.ResMyTasksReg)
	r.DB.Select("distinct t.id, tn.id as types_id, u.plat_no as name, tnl.address, tn.created_at as start_date, "+
		"t.updated_at as end_date, tn.created_at as schedule_date, "+
		"tnl.latitude, tnl.longitude, "+
		"t.status, t.types").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").Joins("JOIN "+
		"users u on mu.users_id = u.id").Joins("JOIN "+
		"tasks_nonregular tn on t.id = tn.tasks_id").Joins("JOIN "+
		"tasks_nonregular_log tnl on tn.id = tnl.tasks_nonregular_id").
		Where("mu.merchant_id = ? and t.status < ? and t.deleted_at is null", id, 4).
		Order("t.id desc").Table("tasks t").Scan(&data)
	return data
}
