package trx

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.Trx) (*gorm.DB, error)
	CountTrx() int64
	GetOrderTrx(int64, string, string, string) *model.ResTrxOrderList
	CreateSyncTrx(data *entity.TrxSync)
	FindByID(string) *entity.Trx
	FindByUsersID(limit, page int, usersID int64, filter string) (*[]entity.Trx, int, int)
	GetSyncTrx() *[]entity.TrxSync
	Update(*entity.Trx, string)
	UpdateByTrxID(*entity.Trx, string) (*gorm.DB, error)
	UpdateStatusSyncTrx(string, int64, int, *gorm.DB) (*gorm.DB, error)
	TrxReport(int64, string, string) []model.ResTransactionReport
	TrxReportSingle(int64, string, string) []model.ResTransactionReportSingle
	TrxReportPagination(int64, string, string, int, int, []string) model.Pagination
	TrxReportSinglePagination(int64, string, string, int, int, []string) model.Pagination
	HardDeleteByID(string)
}

// ServiceInterface ..
type ServiceInterface interface {
	CreateTrxVisit(req *model.ReqCreateVisitTrx) (*entity.TrxVisit, error)
	CreateTrxOrder(trxSync *model.ReqTrxOrderList, usersID int64, merchantID int64, uniqueID string) error
	GetOrderTrx(int64, string, string, string) *model.ResTrxOrderList
	CreateSyncTrx(*model.ReqCreateSyncTrx, int64) *entity.TrxSync
	GetOrderHistoryByUsersID(limit, page int, usersID int64, filter string) model.Pagination
	CreateTrxOrderOnline(*model.ReqTrxOrderOnline) (*entity.Trx, error)
	GetSyncTrx() *[]entity.TrxSync
	CountTrx() int64
	GetOnlineOrder(merchantUsers entity.MerchantUsers) *model.ResTrxOrderList
	UpdateStatusSyncTrx(string, int64, int, *gorm.DB) (*gorm.DB, error)
	GetOneTrxOrderByTrxID(string) *model.TrxOrderMerchant
	GetMerchantIDByTrxID(string) int64
	ListBookingTrxVisitSale(merchantID int64, limit, page int, sort []string, filter string) model.Pagination
	ListBookingTrxVisitSalesByID(ID int64, merchantID int64) *model.ResHomeVisitBookingDetailTimeNew
	ClosedOnlineOrderByTrxID(string) *entity.TrxOrder
	GetTrxByID(string) *entity.Trx
	TrxReport(int64, string, string) []model.ResTransactionReport
	TrxReportSingle(int64, string, string) []model.ResTransactionReportSingle
	TrxReportPagination(int64, string, string, int, int, []string) model.Pagination
	TrxReportSinglePagination(int64, string, string, int, int, []string) model.Pagination
	DeleteTrxByID(string)
	VoidTrxByID(string)
	CheckTrxStatusPendingByTrxHomeVisitSalesID(id int64) int
}
