package trxvisitmenusales

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/entity"
)

type Repository interface {
	Create(db *gorm.DB, data *entity.TrxHomevisitMenuSales) (err error)
}
