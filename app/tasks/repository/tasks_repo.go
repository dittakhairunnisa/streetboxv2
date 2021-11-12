package repository

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/tasks"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TasksRepo ..
type TasksRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) tasks.RepoInterface {
	return &TasksRepo{db}
}

// FindOne get by id
func (r *TasksRepo) FindOne(id int64) *entity.Tasks {
	data := new(entity.Tasks)
	r.DB.Find(&data, "id = ?", id)
	if data.ID == 0 {
		return nil
	}
	return data
}

// Create ..
func (r *TasksRepo) Create(data *entity.Tasks) (*gorm.DB, error) {
	db := r.DB.Begin()
	if err := db.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		db.Rollback()
		return nil, err
	}
	log.Printf("INFO: Created Tasks: %+v", data)
	return db, nil
}

// FindNonRegByUsersID ..
func (r *TasksRepo) FindNonRegByUsersID(id int64) *entity.Tasks {
	data := new(entity.Tasks)
	r.DB.Select("t.*").Joins("JOIN merchant_users mu "+
		"on t.merchant_users_id = mu.id").
		Where("mu.users_id = ? and t.deleted_at is null "+
			"and t.types = ?", id, "NONREGULAR").
		Order("t.updated_at desc").Table("tasks t").Limit(1).Scan(&data)
	if data.ID == 0 || data.Status == 4 {
		return nil
	}
	return data
}

// FindByUsersID ..
func (r *TasksRepo) FindByUsersID(id int64) *entity.Tasks {
	data := new(entity.Tasks)
	r.DB.Select("t.*").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").
		Where("mu.users_id = ? and t.deleted_at is null", id).
		Order("t.id desc").Table("tasks t").Limit(1).Scan(&data)
	if data.ID == 0 {
		return nil
	}
	return data
}

// UpdateByID ..
func (r *TasksRepo) UpdateByID(data *entity.Tasks, id int64) error {
	if err := r.DB.Model(&entity.Tasks{ID: id}).Updates(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Updated Tasks: %+v", data)
	return nil
}

// UpdateByMerchantUsersID ..
func (r *TasksRepo) UpdateByMerchantUsersID(data *entity.Tasks, id int64) error {
	if err := r.DB.Model(&entity.Tasks{MerchantUsersID: id}).
		Updates(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Updated Tasks: %+v", data)
	return nil
}

// FindUncompletedRegByMerchantUsersID ..
func (r *TasksRepo) FindUncompletedRegByMerchantUsersID(id int64) *[]entity.Tasks {
	data := new([]entity.Tasks)
	dateNow := time.Now().Format("2006-01-02")
	r.DB.Select("t.*").Joins("JOIN "+
		"tasks_regular tr on t.id = tr.tasks_id").Table("tasks t").
		Where("t.deleted_at is null and t.status < 4 and t.types = ? "+
			"and tr.schedule_date < ? and t.merchant_users_id = ?",
			util.TasksRegular, dateNow, id).Scan(&data)
	return data
}

// FindUncompletedVisitByMerchantUsersID ..
func (r *TasksRepo) FindUncompletedVisitByMerchantUsersID(id int64) *[]model.TasksVisit {
	data := new([]model.TasksVisit)
	r.DB.Select("t.*,th.trx_homevisit_sales_id").Joins("JOIN "+
		"tasks_homevisit th on t.id = th.tasks_id").Joins("JOIN "+
		"trx_homevisit_sales ths on th.trx_homevisit_sales_id = ths.id").Joins("JOIN "+
		"trx_visit tv on ths.trx_visit_id = tv.id").Joins("JOIN "+
		"homevisit_sales hs on ths.homevisit_sales_id = hs.id").
		Table("tasks t").
		Where("t.deleted_at is null and t.status < 4 and t.types = ? "+
			"and hs.end_date < ? and t.merchant_users_id = ? and ths.status = ?",
			util.TasksHomeVisit, time.Now(), id, util.TrxVisitStatusOpen).Scan(&data)
	return data
}

// FindUncompletedNonRegByMerchantUsersID ..
func (r *TasksRepo) FindUncompletedNonRegByMerchantUsersID(id int64) *[]entity.Tasks {
	data := new([]entity.Tasks)
	r.DB.Find(&data, "status < ? and types = ? and merchant_users_id = ?",
		4, util.TasksNonRegular, id)
	return data
}

// FindByStatusUsersID get tasks for get all foodtrucks
func (r *TasksRepo) FindByStatusUsersID(usersID int64, status int) *entity.Tasks {
	data := new(entity.Tasks)
	r.DB.Select("t.*").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").
		Table("tasks t").Where("t.deleted_at is null and "+
		"t.status = ? and mu.users_id = ?", status, usersID).Order("t.id desc").
		Limit(1).Scan(&data)
	return data
}
