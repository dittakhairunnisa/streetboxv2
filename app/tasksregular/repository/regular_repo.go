package repository

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/tasksregular"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// TasksRegularRepo ..
type TasksRegularRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) tasksregular.RepoInterface {
	return &TasksRegularRepo{db}
}

// Create ..
func (r *TasksRegularRepo) Create(
	trx *gorm.DB, data *entity.TasksRegular) error {
	var daily int
	row := r.DB.Raw("select * from daily_quota_spot(?,?)",
		data.TrxParkingSpaceSalesID, data.ScheduleDate).Row()
	row.Scan(&daily)
	if daily > 0 {
		if err := trx.Create(data).Error; err != nil {
			log.Printf("ERROR: %s", err.Error())
			trx.Rollback()
			return err
		}
		return nil
	}
	trx.Rollback()
	return errors.New("Daily Quota Total Spot Full")
}

// FindByID ..
func (r *TasksRegularRepo) FindByID(int64) *entity.TasksRegular {
	return nil
}

// IsAssigned ..
func (r *TasksRegularRepo) IsAssigned(
	id int64, schedule time.Time, foodtruckID int64) bool {
	data := new(entity.TasksRegular)
	r.DB.Select("t.*").Joins("JOIN "+
		"tasks_regular tr on t.id = tr.tasks_id").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").
		Where("t.deleted_at is null and tr.trx_parking_space_sales_id = ? and "+
			"tr.schedule_date = ? and mu.users_id = ?",
			id, schedule, foodtruckID).Table("tasks t").
		Scan(&data)
	if data.ID == 0 {
		return false
	}
	return true
}

// MyTasks ..
func (r *TasksRegularRepo) MyTasks(usersID int64) *[]model.ResMyTasksReg {
	data := new([]model.ResMyTasksReg)
	dateNow := time.Now().Format("2006-01-02")
	r.DB.Select("t.id, tg.id as types_id, ps.name, ps.address,  pss.start_date, "+
		"pss.end_date, tg.schedule_date, ps.latitude, ps.longitude, t.status, t.types").Joins("JOIN "+
		"parking_space_sales pss on ps.id = pss.parking_space_id").Joins("JOIN "+
		"trx_parking_space_sales trx on pss.id = trx.parking_space_sales_id").Joins("JOIN "+
		"tasks_regular tg on trx.id = tg.trx_parking_space_sales_id").Joins("JOIN "+
		"tasks t on t.id = tg.tasks_id").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").Where(
		"mu.users_id = ? and t.status < ? and "+
			"pss.end_date >= ? and t.deleted_at is null and "+
			"tg.schedule_date >= ?", usersID, 4, time.Now(), dateNow).
		Table("parking_space ps").Scan(&data)
	return data
}

// MyTasksByMerchantID ..
func (r *TasksRegularRepo) MyTasksByMerchantID(id int64) *[]model.ResMyTasksReg {
	data := new([]model.ResMyTasksReg)
	dateNow := time.Now().Format("2006-01-02")
	r.DB.Select("t.id, tg.id as types_id, ps.name, ps.address,  pss.start_date, "+
		"pss.end_date, tg.schedule_date, ps.latitude, ps.longitude, t.status, t.types").Joins("JOIN "+
		"parking_space_sales pss on ps.id = pss.parking_space_id").Joins("JOIN "+
		"trx_parking_space_sales trx on pss.id = trx.parking_space_sales_id").Joins("JOIN "+
		"tasks_regular tg on trx.id = tg.trx_parking_space_sales_id").Joins("JOIN "+
		"tasks t on t.id = tg.tasks_id").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").Where(
		"mu.merchant_id = ? and t.status < ? and "+
			"pss.end_date >= ? and t.deleted_at is null and "+
			"tg.schedule_date >= ?", id, 4, time.Now(), dateNow).
		Table("parking_space ps").Scan(&data)
	return data
}

// FindByTasksID ..
func (r *TasksRegularRepo) FindByTasksID(id int64) *entity.TasksRegular {
	data := new(entity.TasksRegular)
	r.DB.Find(&data, "tasks_id = ?", id)
	if data.ID == 0 {
		return nil
	}
	return data
}

// FindByMerchantUsersID ..
func (r *TasksRegularRepo) FindByMerchantUsersID(id int64) *entity.TasksRegular {
	data := new(entity.TasksRegular)
	r.DB.Select("tr.*").Joins("JOIN "+
		"tasks_regular tr on t.id = tr.tasks_id").
		Where("t.merchant_users_id = ? and t.status < 4", id).
		Table("tasks t").Scan(&data)
	if data.ID == 0 {
		return nil
	}
	return data
}

// CountBySalesID ..
func (r *TasksRegularRepo) CountBySalesID(id int64) int {
	count := 0
	data := new(entity.TasksRegular)
	timeNow := time.Now().Format("2006-01-02")
	r.DB.Select("tr.*").Joins("JOIN "+
		"tasks_regular_log trl on tr.id = trl.tasks_regular_id").Joins("JOIN "+
		"trx_parking_space_sales tpss on tr.trx_parking_space_sales_id = tpss.id").Joins("JOIN "+
		"parking_space_sales pss on tpss.parking_space_sales_id = pss.id").
		Where("pss.id = ? and "+
			"TO_CHAR(trl.updated_at,'yyyy-mm-dd') like ? and trl.activity = ?",
			id, timeNow, 1).Table("tasks_regular tr").Scan(&data).Count(&count)
	return count
}
