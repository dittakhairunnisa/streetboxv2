package homevisitsales

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

type RepoInterface interface {
	GetAll(merchantID int64) *[]entity.HomevisitSales
	Create(*entity.HomevisitSales) (*entity.HomevisitSales, error)
	GetAllByMerchantAndDate(merchantID int64, startDate string, endDate string) *[]entity.HomevisitSales
	GetInfoByDate(merchantID int64, date string) *model.ResHomeVisitGetInfo
	GetByID(ID int64) *entity.HomevisitSales
	GetMenuByTrxVisitSalesID(ID int64) (menus []model.ResMenu)
	Update(*entity.HomevisitSales) (*entity.HomevisitSales, error)
	GetAllList(limit, page int) (*[]model.ResVisitSales, int, int)
	DeleteByID(int64, int64) (*gorm.DB, error)
	DeleteByDate(string, int64) (*gorm.DB, error)
	GetAvailableByMerchantID(int64) *[]model.ResVisitSalesDetail
	UpdateAvailableByID(int, int64) error
	UpdateByID(*entity.TrxHomevisitSales, int64) error
	CheckDate(date string, merchantID int64) int
}

type ServiceInterface interface {
	GetAll(merchantID int64) *[]entity.HomevisitSales
	Create(*entity.HomevisitSales) (*entity.HomevisitSales, error)
	GetAllByMerchantAndDate(merchantID int64, startDate string, endDate string) *[]entity.HomevisitSales
	GetInfoByDate(merchantID int64, date string) *model.ResHomeVisitGetInfo
	GetByID(ID int64) *entity.HomevisitSales
	Update(req *entity.HomevisitSales) (*entity.HomevisitSales, error)
	DeleteByDate(string, int64) (*gorm.DB, error)
	DeleteByID(int64, int64) (*gorm.DB, error)
	GetAllEndUser(limit, page int) model.Pagination
	GetAvailableByMerchantID(int64) *[]model.ResVisitSalesDetail
	UpdateByTrxID(string) error
	CheckDate(date string, merchantID int64) int
}
