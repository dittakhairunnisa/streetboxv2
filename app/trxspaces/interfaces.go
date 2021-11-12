package trxspaces

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*entity.TrxParkingSpaceSales) (*gorm.DB, error)
	Update(*entity.TrxParkingSpaceSales, int64) (*gorm.DB, error)
	GetByUserID(int64) (*[]entity.TrxParkingSpaceSales, error)
	GetMyParking(int64) (*[]model.ResMyParking, error)
	GetSlotMyParking(pspaceSalesID []int64, usersID int64) (*[]model.ResSlotMyParking, error)
	GetAll() *[]model.ResTrxList
	DeleteByID(int64) error
	GetByID(int64) *model.ResTrxList
	GetByMerchantIDAndParkingSalesID(int64, int64) (*entity.TrxParkingSpaceSales, error)
	GetMerchantBySpaceID(int64) *[]model.MerchantSpace
	GetMerchantSchedules(int64, int64) *[]model.Schedules
	GetSchedulesByTypesID(int64) *[]model.Schedules
	GetList(limit, page int, sort []string, filter string) (*[]model.ResTrxList, int, int)
}

// ServiceInterface ...
type ServiceInterface interface {
	CreateTrx(*model.ReqCreateTrxSales, int64) error
	UpdateTrx(*model.ReqCreateTrxSales, int64, int64) error
	GetByUserID(int64) (*[]entity.TrxParkingSpaceSales, error)
	GetMyParking(usersID, merchantID int64) *[]model.ResMyParkingList
	GetSlotMyParking(pspaceSalesID, usersID int64) (*[]model.ResSlotMyParking, error)
	GetAll() *[]model.ResTrxList
	DeleteByID(int64) error
	GetByID(int64) *model.ResTrxList
	GetByMerchantIDAndParkingSalesID(merchantID int64, parkingSpaceSalesID int64) (*entity.TrxParkingSpaceSales, error)
	GetList(limit, page int, sort []string, filter string) model.Pagination
}
