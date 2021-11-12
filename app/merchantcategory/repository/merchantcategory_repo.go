package merchantcategory

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

type MerchantCategoryRepo struct {
	db *gorm.DB
}

func (m *MerchantCategoryRepo) Create(cat *entity.MerchantCategory) (err error) {
	if err = m.db.Create(&cat).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	} else {
		log.Printf("INFO: Created merchant category: %s", cat.Category)
	}
	return
}

func (m *MerchantCategoryRepo) GetAll() (cats []entity.MerchantCategory, err error) {
	if err = m.db.Find(&cats).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (m *MerchantCategoryRepo) Update(cat *entity.MerchantCategory) (err error) {
	if err = m.db.Save(&cat).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func (m *MerchantCategoryRepo) Delete(id int64) (err error) {
	if err := m.db.Where("id = ?", id).
		Delete(&entity.MerchantCategory{}).Error; err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	return
}

func New(db *gorm.DB) *MerchantCategoryRepo {
	return &MerchantCategoryRepo{db: db}
}
