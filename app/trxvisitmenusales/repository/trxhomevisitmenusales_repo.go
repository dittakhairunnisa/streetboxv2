package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

type TrxHomevisitMenuSalesRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *TrxHomevisitMenuSalesRepo {
	return &TrxHomevisitMenuSalesRepo{
		db: db,
	}
}

func (t *TrxHomevisitMenuSalesRepo) Create(db *gorm.DB, data *entity.TrxHomevisitMenuSales) (err error) {
	if err := db.Create(data).Error; err != nil {
		db.Rollback()
		log.Printf("ERROR: %s", err.Error())
		return err
	}
	return
}
