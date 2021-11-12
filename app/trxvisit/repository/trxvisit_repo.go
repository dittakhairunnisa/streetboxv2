package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxvisit"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// TrxVisitRepo ..
type TrxVisitRepo struct {
	DB *gorm.DB
}

// New ..
func New(db *gorm.DB) trxvisit.RepoInterface {
	return &TrxVisitRepo{db}
}

// Create ..
func (r *TrxVisitRepo) Create(data *entity.TrxVisit) (*gorm.DB, error) {
	db := r.DB.Begin()
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Create TrxVisit : %+v", data)
	return db, nil
}

// FindOne ..
func (r *TrxVisitRepo) FindOne(id int64) *entity.TrxVisit {
	data := new(entity.TrxVisit)
	r.DB.Find(&data)
	return data
}

// Find ..
func (r *TrxVisitRepo) Find(model *entity.TrxVisit) *entity.TrxVisit {
	data := new(entity.TrxVisit)
	r.DB.Where(model).Find(&data)
	return data
}

// GetMerchantIDByTrxID for push notif
func (r *TrxVisitRepo) GetMerchantIDByTrxID(id string) int64 {
	var merchantID int64
	row := r.DB.Select("hs.merchant_id").Joins("JOIN "+
		"trx_homevisit_sales ths on tv.id = ths.trx_visit_id").Joins("JOIN "+
		"homevisit_sales hs on ths.homevisit_sales_id = hs.id").
		Table("trx_visit tv").Where("tv.trx_id = ?", id).Row()
	row.Scan(&merchantID)
	return merchantID
}

// FindByTrxID method for get order detail history end user
func (r *TrxVisitRepo) FindByTrxID(id string) *model.TrxVisit {
	data := new(model.TrxVisit)
	r.DB.Select("tv.trx_id, m.logo, m.name, tv.address, t.created_at, hv.deposit, t.address, "+
		"tv.notes,tv.id, m.phone").Joins("JOIN "+
		"trx_visit tv on t.id = tv.trx_id").Joins("JOIN "+
		"trx_homevisit_sales thv on thv.trx_visit_id = tv.id").Joins("JOIN "+
		"homevisit_sales hv on thv.homevisit_sales_id = hv.id").Joins("JOIN "+
		"merchant m on hv.merchant_id = m.id").Table("trx t").Where("t.deleted_at is null and hv.deleted_at is null and t.id = ?", id).Scan(&data)
	return data
}

// CheckTrxStatusPendingByTrxHomeVisitSalesID to check trx status by trx home visit sales id
func (r *TrxVisitRepo) CheckTrxStatusPendingByTrxHomeVisitSalesID(id int64) int {
	var count int
	r.DB.Select("count(*)").Joins("JOIN "+
		"trx_visit tv on ths.trx_visit_id = tv.id").Joins("JOIN "+
		"trx on tv.trx_id = trx.id").Table("trx_homevisit_sales ths").
		Where("ths.deleted_at is null and ths.id = ? and trx.status = ?", id, "PENDING").Count(&count)
	return count
}
