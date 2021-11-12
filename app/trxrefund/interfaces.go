package trxrefund

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ...
type RepoInterface interface {
	CreateRefund(*entity.TrxRefund) (int64, *gorm.DB, error)
	CreateRefundSpace(*entity.TrxRefundSpace, *gorm.DB) (*gorm.DB, error)
	CreateRefundVisit(data *entity.TrxRefundVisit, db *gorm.DB) (*gorm.DB, error)
}

// ServiceInterface ...
type ServiceInterface interface {
	CreateRefundParkingSpace(req *model.ReqRefundParkingSpaceSales) error
	CreateRefundHomeVisit(req *model.ReqRefundHomeVisit, merchantID int64) error
}
