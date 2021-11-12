package userrole

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*gorm.DB, *entity.UsersRole) error
	DeleteByID(*gorm.DB, int64) error
	DeleteByMultipleID(*gorm.DB, []int64) error
	GetNameByUserID(int64) string
	Update(usersID, roleID int64) error
}
