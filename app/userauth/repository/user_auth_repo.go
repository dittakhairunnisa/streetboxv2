package repository

import (
	"log"

	"streetbox.id/app/userauth"
	"streetbox.id/util"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

// UserAuthRepo ...
type UserAuthRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) userauth.RepoInterface {
	return &UserAuthRepo{db}
}

// Create ...
func (repo *UserAuthRepo) Create(trx *gorm.DB, userAuth *entity.UsersAuth) error {
	if err := trx.Create(userAuth).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		trx.Rollback()
		return err
	}
	log.Printf("INFO: Created UsersAuth: %+v", userAuth)
	return nil
}

// DeleteByID ...
func (repo *UserAuthRepo) DeleteByID(trx *gorm.DB, userName string) error {
	data := new(entity.UsersAuth)
	if err := trx.Where("user_name = ?", userName).
		Delete(&data).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	log.Printf("INFO: Deleted UsersAuth: %+v", data)
	return nil
}

// DeleteByMultipleUsername ...
func (repo *UserAuthRepo) DeleteByMultipleUsername(trx *gorm.DB, userNames []string) (*gorm.DB, error) {
	data := new([]entity.UsersAuth)
	if err := trx.Where("user_name IN (?)", userNames).
		Delete(&data).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: Deleted UsersAuth: %+v", data)
	return trx, nil
}

// FindByUsernameAndPassword ...
func (repo *UserAuthRepo) FindByUsernameAndPassword(userName, password string) bool {
	userAuth := new(entity.UsersAuth)
	if err := repo.DB.Find(&userAuth, "user_name = ?", userName).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	match := util.CheckPasswordHash(password, userAuth.Password)
	return match
}

// ResetPassword ...
func (repo *UserAuthRepo) ResetPassword(userName, newPassword string) error {
	userAuth := new(entity.UsersAuth)
	userAuth.Password = newPassword
	return repo.DB.Model(&entity.UsersAuth{UserName: userName}).
		Update(&userAuth).Error
}
