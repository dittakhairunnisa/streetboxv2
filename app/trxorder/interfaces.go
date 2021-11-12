package trxorder

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ...
type RepoInterface interface {
	Create(string, *model.TrxOrder, *gorm.DB) (int64, *gorm.DB, error)
	FindByTrxID(string) *model.TrxOrderMerchant
	CreateOnline(req *model.ReqTrxOrderOnline) (*gorm.DB, int64, error)
	UpdateByTrxID(data *entity.TrxOrder, trxID string) *entity.TrxOrder
	FindAll(*entity.TrxOrder) *[]entity.TrxOrder
	FindOpenByMerchantUsersID(int64) *[]entity.TrxOrder
}
