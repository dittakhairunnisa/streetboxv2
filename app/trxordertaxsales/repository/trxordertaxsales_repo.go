package repository

import (
	"log"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxordertaxsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxOrderTaxSalesRepo ..
type TrxOrderTaxSalesRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) trxordertaxsales.RepoInterface {
	return &TrxOrderTaxSalesRepo{db}
}

// CreateOnline from end user apps
func (r *TrxOrderTaxSalesRepo) CreateOnline(
	req *model.TrxOrderTaxSales,
	trxOrderID, trxOrderBillID int64, db *gorm.DB) error {
	data := new(entity.TrxOrderTaxSales)
	copier.Copy(&data, req)
	updatedAt := util.MillisToTime(req.UpdatedAt)
	data.UpdatedAt = &updatedAt
	data.CreatedAt = util.MillisToTime(req.CreatedAt)
	data.TrxOrderBillID = trxOrderBillID
	data.ID = 0
	if err := db.Create(data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created TrxOrderTaxSales: %+v", data)
	return nil
}

// CreateOffline ..
func (r *TrxOrderTaxSalesRepo) CreateOffline(data *entity.TrxOrderTaxSales, db *gorm.DB) (*gorm.DB, error) {
	if err := db.Create(data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return db, nil
}

// FindAll ..
func (r *TrxOrderTaxSalesRepo) FindAll(model *entity.TrxOrderTaxSales) *[]entity.TrxOrderTaxSales {
	data := new([]entity.TrxOrderTaxSales)
	r.DB.Where(model).Find(&data)
	return data
}

// Find ..
func (r *TrxOrderTaxSalesRepo) Find(model *entity.TrxOrderTaxSales) *entity.TrxOrderTaxSales {
	data := new(entity.TrxOrderTaxSales)
	r.DB.Where(model).Find(&data)
	return data
}
