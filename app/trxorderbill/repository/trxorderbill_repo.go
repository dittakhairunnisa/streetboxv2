package repository

import (
	"log"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxorderbill"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxOrderBillRepo ...
type TrxOrderBillRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) trxorderbill.RepoInterface {
	return &TrxOrderBillRepo{db}
}

// Create ...
func (r *TrxOrderBillRepo) Create(entity *entity.TrxOrderBill, db *gorm.DB) (int64, *gorm.DB, error) {
	if err := db.Create(entity).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return 0, nil, err
	}
	return entity.ID, db, nil
}

// CreateOnline from end user apps
func (r *TrxOrderBillRepo) CreateOnline(
	req *model.TrxOrderBill, trxOrderID int64, db *gorm.DB) int64 {
	data := new(entity.TrxOrderBill)
	copier.Copy(&data, req)
	updatedAt := util.MillisToTime(req.UpdatedAt)
	data.BusinessDate = util.MillisToTime(req.BusinessDate)
	data.CreatedAt = util.MillisToTime(req.CreatedAt)
	data.UpdatedAt = &updatedAt
	data.TrxOrderID = trxOrderID
	// data.IsClose = req.IsClose
	data.TotalDiscount = req.TotalDiscount
	data.SubTotal = req.SubTotal
	data.TotalTax = req.TotalTax
	data.GrandTotal = req.GrandTotal

	data.ID = 0
	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return 0
	}
	log.Printf("INFO: Created TrxOrderBill: %+v", data)
	return data.ID
}

// FindAll ..
func (r *TrxOrderBillRepo) FindAll(model *entity.TrxOrderBill) *[]entity.TrxOrderBill {
	data := new([]entity.TrxOrderBill)
	r.DB.Where(model).Find(&data)
	return data
}

// FindByTrxOrderID ..
func (r *TrxOrderBillRepo) FindByTrxOrderID(id int64) *[]entity.TrxOrderBill {
	data := new([]entity.TrxOrderBill)
	r.DB.Where("trx_order_id = ?", id).Find(&data)
	return data
}
