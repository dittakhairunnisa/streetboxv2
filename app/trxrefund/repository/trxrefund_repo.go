package repository

import (
	"log"

	"streetbox.id/app/trxrefund"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// TrxRefundRepo ...
type TrxRefundRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) trxrefund.RepoInterface {
	return &TrxRefundRepo{db}
}

// CreateRefund ...
func (r *TrxRefundRepo) CreateRefund(data *entity.TrxRefund) (int64, *gorm.DB, error) {
	db := r.DB.Begin()
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("Cannot Create Refund")
		return 0, nil, err
	}
	return data.ID, db, nil
}

// CreateRefundSpace ...
func (r *TrxRefundRepo) CreateRefundSpace(data *entity.TrxRefundSpace, db *gorm.DB) (*gorm.DB, error) {
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("Cannot Create Refund Space")
		return nil, err
	}
	return db, nil
}

// CreateRefundVisit ...
func (r *TrxRefundRepo) CreateRefundVisit(data *entity.TrxRefundVisit, db *gorm.DB) (*gorm.DB, error) {
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("Cannot Create Refund Visit")
		return nil, err
	}
	return db, nil
}
