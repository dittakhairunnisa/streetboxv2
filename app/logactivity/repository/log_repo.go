package repository

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/logactivity"
	"streetbox.id/entity"
	"streetbox.id/util"
)

// LogRepo ..
type LogRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) logactivity.RepoInterface {
	return &LogRepo{db}
}

// Create ..
func (r *LogRepo) Create(activity string) {
	data := new(entity.LogActivity)
	data.LogTime = time.Now()
	data.Activity = activity
	if err := r.DB.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	log.Printf("INFO: Created LogActivity: %+v", data)
}

// GetAllPagination ..
func (r *LogRepo) GetAllPagination(limit int, page int, sort []string) (*[]entity.LogActivity, int, int) {
	data := new([]entity.LogActivity)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Limit(limit).Offset(offset)
	//sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	} else {
		qry = qry.Order("log_time desc")
	}
	qry = qry.Find(&data)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}

// GetList ..
func (r *LogRepo) GetList() *[]entity.LogActivity {
	data := new([]entity.LogActivity)
	if err := r.DB.Order("log_time desc").Find(&data).Error; err != nil {
		return nil
	}
	return data
}
