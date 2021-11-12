package logactivitymerchant

import (
	"streetbox.id/entity"
	"streetbox.id/model"
)

type ServiceInterface interface {
	Create(int64, string)
	GetAll(limit, page int, sort []string, merchantID int64) model.Pagination
}

// RepoInterface ..
type RepoInterface interface {
	Create(*entity.LogActivityMerchant) error
	GetAll(limit, page int, sort []string, merchantID int64) (*[]entity.LogActivityMerchant, int, int)
}
