package logactivity

import (
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ..
type RepoInterface interface {
	Create(activity string)
	GetAllPagination(limit int, page int, sort []string) (*[]entity.LogActivity, int, int)
	GetList() *[]entity.LogActivity
}

// ServiceInterface ..
type ServiceInterface interface {
	GetAll(limit int, page int, sort []string) model.Pagination
	GetList() *[]entity.LogActivity
}
