package repository

import (
	"log"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"streetbox.id/app/trxorderpaymentsales"
	"streetbox.id/entity"
	"streetbox.id/model"
	"streetbox.id/util"
)

// TrxOrderPaymentSalesRepo ..
type TrxOrderPaymentSalesRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) trxorderpaymentsales.RepoInterface {
	return &TrxOrderPaymentSalesRepo{db}
}

// CreateOnline from end user apps
func (r *TrxOrderPaymentSalesRepo) CreateOnline(
	req *model.TrxOrderPaymentSales, trxOrderBillID int64, db *gorm.DB) error {
	data := new(entity.TrxOrderPaymentSales)
	copier.Copy(&data, req)
	updatedAt := util.MillisToTime(req.UpdatedAt)
	data.CreatedAt = util.MillisToTime(req.CreatedAt)
	data.UpdatedAt = &updatedAt
	data.TrxOrderBillID = trxOrderBillID
	data.ID = 0
	if err := db.Create(data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Created TrxOrderPaymentSales: %+v", data)
	return nil
}

// CreateOffline ...
func (r *TrxOrderPaymentSalesRepo) CreateOffline(data *entity.TrxOrderPaymentSales, db *gorm.DB) (*gorm.DB, error) {
	if err := db.Create(data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return db, nil
}

// FindAll ..
func (r *TrxOrderPaymentSalesRepo) FindAll(model *entity.TrxOrderPaymentSales) *[]entity.TrxOrderPaymentSales {
	data := new([]entity.TrxOrderPaymentSales)
	r.DB.Where(model).Find(&data)
	return data
}

// FindPaymentName method for get payment name order detail history
func (r *TrxOrderPaymentSalesRepo) FindPaymentName(trxOrderBillID int64) *string {
	var paymentName string
	row := r.DB.Select("pm.name").Joins("JOIN "+
		"payment_method pm on tops.payment_method_id = pm.id").
		Table("trx_order_payment_sales tops").
		Where("tops.trx_order_bill_id = ?", trxOrderBillID).Row()
	row.Scan(&paymentName)
	return &paymentName
}
