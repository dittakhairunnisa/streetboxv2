package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/logxenditreq"
	"streetbox.id/entity"
)

// LogXenditReqRepo ..
type LogXenditReqRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) logxenditreq.RepoInterface {
	return &LogXenditReqRepo{db}
}

// Create ..
func (r *LogXenditReqRepo) Create(db *gorm.DB, data *entity.LogXenditRequest) error {
	trx := r.DB
	if db != nil {
		trx = db
	}
	if err := trx.Create(&data).Error; err != nil {
		return err
	}
	log.Printf("INFO: Create Log Xendit Request Success")
	return nil
}
