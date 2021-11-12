package merchant

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// ServiceInterface ...
type ServiceInterface interface {
	CreateMerchant(req *model.ReqCreateMerchant, usersID int64) (*entity.Merchant, error)
	XenditGenerateSubAccount(req *model.ReqXenditGenerateSubAccount, usersID int64) (*entity.Merchant, error)
	CreateFoodtruck(req *model.ReqCreateFoodtruck, merchantID int64) (*entity.Users, error)
	CreateMenu(req *model.ReqCreateMerchantMenu, merchantID int64) (*entity.MerchantMenu, error)
	CreateMenuWithImage(req *model.ReqCreateMerchantMenu, image string, merchantID int64) (*entity.MerchantMenu, error)
	CreateTax(req *model.MerchantTax, merchantID int64) (*entity.MerchantTax, error)
	GetAllFoodtruck(merchantID int64) *[]entity.Users
	GetInfo(usersID int64) *model.Merchant
	IsExist(usersID int64) bool
	CreateShift(usersID int64, shift string) (*entity.MerchantUsersShift, error)
	Update(*model.ReqUpdateMerchant, int64) (*entity.Merchant, error)
	UpdateFoodtruck(*model.ReqUserUpdate, int64) (*entity.Users, error)
	UpdateMenu(req *model.ReqUpdateMerchantMenu, merchantID int64, ID int64) (*entity.MerchantMenu, error)
	UpdateTax(req *model.MerchantTax, ID int64, merchantID int64) (*entity.MerchantTax, error)
	GetAll() *[]model.Merchant
	DeleteByMerchantID(id int64, UserID int64) error
	GetByID(int64) *entity.Merchant
	IsUsersShiftIn(usersID int64) bool
	UploadLogo(string, int64) error
	UploadBanner(string, int64) error
	UploadMenu(string, int64, int64) error
	GetFoodtruckByID(int64) *model.ResUserAll
	DeleteFoodTruckByID(id int64) error
	DeleteMenu(ID int64) error
	GetFoodtruckTasks(merchantID int64) *[]model.ResGetFoodtruckTasks
	GetMenuPagination(int64, int, int, []string) model.Pagination
	CountFoodtruckByMerchantID(merchantID int64) int
	GetMenuList(merchantID int64, nearby, visit bool) *[]entity.MerchantMenu
	GetMenuByID(merchantID int64, ID int64) *entity.MerchantMenu
	GetTax(merchantID int64) *entity.MerchantTax
	RegistrationToken(string, int64) error
	GetMerchantUsersByUsersID(int64) *entity.MerchantUsers
	GetMerchantUsersAdminByMerchantID(int64) *entity.MerchantUsers
	GetMerchantUsersByID(int64) *entity.MerchantUsers
	RemoveImage(filename, types string, merchantID int64) error
	RemoveImageMenu(menu *entity.MerchantMenu, ID int64) error
	CreateCategory(cat *entity.MerchantCategory) (err error)
	GetAllCategory() (cats []entity.MerchantCategory, err error)
	UpdateCategory(cat *entity.MerchantCategory) (err error)
	DeleteCategory(id int64) (err error)
	CheckStock(reqProductSales []model.TrxOrderProductSales) (err error)
}

// RepoInterface ...
type RepoInterface interface {
	Create(*entity.Merchant) (*gorm.DB, error)
	GetAllFoodtruck(merchantID int64) *[]entity.Users
	GetInfo(usersID int64) *model.Merchant
	Update(*entity.Merchant, int64) error
	GetAll() *[]model.Merchant
	DeleteByID(id int64) (*gorm.DB, error)
	GetByID(int64) *entity.Merchant
	GetFoodtruckTasks(merchantID int64) *[]model.ResGetFoodtruckTasks
	GetNearby(limit int, page int, lat float64,
		lon float64, distance float64) (*model.NearbySorted, int, int)
	RemoveImage(string, int64) error
}
