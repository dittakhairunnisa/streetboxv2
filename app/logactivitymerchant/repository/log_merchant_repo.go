package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/logactivitymerchant"
	"streetbox.id/entity"
	"streetbox.id/util"
)

// LogActivityMerchant ..
type LogActivityMerchant struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) logactivitymerchant.RepoInterface {
	return &LogActivityMerchant{db}
}

// Create ..
func (r *LogActivityMerchant) Create(data *entity.LogActivityMerchant) error {
	if err := r.DB.Create(data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return nil
}

// GetAll ..
func (r *LogActivityMerchant) GetAll(
	limit, page int, sort []string, merchantID int64) (*[]entity.LogActivityMerchant, int, int) {
	data := new([]entity.LogActivityMerchant)
	count := 0
	offset := util.Offset(page, limit)
	qry := r.DB.Limit(limit).Offset(offset)
	// sorting
	if len(sort) > 0 {
		for _, o := range sort {
			qry = qry.Order(o)
		}
	} else {
		qry = qry.Order("log_time desc")
	}
	qry = qry.Find(&data, "merchant_id = ?", merchantID)
	qry = qry.Offset(0).Count(&count)
	return data, count, offset
}
