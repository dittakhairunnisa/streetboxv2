package trxvisit

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.TrxVisit) (*gorm.DB, error)
	FindOne(int64) *entity.TrxVisit
	Find(*entity.TrxVisit) *entity.TrxVisit
	GetMerchantIDByTrxID(string) int64
	FindByTrxID(string) *model.TrxVisit
	CheckTrxStatusPendingByTrxHomeVisitSalesID(id int64) int
}
