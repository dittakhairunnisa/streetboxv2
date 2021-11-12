package merchantusers

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// RepoInterface ...
type RepoInterface interface {
	Create(db *gorm.DB, merchantID, usersID int64) error
	IsExist(usersID int64) bool
	GetByUsersID(usersID int64) *entity.MerchantUsers
	GetAdminByMerchantID(merchantID int64) *entity.MerchantUsers
	DeleteByMerchantID(*gorm.DB, int64) error
	DeleteByFoodtruckID(int64) error
	GetUserIdsByMerchantID(merchantID int64) *[]int64
	CountFoodtruckByMerchantID(merchantID int64) int
	Update(*entity.MerchantUsers, int64) error
	GetOne(int64) *entity.MerchantUsers
}
