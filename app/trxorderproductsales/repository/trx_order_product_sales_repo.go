package repository

import (
	"log"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxorderproductsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxOrderProductSalesRepo ..
type TrxOrderProductSalesRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) trxorderproductsales.RepoInterface {
	return &TrxOrderProductSalesRepo{db}
}

// Create ...
func (r *TrxOrderProductSalesRepo) Create(entity *entity.TrxOrderProductSales, db *gorm.DB) (int64, *gorm.DB, error) {
	if err := db.Create(entity).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return 0, nil, err
	}
	return entity.ID, db, nil
}

// CreateOnline from end user apps
func (r *TrxOrderProductSalesRepo) CreateOnline(
	req *model.TrxOrderProductSales, db *gorm.DB, trxOrderID, trxOrderBillID int64) error {
	data := new(entity.TrxOrderProductSales)
	copier.Copy(&data, req)
	updatedAt := util.MillisToTime(req.UpdatedAt)
	data.BusinessDate = util.MillisToTime(req.BusinessDate)
	data.CreatedAt = util.MillisToTime(req.CreatedAt)
	data.UpdatedAt = &updatedAt
	data.TrxOrderBillID = trxOrderBillID
	data.ID = 0
	if err := db.Create(data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created TrxOrderProductSales: %+v", data)
	return nil
}

// FindAll ..
func (r *TrxOrderProductSalesRepo) FindAll(model *entity.TrxOrderProductSales) *[]entity.TrxOrderProductSales {
	data := new([]entity.TrxOrderProductSales)
	r.DB.Where(model).Find(&data)
	return data
}
