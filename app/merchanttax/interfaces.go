package merchanttax

import (
	"streetbox.id/entity"
)

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.MerchantTax) (*entity.MerchantTax, error)
	Update(*entity.MerchantTax, int64, int64) (*entity.MerchantTax, error)
	GetTax(int64) *entity.MerchantTax
	Find(*entity.MerchantTax) *entity.MerchantTax
}
