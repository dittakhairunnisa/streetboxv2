package trxorderpaymentsales

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	CreateOnline(req *model.TrxOrderPaymentSales, trxOrderBillID int64, db *gorm.DB) error
	CreateOffline(req *entity.TrxOrderPaymentSales, db *gorm.DB) (*gorm.DB, error)
	FindAll(*entity.TrxOrderPaymentSales) *[]entity.TrxOrderPaymentSales
	FindPaymentName(trxOrderBillID int64) *string
}
