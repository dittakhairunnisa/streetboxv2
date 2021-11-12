package trxordertaxsales

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	CreateOnline(req *model.TrxOrderTaxSales,
		trxOrderID, trxOrderBillID int64, db *gorm.DB) error
	CreateOffline(req *entity.TrxOrderTaxSales, db *gorm.DB) (*gorm.DB, error)
	FindAll(*entity.TrxOrderTaxSales) *[]entity.TrxOrderTaxSales
	Find(*entity.TrxOrderTaxSales) *entity.TrxOrderTaxSales
}
