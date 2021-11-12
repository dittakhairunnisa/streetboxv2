package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

type UserAddressRepo struct {
	db *gorm.DB
}

func (u *UserAddressRepo) Create(addr *entity.UsersAddress) (err error) {
	if err = u.db.Create(&addr).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	} else {
		log.Printf("INFO: Created users address: %v", addr)
	}
	return
}

func (u *UserAddressRepo) GetByUserID(userID int64) (addrs []entity.UsersAddress, err error) {
	if err = u.db.Where("user_id = ?", userID).Order("id desc").Find(&addrs).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (u *UserAddressRepo) GetPrimaryByUserID(userID int64) (addrs entity.UsersAddress, err error) {
	if err = u.db.Where("user_id = ? and \"primary\"", userID).First(&addrs).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (u *UserAddressRepo) Update(addr entity.UsersAddress) (err error) {
	if err = u.db.Save(&addr).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (u *UserAddressRepo) Delete(id, userID int64) (err error) {
	if err = u.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.UsersAddress{}).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (u *UserAddressRepo) Switch(id, userID int64) (err error) {
	tx := u.db.Begin()
	if err = tx.Table("users_address").Where("user_id = ?", userID).UpdateColumn("primary", false).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		tx.Rollback()
		return
	} else if err = tx.Table("users_address").Where("id = ? AND user_id = ?", id, userID).UpdateColumn("primary", true).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func New(db *gorm.DB) (repo *UserAddressRepo) {
	return &UserAddressRepo{db: db}
}
