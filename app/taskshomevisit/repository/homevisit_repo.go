package repository

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/taskshomevisit"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// TasksHomevisitRepo ..
type TasksHomevisitRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) taskshomevisit.RepoInterface {
	return &TasksHomevisitRepo{db}
}

// Create ..
func (r *TasksHomevisitRepo) Create(
	db *gorm.DB, data *entity.TasksHomevisit) error {
	if err := db.Create(data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// MyTasks ..
func (r *TasksHomevisitRepo) MyTasks(usersID int64) *[]model.ResMyTasksReg {
	data := new([]model.ResMyTasksReg)
	r.DB.Select("ts.id, tsv.id as types_id, tv.customer_name as name, tv.address, hv.start_date, "+
		"hv.end_date, hv.start_date as schedule_date, "+
		"tv.latitude, tv.longitude, ts.status, ts.types").Joins("JOIN "+
		"trx_visit tv on t.id = tv.trx_id").Joins("JOIN "+
		"users u on t.users_id = u.id").Joins("JOIN "+
		"trx_homevisit_sales thv on tv.id = thv.trx_visit_id").Joins("JOIN "+
		"homevisit_sales hv on thv.homevisit_sales_id = hv.id").Joins("JOIN "+
		"merchant m on hv.merchant_id = m.id").Joins("JOIN "+
		"tasks_homevisit tsv on thv.id = tsv.trx_homevisit_sales_id").Joins("JOIN "+
		"tasks ts on tsv.tasks_id = ts.id").Joins("JOIN "+
		"merchant_users mu on ts.merchant_users_id = mu.id").
		Table("trx t").Where(
		"mu.users_id = ? and ts.status < ? and "+
			"hv.end_date >= ? and ts.deleted_at is null", usersID, 4, time.Now()).
		Scan(&data)
	return data
}

// IsAssigned ..
func (r *TasksHomevisitRepo) IsAssigned(id int64) bool {
	data := new(entity.TasksHomevisit)
	timeNow := time.Now().Format("2006-01-02")
	r.DB.Find(&data, "trx_homevisit_sales_id = ? and "+
		"TO_CHAR(updated_at,'yyyy-MM-dd') like ?", id, timeNow)
	if data.ID == 0 {
		return false
	}
	return true
}

// MyTasksByMerchantID ..
func (r *TasksHomevisitRepo) MyTasksByMerchantID(id int64) *[]model.ResMyTasksReg {
	data := new([]model.ResMyTasksReg)
	r.DB.Select("t.id, th.id as types_id, u.name, trx.address, hs.start_date, "+
		"hs.end_date, hs.start_date as schedule_date, "+
		"trx.latitude, trx.longitude, t.status, t.types").Joins("JOIN "+
		"trx_homevisit_sales trx on hs.id = trx.homevisit_sales_id").Joins("JOIN "+
		"users u on trx.users_id = u.id").Joins("JOIN "+
		"tasks_homevisit th on trx.id = th.trx_homevisit_sales_id").Joins("JOIN "+
		"tasks t on t.id = th.tasks_id").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").Where(
		"mu.merchant_id = ? and t.status < ? and "+
			"hs.end_date >= ? and t.deleted_at is null", id, 4, time.Now()).
		Table("homevisit_sales hs").Scan(&data)
	return data
}

// FindOne ..
func (r *TasksHomevisitRepo) FindOne(qry *entity.TasksHomevisit) *entity.TasksHomevisit {
	data := new(entity.TasksHomevisit)
	r.DB.Where(qry).Find(&data)
	return data
}
