package repository

import (
	"github.com/jinzhu/gorm"
	"streetbox.id/app/paymentmethod"
	"streetbox.id/model"
)

// PaymentMethodRepo ...
type PaymentMethodRepo struct {
	DB *gorm.DB
}

// New ...
func New(db *gorm.DB) paymentmethod.RepoInterface {
	return &PaymentMethodRepo{db}
}

// FindByActive ...
func (r *PaymentMethodRepo) FindByActive() *[]model.ResPaymentMethod {
	data := new([]model.ResPaymentMethod)
	r.DB.Select("pm.id, pm.name, pm.types, pm.is_active, "+
		"p.name as provider_name").Joins("JOIN "+
		"payment_method pm on p.id = pm.payment_provider_id").
		Table("payment_provider p").Where("pm.deleted_at is null and "+
		"p.deleted_at is null and pm.is_active = ?", true).Scan(&data)
	return data
}
