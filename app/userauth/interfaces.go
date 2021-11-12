package userauth

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*gorm.DB, *entity.UsersAuth) error
	DeleteByID(*gorm.DB, string) error
	DeleteByMultipleUsername(*gorm.DB, []string) (*gorm.DB, error)
	FindByUsernameAndPassword(userName, password string) bool
	ResetPassword(userName, newPassword string) error
}
