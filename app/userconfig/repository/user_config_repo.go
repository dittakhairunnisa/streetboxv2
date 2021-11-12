package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

type UsersConfigRepo struct {
	db *gorm.DB
}

func (u *UsersConfigRepo) UpdateRadius(rad int) (err error) {
	if err = u.db.Table("users_config").UpdateColumn("radius", rad).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	log.Printf("INFO: Update users radius configuration: %d", rad)
	return
}

func (u *UsersConfigRepo) GetConfig() (cfg entity.UsersConfig, err error) {
	if err = u.db.First(&cfg).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func New(db *gorm.DB) *UsersConfigRepo {
	return &UsersConfigRepo{db: db}
}