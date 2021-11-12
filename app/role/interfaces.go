package role

import (
	"streetbox.id/entity"
	"streetbox.id/model"
)

// RepoInterface ...
type RepoInterface interface {
	Create(*entity.Role) error
	FindByName(string) *entity.Role
	DeleteByID(int64) error
	GetAll() *[]entity.Role
	GetAllExclude() *[]entity.Role
	GetOne(int64) *entity.Role
}

// ServiceInterface ...
type ServiceInterface interface {
	Create(model.ReqRoleCreate) (*entity.Role, error)
	SearchByName(string) *entity.Role
	DeleteByID(int64) error
	GetAll() *[]entity.Role
	GetAllExclude() *[]entity.Role
}
