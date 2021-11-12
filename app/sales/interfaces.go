package sales

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface Sales Repo Interface
type RepoInterface interface {
	Create(*entity.ParkingSpaceSales) error
	GetOne(int64) *entity.ParkingSpaceSales
	FindBySpaceID(parkingSpaceID int64,
		limit, page int, sort []string) (*[]model.ResSales, int, int)
	UpdateAvailableSlot(db *gorm.DB, qty int, id int64) (*gorm.DB, error)
	Update(*entity.ParkingSpaceSales, int64) (*entity.ParkingSpaceSales, error)
	DeleteByID(int64) error
	FindLikeName(name string, limit,
		page int, sort []string) (*[]model.ResSales, int, int)
	FindLikeNameBackoffice(name string, limit,
		page int, sort []string) (*[]model.ResSales, int, int)
	GetAll(limit, page int, sort []string) (*[]model.ResSales, int, int)
	GetAllBackoffice(limit, page int, sort []string) (*[]model.ResSales, int, int)
	GetAllBackofficeNonPaginate(string) *[]model.ResSales
	GetAllList() *[]entity.ParkingSpaceSales
	GetSalesBySpace(salesID int64, startDate string, endDate string) *[]entity.ParkingSpaceSales
	GetSalesIDByPSpaceID(id int64, usersID int64) []int64
	GetSalesNearby(userLat, userLon, distance float64) *[]model.ResParkingSpace
}

// ServiceInterface ...
type ServiceInterface interface {
	CreateSales(*model.ReqSalesCreate, int64) (*entity.ParkingSpaceSales, error)
	GetAll(limit, page int, sort []string) model.Pagination
	GetAllBackoffice(limit, page int, sort []string) model.Pagination
	GetBySpaceID(id int64, limit, page int, sort []string) model.Pagination
	UpdateAvailableSlot(*gorm.DB, int64, int) (*gorm.DB, error)
	GetByID(id int64) *entity.ParkingSpaceSales
	Update(*model.ReqSalesUpdate, int64, int64) (*entity.ParkingSpaceSales, error)
	DeleteByID(id int64, userID int64) error
	FindLikeName(name string,
		limit, page int, sort []string) model.Pagination
	FindLikeNameBackoffice(name string,
		limit, page int, sort []string) model.Pagination
	GetAllList() *[]entity.ParkingSpaceSales
	GetSalesBySpace(salesID int64, startDate string, endDate string) *[]entity.ParkingSpaceSales
	GetAllListNonPaginate(string) *[]model.ResSales
}
