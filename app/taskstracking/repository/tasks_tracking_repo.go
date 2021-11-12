package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/taskstracking"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TasksTrackingRepo ..
type TasksTrackingRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) taskstracking.RepoInterface {
	return &TasksTrackingRepo{db}
}

// Create ..
func (r *TasksTrackingRepo) Create(data *entity.TasksTracking) error {
	if err := r.DB.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	// log.Printf("INFO: Created TasksTracking: %+v", data)
	return nil
}

// GetTrackingByID need changes because tracking should have status 2
func (r *TasksTrackingRepo) GetTrackingByID(id int64) *model.ResTasksTracking {
	data := new(model.ResTasksTracking)
	dataTask := new(model.ResTasksStatus)
	r.DB.Select("tt.*,t.status").Joins("JOIN "+
		"tasks_tracking tt on t.id = tt.tasks_id").Joins("JOIN "+
		"merchant_users mu on t.merchant_users_id = mu.id").
		Where("mu.id = (select merchant_users_id from tasks t2 where t2.id = ?)", id).Table("tasks t").Order("tt.log_time desc").Limit("1").Scan(&data)
	r.DB.Select("t.*").
		Where("t.id = ?", id).Table("tasks t").Scan(&dataTask)
	data.TasksID = id
	data.Status = util.ParamIDToInt(dataTask.Status)
	return data
}

// DeleteByTasksID Hard Delete for tasks completed
func (r *TasksTrackingRepo) DeleteByTasksID(id int64) {
	r.DB.Delete(entity.TasksTracking{}, "tasks_id = ?", id)
}

// GetLiveTracking consumer get live tracking foodtruck on going (tasks.status = 2) in map
func (r *TasksTrackingRepo) GetLiveTracking(lat, lon, distance float64) *[]model.ResLiveTracking {
	data := make([]model.ResLiveTracking, 0)
	rows, _ := r.DB.Raw("select * from livetracking(?,?,?) order by nearby asc", lat, lon, distance).Rows()
	for rows.Next() {
		var dat model.ResLiveTracking
		r.DB.ScanRows(rows, &dat)
		data = append(data, dat)
	}
	return &data
}

// DeleteAll method for clear all records in tasks_tracking daily at 00:00:01
func (r *TasksTrackingRepo) DeleteAll() {
	r.DB.Exec("delete from tasks_tracking;")
	log.Printf("INFO: Clear All Tasks Tracking")
}
