package merchantcategory

import "streetbox.id/entity"

type RepoInterface interface {
	Create(cat *entity.MerchantCategory) (err error)
	GetAll() (cats []entity.MerchantCategory, err error)
	Update(cat *entity.MerchantCategory) (err error)
	Delete(id int64) (err error)
}
