package trxorderproductsales

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.TrxOrderProductSales, *gorm.DB) (int64, *gorm.DB, error)
	CreateOnline(req *model.TrxOrderProductSales,
		db *gorm.DB, trxOrderID, trxOrderBillID int64) error
	FindAll(*entity.TrxOrderProductSales) *[]entity.TrxOrderProductSales
}
