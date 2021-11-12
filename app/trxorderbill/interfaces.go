package trxorderbill

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*entity.TrxOrderBill, *gorm.DB) (int64, *gorm.DB, error)
	CreateOnline(*model.TrxOrderBill, int64, *gorm.DB) int64
	FindAll(*entity.TrxOrderBill) *[]entity.TrxOrderBill
	FindByTrxOrderID(int64) *[]entity.TrxOrderBill
}
