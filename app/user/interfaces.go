package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*entity.Users) (*gorm.DB, error)
	GetAll(string) *[]model.ResUserAll
	// GetAllPagination(offset, perPage int, sort string) ([]*entity.Users, int, error)
	FindByID(int64) *model.ResUserAll
	FindUsernameByMultipleID([]int64) *[]string
	DeleteByID(int64) (*gorm.DB, error)
	DeleteByMultipleID(*gorm.DB, []int64) (*gorm.DB, error)
	FindByUsername(string) *entity.Users
	Update(*entity.Users, int64) error
	GetUserAdmin() *[]model.ResUserMerchant
	FindEndUserByID(int64) *entity.Users
	GetByMerchantUsersID(int64) *entity.Users
	GetFoodtruckByTrxVisitID(int64) *entity.Users
}

// ServiceInterface ...
type ServiceInterface interface {
	CreateUser(req model.ReqUserCreate, usersID int64) *entity.Users
	GetAllUser(filter string) *[]model.ResUserAll
	Login(req model.ReqUserLogin, clientID string) string
	LoginGoogle(*model.ReqUserLoginGoogle) string
	GetUserByUserName(string) *entity.Users
	ResetPassword(newPassword, userName string) error
	SendEmailResetPassword(string) bool
	UpdateUser(req model.ReqUserUpdate, id int64) (*entity.Users, error)
	GetUserByID(int64) *model.ResUserAll
	ChangePassword(string, int64) error
	ResetForgotPassword(string, string) error
	GetUserAdmin() *[]model.ResUserMerchant
	DeleteByID(int64, int64) error
	UpdateRole(usersID int64, roleID int64, userID int64) error
	CheckJwt(string) (*jwt.Token, error)
	CreateAddress(addr *entity.UsersAddress) (err error)
	GetPrimaryAddressByUserID(userID int64) (addrs entity.UsersAddress, err error)
	GetAddressByUserID(userID int64) (addrs []entity.UsersAddress, err error)
	DeleteAddress(id, userID int64) (err error)
	UpdateAddress(addr entity.UsersAddress) (err error)
	SwitchAddress(id, userID int64) (err error)
	UpdateRadius(rad int) (err error)
	GetConfig() (cfg entity.UsersConfig, err error)
}
