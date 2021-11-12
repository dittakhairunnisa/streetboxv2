package trxvisitsales

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*gorm.DB, *entity.TrxHomevisitSales) error
	FindByTrxID(string) *[]model.TrxHomevisitSales
	FindByMerchantID(int64) *[]model.HomeVisitSales
	ListBookingTrxVisitSale(merchantID int64, limit, page int, sort []string, filter string) model.Pagination
	ListBookingTrxVisitSalesByID(ID int64, merchantID int64) *model.ResHomeVisitBookingDetailTimeNew
	GetHomeVisitData(date string, merchantID int64) []int64
	UpdateByID(*entity.TrxHomevisitSales, int64) error
}
