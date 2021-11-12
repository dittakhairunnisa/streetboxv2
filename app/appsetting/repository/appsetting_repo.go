package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/app/appsetting"
	"streetbox.id/model"
)

// AppSettingRepo ...
type AppSettingRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) appsetting.RepoInterface {
	return &AppSettingRepo{db}
}

// GetByKey ... appsetting
func (r *AppSettingRepo) GetByKey(key string) *model.AppSetting {
	data := new(model.AppSetting)
	if err := r.DB.Select("a.*").Where("a.key = ?", key).Table("app_setting a").Scan(&data).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	if data.ID == 0 {
		return nil
	}
	return data
}

// UpdateByKey ...
func (r *AppSettingRepo) UpdateByKey(key string, value string) error {
	trx := r.DB.Begin()
	if err := trx.Where("a.key = ?", key).Table("app_setting a").Update("value", value).Error; err != nil {
		trx.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return nil
	}
	trx.Commit()
	return nil
}